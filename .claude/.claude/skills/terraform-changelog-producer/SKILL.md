---
name: terraform-changelog-producer
description: Use this skill when the user wants generate a changelog message for a new release.
---

# Terraform Issue Resolver

## When to Use
Use this skill when the user wants generate a changelog message for a new release.

## Dependency Management

The script dependencies are managed with [uv](https://docs.astral.sh/uv/) via the `scripts/pyproject.toml`.

### First-time setup (or after cloning)
```bash
cd .claude/skills/terraform-changelog-producer/scripts
uv sync
```
This creates a `.venv` and installs all pinned dependencies from `uv.lock`.

### Adding a new dependency
```bash
cd .claude/skills/terraform-changelog-producer/scripts
uv add <package>
```

### Upgrading dependencies
```bash
cd .claude/skills/terraform-changelog-producer/scripts
uv lock --upgrade      # regenerate uv.lock with latest versions
uv sync                # apply the updated lock file
```

### Running the script with the managed environment
```bash
cd .claude/skills/terraform-changelog-producer/scripts
uv run python gather_all_commits_since.py <commit_id>
```
`uv run` automatically activates the `.venv` for the invocation without requiring a manual `source .venv/bin/activate`.

---

## Workflow
STEP 1 : Run the script using `uv run` so the managed virtual environment is used automatically. The user must pass in the commit ID:
```bash
cd .claude/skills/terraform-changelog-producer/scripts
uv run python gather_all_commits_since.py <commit_id>
```
STEP 2 : Use the output of the script executed in STEP 1 to generate a changelog message, it should be a list of changes where each change is of the form 
* <CHANGE_MADE> <PR_LINK> by <PR_CREATOR>
where <CHANGE_MADE> can be picked up from the commit message itself. You must use code blocks for anything that is a field or a resource or data source name.

This is an example of what an entry should look like
* New resource and data source `okta_session_violation_policy` and `okta_session_violation_policy_rule` [#2759](https://github.com/okta/terraform-provider-okta/pull/2759) by [pranav-okta](https://github.com/pranav-okta)

The changes should be categorized into 3 main categories, namely "FEATURES", "BUG FIXES" and "ENHANCEMENTS".
When a new resource or data source is being added, it should fall under the category of "FEATURE"
When new fields are being supported for an existing resource or data source, it should fall under the category of "BUG FIXES"
If it doesn't fall under the above 2 categories, it is a "BUG FIX"

STEP 3 : Write the output to a local file named CHANGELOG_<TS>.md where <TS> should be replaced by the current timestamp, which can be obtained by running the command "date -u +%s"
## References
- @./references/notes.md
