#!/usr/bin/env bash
set -euo pipefail

pr_title_raw="$1"
changed_files="$2"

# Sanitize pr_title_raw by removing CR/LF characters that might come from jq output
pr_title=$(echo "$pr_title_raw" | tr -d '\r\n')

if [ -z "$pr_title" ]; then
  error_message="PR title is empty. Please provide a title in the format 'type(scope): message'."
  echo "::error::$error_message"
  comment_body="**Scope Check Failed!** 🚨\n\n${error_message}\n\nExample: \`feat(my-scope): add new feature\`."
  echo "$comment_body" > comment_file.txt
  {
    echo "comment_body<<EOF_CMT"
    cat comment_file.txt
    echo "EOF_CMT"
  } >> "$GITHUB_ENV"
  exit 1
fi

# Extract scope from PR title (e.g., "fix(scope1,scope2): ...")
pr_scopes_str=$(echo "$pr_title" | grep -oP '^\w+\(\K[^\)]+' || echo "")
IFS=',' read -r -a pr_scopes <<< "$pr_scopes_str"
# Trim whitespace from scopes
mapfile -t pr_scopes < <(for scope in "${pr_scopes[@]}"; do echo "$scope" | xargs; done)

if [ ${#pr_scopes[@]} -eq 0 ] && [[ "$pr_scopes_str" != "*" ]]; then
  error_message="PR title does not contain a valid scope. Please use the format 'type(scope): message' or 'type(*): message'."
  echo "::error::$error_message"
  comment_body="**Scope Check Failed!** 🚨\n\n${error_message}\n\nTitle received: \`$pr_title\`\n\nExample: \`feat(my-scope): add new feature\` or \`fix(*): resolve an issue\`."
  echo "$comment_body" > comment_file.txt
  {
    echo "comment_body<<EOF_CMT"
    cat comment_file.txt
    echo "EOF_CMT"
  } >> "$GITHUB_ENV"
  exit 1
fi

echo "PR Scopes: ${pr_scopes[*]}"

# Load scope rules
if [ ! -f "scope-rules.json" ]; then
  echo "::error::scope-rules.json not found."
  exit 1
fi
rules=$(jq -r '.rules' < scope-rules.json)

required_scopes=()
while IFS= read -r file; do
  matched_scope=""
  for i in $(seq 0 $(($(echo "$rules" | jq length) - 1))); do
    pattern=$(echo "$rules" | jq -r ".[$i].pattern")
    scope_template=$(echo "$rules" | jq -r ".[$i].scope")

    if [[ "$file" =~ $pattern ]]; then
      # Handle scope template with capture groups (e.g., $1)
      if [[ "$scope_template" == "\$1" ]]; then
        matched_scope="${BASH_REMATCH[1]}"
      else
        matched_scope="$scope_template"
      fi
      break
    fi
  done

  if [ -n "$matched_scope" ]; then
    # Add to required_scopes if not already present
    if [[ ! " ${required_scopes[*]} " =~  ${matched_scope}  ]]; then
      required_scopes+=("$matched_scope")
    fi
  else
    # If no rule matches, consider it an error or a default scope based on requirements
    # For now, let's assume if a file doesn't match any rule, it's an issue.
    # Or, define a default scope in scope-rules.json like { "pattern": ".*", "scope": "default" }
    echo "::warning::No matching scope rule for file: $file"
  fi
done <<< "$changed_files"

echo "Required Scopes based on changed files: ${required_scopes[*]}"

# Output required_scopes for GitHub Actions output
if [ -n "${GITHUB_OUTPUT:-}" ]; then
  echo "required_scopes=${required_scopes[*]}" >> "$GITHUB_OUTPUT"
  echo "pr_scopes=${pr_scopes[*]}" >> "$GITHUB_OUTPUT"
fi

if [[ ${pr_scopes[*]} =~ \* || ${pr_scopes[*]} =~ deps ]]; then
  echo "Wildcard scope '*' or 'deps' in PR title allows all changes."
  exit 0
fi

missing_scopes=()
for req_scope in "${required_scopes[@]}"; do
  is_covered=false
  for pr_scope in "${pr_scopes[@]}"; do
    if [ "$req_scope" == "$pr_scope" ]; then
      is_covered=true
      break
    fi
  done
  if [ "$is_covered" == false ]; then
    if [[ ! " ${pr_scopes[*]} " =~  $req_scope  ]]; then
      missing_scopes+=("$req_scope")
    fi
  fi
done

if [ ${#missing_scopes[@]} -gt 0 ]; then
  echo "::error::PR title scopes do not cover all changed files. Missing scopes for: ${missing_scopes[*]}"
  # Prepare message for PR comment
  comment_body="**Scope Check Failed!** 🚨\n\nPR title scopes do not cover all changed files.\n\n"
  if [ ${#missing_scopes[@]} -gt 0 ]; then
    comment_body+="Missing required scopes in PR title: \`${missing_scopes[*]}\`\n\n" # Use [*] for space separated, or loop for comma separated

    suggestion_scopes=""
    if [ -n "$pr_scopes_str" ] && [[ "$pr_scopes_str" != "*" ]]; then
      suggestion_scopes="${pr_scopes_str}"
    fi

    # Build comma-separated list for suggestion_scopes
    for missing_scope in "${missing_scopes[@]}"; do
      if [ -n "$suggestion_scopes" ]; then
        suggestion_scopes+=",";
      fi
      suggestion_scopes+="$missing_scope"
    done

    comment_body+="Please update your PR title to include these scopes. For example: \`type($suggestion_scopes): your message\` or \`type(*): your message\` if a wildcard is appropriate.\n\n"
  fi
  comment_body+="<details><summary>Details</summary>\n"
  # Ensure pr_scopes are displayed correctly, even if empty or just wildcard
  pr_scopes_display="${pr_scopes[*]}"
  if [ -z "$pr_scopes_display" ] && [ "$pr_scopes_str" == "*" ]; then
    pr_scopes_display="*"
  elif [ -z "$pr_scopes_display" ]; then
    pr_scopes_display="(none)"
  fi
  comment_body+="PR Title Scopes: \`${pr_scopes_display}\`\n"
  comment_body+="Required Scopes from changed files: \`${required_scopes[*]}\`\n"
  comment_body+="Changed Files:\n"
  comment_body+="\`\`\`\n${changed_files}\n\`\`\`\n"
  comment_body+="</details>"

  echo "$comment_body" > comment_file.txt
  {
    echo "comment_body<<EOF_CMT"
    cat comment_file.txt
    echo "EOF_CMT"
  } >> "$GITHUB_ENV"
  exit 1
else
  echo "PR title scopes are valid."
fi
