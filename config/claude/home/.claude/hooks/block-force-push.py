#!/usr/bin/env python3
"""Deny only genuine `git push` force/delete invocations.

Regex-matching the whole command string blocks unrelated commands
(`git clean -f`, `git checkout -f`, `git branch -d`, `rm -f`, ...), so this
tokenizes the command, isolates each `git push` invocation, and inspects only
that invocation's arguments. Nested shells (`bash -c '...'`) and command
substitutions are recursed into, otherwise they are a trivial bypass.
"""
import json
import re
import shlex
import sys

# Wrappers that may precede `git` without changing the meaning of the command.
WRAPPERS = {"sudo", "command", "env", "time", "nohup", "rtk", "proxy"}
# Shells whose `-c` argument is itself a command string.
SHELLS = {"bash", "sh", "zsh", "dash", "ksh", "fish"}
ENV_ASSIGN = re.compile(r"^[A-Za-z_][A-Za-z0-9_]*=")

# `git` global options that consume a separate value.
GIT_OPTS_WITH_VALUE = {
    "-C",
    "-c",
    "--git-dir",
    "--work-tree",
    "--namespace",
    "--exec-path",
    "--super-prefix",
}
# `git push` options that consume a separate value.
PUSH_OPTS_WITH_VALUE = {"-o", "--push-option", "--repo", "--receive-pack", "--exec"}

FORCE_OPTS = {"--force", "--force-with-lease", "--force-if-includes", "--mirror", "--delete"}

DENY_REASON = "git push force/delete operations are blocked by your global Claude Code hook."

SHELL_OPERATOR_CHARS = set(";&|()\n")

MAX_DEPTH = 4

# Fallback for command strings shlex cannot parse (unbalanced quotes). A bare
# `-f` must not be enough: `git` and `push` have to appear before the force
# flag inside one operator-free window, or `grep -f` in an unparseable command
# would be denied.
COARSE_FORCE_PUSH = re.compile(
    r"\bgit\b[^;&|\n]{0,120}?\bpush\b"
    r"[^;&|\n]*?(?:\s(?:--force[-\w]*|--mirror|--delete|-[A-Za-z]*[fd])(?=\s|=|$)"
    r"|\s[:+]\S)"
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


def push_args(segment):
    """Return the argument list of `git push` for this segment, else None."""
    segment = strip_prefix(segment)
    if not segment or segment[0].rsplit("/", 1)[-1] != "git":
        return None
    i = 1
    # `git` global options, before the subcommand.
    while i < len(segment) and segment[i].startswith("-"):
        i += 2 if segment[i] in GIT_OPTS_WITH_VALUE else 1
    if i >= len(segment) or segment[i] != "push":
        return None
    return segment[i + 1 :]


def is_force_push(args):
    positional_only = False
    i = 0
    while i < len(args):
        token = args[i]
        i += 1
        if token == "--":
            positional_only = True
            continue
        if not positional_only and token.startswith("-") and token != "-":
            name = token.split("=", 1)[0]
            if name in FORCE_OPTS:
                return True
            # Clustered short flags, e.g. `-fu` / `-d`.
            if re.fullmatch(r"-[A-Za-z]+", token) and ("f" in token[1:] or "d" in token[1:]):
                return True
            if name in PUSH_OPTS_WITH_VALUE and "=" not in token:
                i += 1
            continue
        # Positional argument: remote or refspec.
        # `:branch` deletes a ref, `+ref` / `+src:dst` force-updates it.
        if token.startswith((":", "+")):
            return True
    return False


def blocks(command, depth=0):
    if "push" not in command or depth > MAX_DEPTH:
        return False
    tokens = tokenize(command)
    if tokens is None:
        # Unparseable (e.g. unbalanced quotes): fail closed on a coarse match.
        return bool(COARSE_FORCE_PUSH.search(command))
    for segment in split_segments(tokens):
        args = push_args(segment)
        if args is not None and is_force_push(args):
            return True
        if any(blocks(nested, depth + 1) for nested in nested_commands(segment)):
            return True
    return any(blocks(body, depth + 1) for body in substitutions(command))


def main():
    try:
        payload = json.load(sys.stdin)
    except (ValueError, OSError):
        return
    if not blocks(payload.get("tool_input", {}).get("command", "")):
        return
    print(
        json.dumps(
            {
                "hookSpecificOutput": {
                    "hookEventName": "PreToolUse",
                    "permissionDecision": "deny",
                    "permissionDecisionReason": DENY_REASON,
                }
            }
        )
    )


if __name__ == "__main__":
    main()
    sys.exit(0)
