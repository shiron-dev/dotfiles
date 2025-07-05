# PR Scope Check CI

This repository contains a GitHub Actions workflow that checks if the scope in a Pull Request (PR) title matches the files changed in the PR.

## Purpose

The primary goal of this CI is to enforce conventional commit message guidelines, specifically ensuring that the scope defined in the PR title accurately reflects the parts of the codebase being modified. This helps in:

-   Improving the clarity and traceability of changes.
-   Automating changelog generation.
-   Making it easier to understand the impact of a PR at a glance.

## How it Works

The CI performs the following steps:

1.  **Extracts PR Title Scope**: It parses the PR title to find the scope(s) defined within the parentheses, e.g., `type(scope1,scope2): message`.
    -   Multiple scopes can be comma-separated.
    -   A wildcard scope `(*)` is also permitted, which will allow changes in any scope.
2.  **Determines Required Scopes from Changed Files**: It analyzes the paths of the files modified in the PR.
3.  **Matches Scopes against Rules**: It uses a set of rules defined in `scope-rules.json` to map changed file paths to their corresponding scopes.
    -   `config/{any_dir}/...` maps to the `{any_dir}` as scope.
    -   `scripts/{any_dir}/...` maps to the `{any_dir}` as scope.
    -   `data/{any_dir}/...` maps to the `{any_dir}` as scope.
    -   `.github/...` maps to `github` scope.
    -   Other files map to `dotfiles` scope.
4.  **Validates PR Title Scopes**: It checks if the scopes extracted from the PR title cover all the scopes determined from the changed files.
    -   If the PR title uses a wildcard `*`, the check automatically passes.
    -   If there's a mismatch (e.g., a required scope is not listed in the PR title, or the PR title has no scope when one is required), the CI job will fail, and a comment will be posted on the PR detailing the issue.

## Configuration

### Scope Rules (`scope-rules.json`)

The mapping between file paths and their scopes is defined in the `scope-rules.json` file in the root of the repository.

The file contains an array of rule objects, each with a `pattern` (regex) and a `scope` (string). The `scope` can use capture groups from the regex (e.g., `$1`).

**Example `scope-rules.json`:**

```json
{
  "rules": [
    {
      "pattern": "^config/(.+)/.*",
      "scope": "$1"
    },
    {
      "pattern": "^scripts/(.+)/.*",
      "scope": "$1"
    },
    {
      "pattern": "^data/(.+)/.*",
      "scope": "$1"
    },
    {
      "pattern": "^\\.github/.*",
      "scope": "github"
    },
    {
      "pattern": ".*",
      "scope": "dotfiles"
    }
  ]
}
```

Rules are evaluated in the order they are defined. The first rule that matches a file path determines its scope. The last rule `(.*)` acts as a catch-all for files not matching any preceding rules.

### GitHub Actions Workflow

The CI is implemented as a GitHub Actions workflow defined in `.github/workflows/check_pr_scope.yml`. It triggers on `pull_request` events (`opened`, `synchronize`, `reopened`, `edited`).

No special setup is required beyond having this workflow file and the `scope-rules.json` in the repository.

## Troubleshooting

-   **CI Failure with PR Comment**: If the CI fails, it will post a comment on your PR explaining the reason.
    -   **"PR title does not contain a valid scope."**: Ensure your PR title follows the format `type(scope): message`. If no specific scope applies, but changes are made, consider if `dotfiles` or another general scope is appropriate, or use `type(*): message`. The comment will provide examples.
    -   **"PR title scopes do not cover all changed files. Missing scopes for: [scope_list]"**: Your PR title is missing one or more scopes that correspond to the files you've changed. Add the missing scopes to your PR title (e.g., `fix(scopeA,scopeB): resolve issues`). The comment will list the missing scopes and suggest how to update the title.
-   **"No matching scope rule for file: [filepath]"**: This is a warning output in the CI logs (not a PR comment). While it doesn't fail the CI, it indicates that a changed file didn't match any specific rule in `scope-rules.json` and wasn't assigned a scope. This might mean a new rule needs to be added to `scope-rules.json` or the default "dotfiles" scope (if it's the last rule) is being applied implicitly. If the file should have a specific scope, update the rules.

This provides a good overview for users of the repository.
