#!/usr/bin/env bash
#
# git clean filter for config/claude/home/.claude/settings.json.
#
# ~/.claude/settings.json is a symlink into this repository, so every tool that
# writes to it writes straight into the working tree. Two of them write things
# that must not reach a public repository:
#
#   - Orca injects its agent hooks on startup: a 3KB generated one-liner per
#     hook event, rewritten whenever the app updates.
#   - Claude Code regenerates an `autoMode` block describing whichever
#     repository was open at the time. It has carried private repository
#     names, GCP project ids, BigQuery dataset names and CI secret names.
#
# This filter removes both on the way into git, so the repository holds only
# the settings it deliberately manages while the working file keeps
# everything. See .gitattributes for the binding and config/git/.config/config
# for the registration.
#
# Output is key-sorted so that a rewrite by one of those tools does not show up
# as a diff just because it reordered the file.
#
# Two things to know:
#
#   - `git checkout`, `git restore` and `git stash` write the filtered version
#     back over the working file, so the Orca hooks are gone until the app adds
#     them again on its next start. Nothing is lost that the app does not
#     rewrite by itself.
#   - A rewrite that only touches the stripped parts still updates the file's
#     mtime, so `git status` can show it as modified while `git diff` is empty.
#     `git add` on the file clears that; it stages identical content.

set -euo pipefail

if ! command -v jq > /dev/null 2>&1; then
  echo "claude-settings-clean: jq is required by the clean filter" >&2
  exit 1
fi

exec jq -S '
  del(.autoMode)
  | if has("hooks") then
      .hooks |= (
        with_entries(
          .value |= (
            map(.hooks |= map(select((.command // "") | test("orca"; "i") | not)))
            | map(select((.hooks | length) > 0))
          )
        )
        | with_entries(select((.value | length) > 0))
      )
    else . end
'
