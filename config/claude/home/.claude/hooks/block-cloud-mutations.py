#!/usr/bin/env python3
"""Deny two classes of dangerous cloud CLI calls from Claude Code.

1. `gcloud spanner ...` against a production instance/database/project.
   Production data is meant to be read through a read-only path (an MCP server
   or BigQuery), never poked at directly with the CLI.
2. `bq` writes: the mutating subcommands, and `bq query` whose SQL is DML/DDL.

The command string is tokenized and each shell segment inspected on its own, so
unrelated commands that merely mention `bq` or `spanner` are not blocked.
Nested shells (`bash -c '...'`) and command substitutions are recursed into.
"""
import json
import re
import shlex
import sys

# Wrappers that may precede the real command without changing its meaning.
WRAPPERS = {"sudo", "command", "env", "time", "nohup", "rtk", "proxy", "xargs"}
# Shells whose `-c` argument is itself a command string.
SHELLS = {"bash", "sh", "zsh", "dash", "ksh", "fish"}
ENV_ASSIGN = re.compile(r"^[A-Za-z_][A-Za-z0-9_]*=")

SHELL_OPERATOR_CHARS = set(";&|()\n")
MAX_DEPTH = 4

# Production-looking identifier inside an instance / database / project name.
PRODUCTION = re.compile(r"(?:^|[-_.=/])(prd|prod|production)(?:[-_.]|$)", re.IGNORECASE)

# `bq` subcommands that write. `extract` (export to GCS), `ls`, `show`, `head`,
# `mkdef` and `wait` stay allowed.
BQ_WRITE_COMMANDS = {
    "insert",
    "load",
    "rm",
    "cp",
    "mk",
    "update",
    "partition",
    "cancel",
    "set-iam-policy",
    "add-iam-policy-binding",
    "remove-iam-policy-binding",
}

# Every `bq` subcommand, so the one in a command line can be found by name.
BQ_COMMANDS = BQ_WRITE_COMMANDS | {
    "query",
    "ls",
    "show",
    "head",
    "extract",
    "mkdef",
    "wait",
    "version",
    "help",
    "init",
    "shell",
    "get-iam-policy",
    "test-iam-permissions",
}

# DML / DDL statements. `CREATE TEMP FUNCTION` is excluded: it is routine inside
# read-only analytical queries.
SQL_WRITE = re.compile(
    r"\b(?:"
    r"insert\s+into|insert\s+all"
    r"|update\s+[`\w]"
    r"|delete\s+from"
    r"|merge\s+(?:into\s+)?[`\w]"
    r"|truncate\s+table"
    r"|create\s+(?!temp\w*\s+function\b)(?:or\s+replace\s+)?(?:temp\w*\s+)?"
    r"(?:table|view|materialized|external|schema|dataset|function|procedure|snapshot|index|model|reservation|assignment|capacity|search|vector|row\s+access)"
    r"|drop\s+(?:table|view|materialized|external|schema|dataset|function|procedure|snapshot|index|model|row\s+access|reservation|assignment|capacity|search|vector|if)"
    r"|alter\s+(?:table|view|materialized|schema|column|organization|project|bi_capacity)"
    r"|grant\s+[`\w]|revoke\s+[`\w]"
    r"|export\s+data|load\s+data"
    # No trailing \b: branches such as `update\s+[`\w]` already end mid-word.
    r")",
    re.IGNORECASE,
)

REASON_SPANNER = (
    "gcloud spanner against production is blocked by your global Claude Code hook. "
    "Read production data through a read-only path (the project's MCP server or "
    "BigQuery) instead."
)
REASON_BQ = (
    "BigQuery write operations are blocked by your global Claude Code hook "
    "(mutating bq subcommands and DML/DDL in bq query)."
)


def tokenize(command):
    """Tokenize a command string, or None when it cannot be parsed."""
    try:
        lexer = shlex.shlex(command, posix=True, punctuation_chars=True)
        lexer.whitespace_split = True
        return list(lexer)
    except ValueError:
        return None


def split_segments(tokens):
    """Split a token list on shell operators into individual commands."""
    segments, current = [], []
    for token in tokens:
        if token and set(token) <= SHELL_OPERATOR_CHARS:
            if current:
                segments.append(current)
                current = []
        else:
            current.append(token)
    if current:
        segments.append(current)
    return segments


def substitutions(command):
    """Yield the bodies of `$(...)` and backtick command substitutions."""
    i = 0
    while i < len(command) - 1:
        if command[i] == "$" and command[i + 1] == "(":
            depth, start = 1, i + 2
            j = start
            while j < len(command) and depth:
                if command[j] == "(":
                    depth += 1
                elif command[j] == ")":
                    depth -= 1
                j += 1
            if not depth:
                yield command[start : j - 1]
            i = j
            continue
        i += 1
    for body in re.findall(r"`([^`]*)`", command):
        yield body


def strip_prefix(segment):
    """Drop leading env assignments and wrapper commands."""
    i = 0
    while i < len(segment) and (ENV_ASSIGN.match(segment[i]) or segment[i] in WRAPPERS):
        i += 1
    return segment[i:]


def nested_commands(segment):
    """Yield command strings passed to a nested shell via `-c`."""
    segment = strip_prefix(segment)
    if not segment or segment[0].rsplit("/", 1)[-1] not in SHELLS:
        return
    for i, token in enumerate(segment[1:], start=1):
        # `-c`, and clustered forms such as `-lc` / `-ec`.
        if re.fullmatch(r"-[A-Za-z]*c", token) and i + 1 < len(segment):
            yield segment[i + 1]
            return


def strip_sql_noise(sql):
    """Remove comments and string literals so keywords inside them do not match."""
    sql = re.sub(r"--[^\n]*", " ", sql)
    sql = re.sub(r"/\*.*?\*/", " ", sql, flags=re.DOTALL)
    sql = re.sub(r"'''.*?'''|\"\"\".*?\"\"\"", " ", sql, flags=re.DOTALL)
    return re.sub(r"'[^']*'|\"[^\"]*\"", " ", sql)


def spanner_production(segment):
    """True for `gcloud spanner ...` that names a production resource."""
    segment = strip_prefix(segment)
    if len(segment) < 2 or segment[0].rsplit("/", 1)[-1] != "gcloud":
        return False
    if "spanner" not in segment[1:]:
        return False
    return any(PRODUCTION.search(token) for token in segment[1:])


def bq_write(segment, command):
    """True for a mutating `bq` invocation."""
    segment = strip_prefix(segment)
    if not segment or segment[0].rsplit("/", 1)[-1] != "bq":
        return False
    # Global flags may take a separate value (`bq --project_id foo query ...`),
    # so locate the subcommand by name rather than by position.
    index = next(
        (i for i, token in enumerate(segment[1:], start=1) if token in BQ_COMMANDS), None
    )
    if index is None:
        return False
    subcommand = segment[index]
    if subcommand in BQ_WRITE_COMMANDS:
        return True
    if subcommand != "query":
        return False
    sql = " ".join(token for token in segment[index + 1 :] if not token.startswith("-"))
    if sql:
        return bool(SQL_WRITE.search(strip_sql_noise(sql)))
    # The SQL may arrive on stdin (`echo '...' | bq query`). Literals are not
    # stripped there: the SQL itself is quoted in that shape.
    return bool(SQL_WRITE.search(command))


def check(command, depth=0):
    """Return a deny reason, or None."""
    if depth > MAX_DEPTH:
        return None
    lowered = command.lower()
    if "spanner" not in lowered and "bq" not in lowered:
        return None
    tokens = tokenize(command)
    if tokens is None:
        # Unparseable (e.g. unbalanced quotes): fail closed only on a shape
        # that clearly is one of these calls.
        if re.search(r"\bgcloud\b[^;&|\n]*\bspanner\b", command) and PRODUCTION.search(command):
            return REASON_SPANNER
        if re.search(r"\bbq\s+(?:-\S+\s+)*(%s)\b" % "|".join(BQ_WRITE_COMMANDS), command):
            return REASON_BQ
        if re.search(r"\bbq\s+(?:-\S+\s+)*query\b", command) and SQL_WRITE.search(
            strip_sql_noise(command)
        ):
            return REASON_BQ
        return None
    for segment in split_segments(tokens):
        if spanner_production(segment):
            return REASON_SPANNER
        if bq_write(segment, command):
            return REASON_BQ
        for nested in nested_commands(segment):
            reason = check(nested, depth + 1)
            if reason:
                return reason
    for body in substitutions(command):
        reason = check(body, depth + 1)
        if reason:
            return reason
    return None


def main():
    try:
        payload = json.load(sys.stdin)
    except (ValueError, OSError):
        return
    reason = check(payload.get("tool_input", {}).get("command", ""))
    if not reason:
        return
    print(
        json.dumps(
            {
                "hookSpecificOutput": {
                    "hookEventName": "PreToolUse",
                    "permissionDecision": "deny",
                    "permissionDecisionReason": reason,
                }
            }
        )
    )


if __name__ == "__main__":
    main()
    sys.exit(0)
