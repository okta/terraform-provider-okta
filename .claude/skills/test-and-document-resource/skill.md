---
name: test-and-document-resource
description: Use this skill when the user wants to write acceptance tests, generate documentation, and record VCR cassettes for an existing Okta Terraform resource or data source. Trigger on phrases like "write a test for this resource", "generate docs for", "record VCR for", "add test coverage for", or when the user points at a resource Go file and asks for tests or docs.
user_invocable: true
---

# Test and Document a Terraform Resource / Data Source

Full end-to-end workflow: read resource → check constant → write test → create fixture → compile → record VCR → verify playback → write docs.

---

## STEP 1 — Read the resource (or data source) file

Read the Go file the user points at:

```
okta/services/<package>/resource_okta_<name>.go
```

Identify:
- **Package** (`idaas`, `governance`, etc.)
- **Resource type name** — value returned by `Metadata` (e.g. `"okta_identity_source_group"`)
- **Schema attributes**: Required, Optional, Computed fields; nested blocks
- **Parent params**: any path-level IDs the resource needs (e.g. `identity_source_id`)
- **What `Read` actually maps back** — compare the fields set in the API response mapping against the full schema. Fields that `Read` silently drops will cause `ImportStateVerify` failures and need `ImportStateVerifyIgnore`.
- **Provider factory version**: look for `ProtoV6ProviderFactoriesForTestAcc` vs `ProtoV5ProviderFactoriesForTestAcc`

---

## STEP 2 — Check for the resource constant

The test file references `resources.<Const>`. Check that the constant exists:

```bash
grep -n "okta_<name>" \
  /Users/aditya.anand/okta-sdks/terraform-provider-okta/okta/resources/resources.go
```

If the constant is missing, add it in alphabetical order:

```go
OktaIDaaS<PascalName>   = "okta_<name>"
```

The constants file is at `okta/resources/resources.go`. Insert at the correct alphabetical position — do not append at the end.

---

## STEP 3 — Check for existing test / fixture / docs

```bash
ls okta/services/<package>/resource_okta_<name>_test.go 2>/dev/null
ls examples/resources/okta_<name>/resource.tf          2>/dev/null
ls docs/resources/<name>.md                            2>/dev/null
```

Do not overwrite files that already exist. If a test exists, check whether it needs a new step or whether VCR recording is missing.

---

## STEP 4 — Choose real fixture values

For resources that need parent IDs (e.g. `identity_source_id`, `resource_id`) the fixture must use real values that exist in the recording org. Ask the user if not already provided.

Credentials org: `OKTA_ORG_NAME=dcp-tf-testing-2026-01-06`, `OKTA_BASE_URL=oktapreview.com`.
Cassette name is always `oie-00` unless the user specifies otherwise.

---

## STEP 5 — Write the acceptance test

Create `okta/services/<package>/resource_okta_<name>_test.go` (package `<package>_test`).

### Naming
- Test function: `TestAccResourceOkta<PascalName>_basic`
- Resource name variable: `fmt.Sprintf("%s.test", resources.OktaIDaaS<PascalName>)`

### Import mechanics
For resources with composite IDs (parent + child), write a custom `ImportStateIdFunc`:

```go
func importStateIdFor<PascalName>(resourceName string) resource.ImportStateIdFunc {
    return func(s *terraform.State) (string, error) {
        rs, ok := s.RootModule().Resources[resourceName]
        if !ok {
            return "", fmt.Errorf("resource not found: %s", resourceName)
        }
        parentId := rs.Primary.Attributes["<parent_param>"]
        id := rs.Primary.ID
        return fmt.Sprintf("%s/%s", parentId, id), nil
    }
}
```

Use `*terraform.State` from `"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"` — not `resource.State`.

### ImportStateVerifyIgnore
Any field that `Read` does not map back from the API must be listed in `ImportStateVerifyIgnore`. Determine this by reading the `Read` function carefully.

### Provider factory
Use `ProtoV6ProviderFactoriesForTestAcc` for Framework resources (most idaas resources generated after 2024). Use `ProtoV5ProviderFactoriesForTestAcc` for SDK v2 resources.

### CheckDestroy
Copy the pattern from a neighboring resource test in the same package. The destroy check should skip in VCR play mode:

```go
func check<PascalName>Destroy(s *terraform.State) error {
    if os.Getenv("OKTA_VCR_TF_ACC") == "play" {
        return nil
    }
    // ... actual destroy check
    return nil
}
```

### Full test template

```go
package <package>_test

import (
    "fmt"
    "testing"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
    "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
    "github.com/okta/terraform-provider-okta/okta/acctest"
    "github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOkta<PascalName>_basic(t *testing.T) {
    mgr := newFixtureManager("resources", resources.OktaIDaaS<PascalName>, t.Name())
    config := mgr.GetFixtures("resource.tf", t)
    resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaS<PascalName>)

    acctest.OktaResourceTest(t, resource.TestCase{
        PreCheck:                 acctest.AccPreCheck(t),
        ErrorCheck:               testAccErrorChecks(t),
        ProtoV6ProviderFactoriesForTestAcc: acctest.ProtoV6ProviderFactoriesForTestAcc(t),
        Steps: []resource.TestStep{
            {
                Config: config,
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckResourceAttrSet(resourceName, "id"),
                    resource.TestCheckResourceAttr(resourceName, "<parent_param>", "<parent-id>"),
                    resource.TestCheckResourceAttr(resourceName, "<required_field>", "<value>"),
                ),
            },
            {
                ResourceName:            resourceName,
                ImportState:             true,
                ImportStateIdFunc:       importStateIdFor<PascalName>(resourceName),
                ImportStateVerify:       true,
                ImportStateVerifyIgnore: []string{"<fields_missing_from_read>"},
            },
        },
    })
}

func importStateIdFor<PascalName>(resourceName string) resource.ImportStateIdFunc {
    return func(s *terraform.State) (string, error) {
        rs, ok := s.RootModule().Resources[resourceName]
        if !ok {
            return "", fmt.Errorf("resource not found: %s", resourceName)
        }
        parentId := rs.Primary.Attributes["<parent_param>"]
        id := rs.Primary.ID
        return fmt.Sprintf("%s/%s", parentId, id), nil
    }
}
```

> **Note**: If there is no composite ID (top-level resource), use `resource.ImportStatePassthroughID` in the resource's `ImportState` method and omit `ImportStateIdFunc` from the test step.

---

## STEP 6 — Create the fixture .tf file

```
examples/resources/okta_<name>/resource.tf
```

Name the Terraform resource block `"test"` (the test looks up `okta_<name>.test`):

```hcl
resource "okta_<name>" "test" {
  <parent_param>  = "<real-id-from-org>"
  <required_field> = "<value>"

  # optional fields
  <optional_field> = "<value>"
}
```

Use real IDs that exist in `dcp-tf-testing-2026-01-06.oktapreview.com`.

---

## STEP 7 — Compile check

```bash
GOTOOLCHAIN=auto go build ./okta/services/<package>/...
```

Fix any errors before proceeding. Common issues:
- Missing import for `"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"` when using `*terraform.State`
- `resource.State` used instead of `*terraform.State` — always use the latter
- Missing constant in `resources.go`

---

## STEP 8 — Record the VCR cassette

Source credentials, then run in record mode. The cassette name is `oie-00` by default.

**For idaas package:**
```bash
source /Users/aditya.anand/okta_env.sh && \
OKTA_VCR_TF_ACC=record \
OKTA_VCR_CASSETTE=oie-00 \
TF_ACC=1 \
GOTOOLCHAIN=auto \
go test -tags unit -mod=readonly -test.v -timeout 120m \
  -run TestAccResourceOkta<PascalName>_basic \
  ./okta/services/idaas
```

**For governance package:**
```bash
source /Users/aditya.anand/okta_env.sh && \
OKTA_VCR_TF_ACC=record \
OKTA_VCR_CASSETTE=oie-00 \
TF_ACC=1 \
GOTOOLCHAIN=auto \
go test -tags unit -mod=readonly -test.v -timeout 120m \
  -run TestAccResourceOkta<PascalName>_basic \
  ./okta/services/governance
```

The cassette is written to:
- `test/fixtures/vcr/idaas/<TestFuncName>/oie-00.yaml`  (idaas)
- `test/fixtures/vcr/governance/<TestFuncName>/oie-00.yaml`  (governance)

### Common recording failures

| Error | Cause | Fix |
|---|---|---|
| `OKTA_ORG_NAME must be set` | Credentials not in env | `source /Users/aditya.anand/okta_env.sh` |
| `ImportStateVerify failed: profile.*` | `Read` doesn't restore field from API | Add to `ImportStateVerifyIgnore` |
| `Provider produced inconsistent result` | Computed field missing `Computed: true` | Edit resource schema |
| `resource not found after apply` | Parent ID doesn't exist in org | Use a valid ID from the org |
| `404 Not Found` during recording | Resource created but Read gets wrong path | Check parent param wiring in resource |

---

## STEP 9 — Verify playback

```bash
OKTA_VCR_TF_ACC=play \
TF_ACC=1 \
GOTOOLCHAIN=auto \
go test -tags unit -mod=readonly -test.v -timeout 120m \
  -run TestAccResourceOkta<PascalName>_basic \
  ./okta/services/<package>
```

Playback must pass before proceeding. If it fails, re-record after fixing the underlying issue.

---

## STEP 10 — Write documentation

Create `docs/resources/<name>.md`. Populate schema sections by reading the `Schema()` function in the resource file.

````markdown
---
page_title: "okta_<name> Resource - terraform-provider-okta"
description: |-
  Manages an Okta <Human Name> resource.
---

# okta_<name> (Resource)

Manages an Okta <Human Name>.

## Example Usage

```hcl
resource "okta_<name>" "example" {
  <parent_param>   = "<parent-id>"
  <required_field> = "<value>"
}
```

## Schema

### Required

- `<parent_param>` (String) <description from schema>. Forces replacement if changed.
- `<required_field>` (String) <description>

### Optional

- `<optional_field>` (String) <description>
- `<nested_block>` (Block, Optional) <description>
  - `<nested_field>` (String) <description>

### Read-Only

- `id` (String) The unique identifier for the resource.
- `<computed_field>` (String) <description>

## Import

Import using `{<parent_param>}/{id}`:

```shell
terraform import okta_<name>.example <parent-id>/<resource-id>
```

> **Note**: Fields not returned by the API `GET` endpoint will not be restored after import.
> Add them back manually or use `lifecycle { ignore_changes = [...] }`.
````

For a **data source**, create `docs/data-sources/<name>.md` instead.

---

## Key project constants

| Item | Value |
|---|---|
| Credentials file | `/Users/aditya.anand/okta_env.sh` |
| Recording org | `dcp-tf-testing-2026-01-06.oktapreview.com` |
| Default cassette name | `oie-00` |
| Cassette dir (idaas) | `test/fixtures/vcr/idaas/<TestFuncName>/` |
| Cassette dir (governance) | `test/fixtures/vcr/governance/<TestFuncName>/` |
| Fixture dir | `examples/resources/okta_<name>/resource.tf` |
| Docs dir (resource) | `docs/resources/<name>.md` |
| Docs dir (data source) | `docs/data-sources/<name>.md` |
| Resource constants | `okta/resources/resources.go` |
| Build env | `GOTOOLCHAIN=auto` (go.mod requires go >= 1.25) |
| idaas proto version | `ProtoV6ProviderFactoriesForTestAcc` |

---

## Checklist

- [ ] Resource Go file read — schema, Read mapping, parent params all understood
- [ ] Resource constant exists in `okta/resources/resources.go`
- [ ] Existing test/fixture/docs checked (don't overwrite)
- [ ] Real fixture values confirmed with user (parent IDs, required fields)
- [ ] Test file written with correct package, imports, and `*terraform.State`
- [ ] `ImportStateVerifyIgnore` set for any field `Read` doesn't restore
- [ ] Fixture `.tf` file created with resource block named `"test"`
- [ ] Compile check passes (`GOTOOLCHAIN=auto go build`)
- [ ] VCR recording passes
- [ ] VCR playback passes
- [ ] Docs written with accurate schema sections and import instructions
