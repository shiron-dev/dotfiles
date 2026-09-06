#!/usr/bin/env bash
#
# Resolve / switch the dotfiles profile for this machine.
#
# The active profile is stored outside the repository so that a machine never
# commits its own identity by accident:
#
#   ~/.config/dotfiles/profile   ->  "work" | "private" | ...
#
# Everything under profiles/<name>/ is overlaid on top of config/ by the
# ansible `config` role, so a profile can either add new files or replace
# common ones.

set -euo pipefail

DOTFILES_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PROFILES_DIR="$DOTFILES_DIR/profiles"
PROFILE_FILE="${DOTFILES_PROFILE_FILE:-$HOME/.config/dotfiles/profile}"
DEFAULT_PROFILE="private"

usage() {
  cat <<EOM
Usage: profile.bash <command>

Commands:
  get           Print the active profile (falls back to "$DEFAULT_PROFILE")
  set <name>    Write <name> to $PROFILE_FILE
  list          List the profiles available in profiles/
  path          Print the path of the active profile directory
  file          Print the path of the profile marker file
EOM
}

list_profiles() {
  [ -d "$PROFILES_DIR" ] || return 0
  find "$PROFILES_DIR" -mindepth 1 -maxdepth 1 -type d -exec basename {} \; | sort
}

get_profile() {
  local profile=""
  if [ -f "$PROFILE_FILE" ]; then
    # First non-empty, non-comment line wins.
    profile="$(grep -v '^[[:space:]]*#' "$PROFILE_FILE" | grep -v '^[[:space:]]*$' | head -n1 | tr -d '[:space:]')"
  fi
  echo "${profile:-$DEFAULT_PROFILE}"
}

set_profile() {
  local name="${1:-}"
  if [ -z "$name" ]; then
    echo "Error: profile name is required" >&2
    usage >&2
    exit 1
  fi
  if [ ! -d "$PROFILES_DIR/$name" ]; then
    echo "Error: unknown profile '$name'. Available:" >&2
    list_profiles | sed 's/^/  - /' >&2
    exit 1
  fi
  mkdir -p "$(dirname "$PROFILE_FILE")"
  printf '%s\n' "$name" > "$PROFILE_FILE"
  echo "[INFO] Profile set to '$name' ($PROFILE_FILE)"
  echo "[INFO] Run the ansible playbook to apply it:"
  echo "       cd $DOTFILES_DIR/scripts/ansible && ansible-playbook -i hosts.yml site.yml"
}

case "${1:-}" in
  get) get_profile ;;
  set)
       shift
              set_profile "${1:-}"
                                   ;;
  list) list_profiles ;;
  path) echo "$PROFILES_DIR/$(get_profile)" ;;
  file) echo "$PROFILE_FILE" ;;
  -h | --help | help | "") usage ;;
  *)
    echo "Error: unknown command '$1'" >&2
    usage >&2
    exit 1
    ;;
esac
