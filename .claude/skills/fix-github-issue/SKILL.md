---
name: fix-github-issue
description: Use this skill when the user provides a GitHub issue URL and wants to replicate the bug, fix it, add a test, record a VCR cassette, push the fix to a new branch, and open a PR. Trigger on phrases like "fix this issue", "fix github issue", "look at this issue and fix it", or when a GitHub URL is given alongside a request to investigate or resolve it.
user_invocable: true
---

# Fix GitHub Issue

Full end-to-end workflow: read issue → replicate → branch → fix → test → VCR record → push → open PR.

---

## STEP 1 — Read the issue

Use the `gh` CLI to fetch the issue. Accept any GitHub issue URL format.

```bash
gh issue view <URL> --json number,title,body,comments
```

Parse out:
- Issue number (`<N>`)
- Title (used for branch name slug)
- Body + comments (the bug description, steps to reproduce, any config snippets)

---

## STEP 2 — Identify affected resource and package

From the issue body, determine:
- Which Terraform resource or data source is affected (e.g. `okta_request_condition`)
- Which Go package it lives in (`okta/services/governance` vs `okta/services/idaas`)
- The relevant source file (e.g. `okta/services/governance/resource_request_condition.go`)

Read the affected source file(s) to understand the current implementation before writing any fix.

---

## STEP 3 — Build a replication TF config

Construct a minimal Terraform config that exercises the reported bug. Use values from the issue body when provided.

**Ask the user** for any values that are environment-specific and not in the issue, for example:
- `resource_id` — the Okta resource/app ID
- `approval_sequence_id` — governance sequence ID
- Group/user IDs referenced in policy blocks
- Any other opaque IDs the issue author likely redacted

Do not guess IDs. Pause and ask before proceeding if required values are missing.

Once you have all values, show the config to the user and confirm it represents the reported scenario.

---

## STEP 4 — Create a branch from master

First, check the current branch and stash any uncommitted changes if needed:

```bash
# Check current branch
git branch --show-current

# If there are uncommitted changes, stash them
git stash list  # check if stash is needed
git diff --stat HEAD  # see what's dirty
```

If the working tree is dirty (modified or untracked files that would be lost), stash before switching:

```bash
git stash push -m "wip: stashed before fix/issue-<N>"
```

Then fetch and create the new branch from master:

```bash
git fetch origin master
git checkout -b fix/issue-<N>-<short-slug> origin/master
```

Where `<short-slug>` is a 3–5 word kebab-case summary of the issue title.
Example: `fix/issue-1234-request-condition-status-null`

> If you stashed changes, remind the user at the end: "Your previous changes were stashed as `wip: stashed before fix/issue-<N>`. Run `git stash pop` on your original branch to restore them."

---

## STEP 5 — Apply the fix

Read and understand the root cause from the issue and the source code before editing.
Make the minimal targeted change that resolves the bug. Do not refactor unrelated code.

---

## STEP 6 — Add an acceptance test

### 6a. Create fixture TF file(s)

Fixture files live at:
```
examples/resources/<resource_name>/<fixture_name>.tf
```

Name the fixture after the issue: e.g. `issue_<N>.tf` or `basic_issue<N>.tf`.

Use the replication config from STEP 3. Replace environment-specific IDs with literal values (the same ones the user provided — cassette recording captures them, playback replays them).

### 6b. Write the test function

Add a new test function to the existing `*_test.go` file for the resource:

```
okta/services/<package>/resource_<name>_test.go
```

Follow the existing pattern exactly:

```go
func TestAcc<ResourcePascalCase>_Issue<N>(t *testing.T) {
    mgr := newFixtureManager("resources", resources.<ResourceConst>, t.Name())
    config := mgr.GetFixtures("issue_<N>.tf", t)
    resourceName := fmt.Sprintf("%s.test", resources.<ResourceConst>)

    acctest.OktaResourceTest(t, resource.TestCase{
        PreCheck:                 acctest.AccPreCheck(t),
        ErrorCheck:               testAccErrorChecks(t),
        ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
        CheckDestroy:             check<ResourcePascalCase>Destroy,
        Steps: []resource.TestStep{
            {
                Config: config,
                Check: resource.ComposeTestCheckFunc(
                    // assertions that verify the bug is fixed
                ),
            },
        },
    })
}
```

The `CheckDestroy` function should skip in VCR play mode:
```go
if os.Getenv("OKTA_VCR_TF_ACC") == "play" {
    return nil
}
```

### 6c. Verify the test compiles

```bash
go build ./okta/services/<package>/...
```

---

## STEP 7 — Record the VCR cassette

The cassette name is the Okta org subdomain you are recording against (e.g. `dev-12345678`).
Ask the user for the cassette name (their org subdomain) if not already known.

Set required environment variables and run the test in record mode:

**For governance resources:**
```bash
OKTA_VCR_TF_ACC=record \
OKTA_VCR_CASSETTE=<cassette-name> \
TF_ACC=1 \
go test -tags unit -mod=readonly -test.v -timeout 120m \
  -run TestAcc<ResourcePascalCase>_Issue<N> \
  ./okta/services/governance
```

**For idaas resources:**
```bash
OKTA_VCR_TF_ACC=record \
OKTA_VCR_CASSETTE=<cassette-name> \
TF_ACC=1 \
go test -tags unit -mod=readonly -test.v -timeout 120m \
  -run TestAcc<ResourcePascalCase>_Issue<N> \
  ./okta/services/idaas
```

The cassette YAML is written to:
- `test/fixtures/vcr/governance/<TestFuncName>/<cassette-name>.yaml`  (governance)
- `test/fixtures/vcr/idaas/<TestFuncName>/<cassette-name>.yaml`       (idaas)

After recording, verify playback passes:

**Governance:**
```bash
OKTA_VCR_TF_ACC=play TF_ACC=1 \
go test -tags unit -mod=readonly -test.v -timeout 120m \
  -run TestAcc<ResourcePascalCase>_Issue<N> \
  ./okta/services/governance
```

**Idaas:**
```bash
OKTA_VCR_TF_ACC=play TF_ACC=1 \
go test -tags unit -mod=readonly -test.v -timeout 120m \
  -run TestAcc<ResourcePascalCase>_Issue<N> \
  ./okta/services/idaas
```

---

## STEP 8 — Push the branch

Stage all changed and new files explicitly (never `git add .`):

```bash
git add okta/services/<package>/resource_<name>.go
git add okta/services/<package>/resource_<name>_test.go
git add examples/resources/<resource_name>/issue_<N>.tf
git add test/fixtures/vcr/<package>/<TestFuncName>/<cassette-name>.yaml
```

Commit and push:
```bash
git commit -m "fix(<resource_name>): <one-line description of the fix> (#<N>)"
git push origin fix/issue-<N>-<short-slug>
```

---

## STEP 9 — Open a pull request

```bash
gh pr create \
  --title "fix(okta_<resource>): <concise one-line description> (#<N>)" \
  --body "$(cat <<'EOF'
## Summary
- <one sentence root cause>
- <one sentence fix applied>
- Fixes #<N>

## Test plan
- [ ] `TestAcc<ResourcePascalCase>_Issue<N>` — acceptance test covering the reported scenario
- [ ] VCR cassette recorded and playback verified

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)" \
  --base master
```

**PR title rules:**
- Under 72 characters
- Format: `fix(<resource>): <what changed> (#<N>)`
- Example: `fix(okta_request_condition): add RequiresReplace to resource_id (#2780)`

**PR body rules:**
- Summary: 2–3 bullet points max — root cause + fix, nothing else
- Test plan: checkbox list of the new test(s) added
- Always close the issue with `Fixes #<N>` in the body so GitHub auto-links it

After creating the PR, return the PR URL to the user.

---

## Key project conventions

| Thing | Pattern |
|---|---|
| Test function name | `TestAcc<ResourcePascalCase>_Issue<N>` |
| Fixture file | `examples/resources/<resource>/<fixture>.tf` |
| Cassette dir (governance) | `test/fixtures/vcr/governance/<TestFuncName>/` |
| Cassette dir (idaas) | `test/fixtures/vcr/idaas/<TestFuncName>/` |
| VCR record env | `OKTA_VCR_TF_ACC=record OKTA_VCR_CASSETTE=<org-subdomain>` |
| VCR play env | `OKTA_VCR_TF_ACC=play` |
| Branch name | `fix/issue-<N>-<short-slug>` |
| Commit message | `fix(<resource>): <description> (#<N>)` |
| Random seed in VCR mode | FNV hash of test name (handled automatically by `newFixtureManager`) |
