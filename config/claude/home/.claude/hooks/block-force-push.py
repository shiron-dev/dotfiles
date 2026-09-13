#!/usr/bin/env python3
"""Deny only genuine `git push` force/delete invocations.

The previous version regex-matched the whole command string, so unrelated
commands (`git clean -f`, `git checkout -f`, `git branch -d`, `rm -f`, ...)
were blocked too. This version tokenizes the command, isolates each `git push`
invocation, and inspects only that invocation's arguments.
"""
import json
import re
import shlex
import sys

# Wrappers that may precede `git` without changing the meaning of the command.
WRAPPERS = {"sudo", "command", "env", "time", "nohup", "rtk", "proxy"}
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


def push_args(segment):
    """Return the argument list of `git push` for this segment, else None."""
    i = 0
    # Leading env assignments and wrapper commands.
    while i < len(segment) and (ENV_ASSIGN.match(segment[i]) or segment[i] in WRAPPERS):
        i += 1
    if i >= len(segment):
        return None
    if segment[i].rsplit("/", 1)[-1] != "git":
        return None
    i += 1
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


def blocks(command):
    if "push" not in command:
        return False
    try:
        lexer = shlex.shlex(command, posix=True, punctuation_chars=True)
        lexer.whitespace_split = True
        tokens = list(lexer)
    except ValueError:
        # Unparseable (e.g. unbalanced quotes): fail closed on a coarse match.
        return bool(re.search(r"(^|\s)(--force\S*|-f|--mirror|--delete)(\s|=|$)", command))
    return any(
        is_force_push(args)
        for args in (push_args(segment) for segment in split_segments(tokens))
        if args is not None
    )


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
