---
name: generate-okta-resource
description: Use this skill when the user provides one or more Okta API paths and wants to generate a Terraform resource and/or data source for them. Also trigger when the user says "regenerate all", "generate all resources", "regenerate everything", or wants to re-run the generator for all existing config entries. Trigger on phrases like "generate resource for", "add resource for this API", "generate data source for", "generate terraform for these endpoints", or when the user pastes API paths like `/api/v1/identity-sources/{identitySourceId}/users`.
user_invocable: true
---

# Generate Okta Terraform Resource / Data Source

Full end-to-end workflow: parse API paths → add to config_full.yaml → build generator → generate ALL files → run goimports on ALL → write examples and docs for ALL.

> **The generator always regenerates files for every entry in config_full.yaml.**
> Whether you're adding new paths or re-running from scratch, Steps 4–8 cover everything — not just newly added entries.

---

## STEP 0 — Determine mode

Two modes:

**A) Add new paths mode** — user provides new API paths → follow Steps 1–8 (add to config, then regenerate everything).

**B) Regenerate-all mode** — user says "regenerate all", "regenerate everything", or similar, with no new paths → skip Steps 1–3 (config already up to date), jump directly to Step 4 (build) and continue through Step 8 for **all** resources and data sources currently in config_full.yaml.

In both modes, Steps 4–8 always operate on the full config — not just the newly added entries.

---

## STEP 1 — Parse the user's API paths

*(Skip this step in regenerate-all mode.)*

The user will provide one or more API paths (e.g. `/api/v1/identity-sources/{identitySourceId}/users`).

For each path, determine:
- **HTTP methods** available (POST=create, GET=read, PUT/PATCH=update, DELETE=delete)
- **Path parameters** (curly-brace tokens like `{identitySourceId}`) — all but the last one are parent params; the last one (if present) is the resource ID param
- **Resource name** (snake_case): derive from the path segments. Drop `/api/v1/`, strip path parameter segments, join remaining segments with `_`. Example: `/api/v1/identity-sources/{identitySourceId}/users` → `identity_source_user` (singular for CRUD resource).
- **API tag**: look up in the OAS spec or ask the user if unclear. Common tags: `IdentitySource`, `Application`, `Policy`, `Group`, `User`.

**Resource vs. Data Source rules:**
- If the path set includes POST + GET + (PUT or PATCH) + DELETE → generate a **resource** (full CRUD).
- If the path set is GET-only (no POST) → generate a **data source** only.
- If the user explicitly requests only a data source for a CRUD path, respect that.
- A **plural GET** (path ends without a path-param `{id}`) → data source with `plural` operation.
- A **singular GET** (path ends with a path-param) → data source with `singular` operation.

**Parent params**: for nested resources, every path parameter except the last `{id}` becomes a `parent_params` entry. Example: `{identitySourceId}` → `name: identity_source_id`, `path_param: "{identitySourceId}"`, `description: "ID of the identity source"`.

---

## STEP 2 — Look up the API tag from the OAS spec

```bash
grep -A2 "operationId" \
  /Users/aditya.anand/okta-sdks/okta-sdk-golang/.generator/okta-management-APIs-oasv3-noEnums-inheritance.yaml \
  | grep -B1 "<path_fragment>"
```

Or search by the last path segment:

```bash
grep "tags:" -A1 \
  /Users/aditya.anand/okta-sdks/okta-sdk-golang/.generator/okta-management-APIs-oasv3-noEnums-inheritance.yaml \
  | grep -B5 "<last_segment>"
```

Use the `tags:` value (first tag listed) as `api_tag`.

---

## STEP 3 — Add entries to config_full.yaml

Config file path:
```
/Users/aditya.anand/okta-sdks/terraform-provider-okta/.generator/config_full.yaml
```

**Before editing**, read the file to find the correct alphabetical insertion point for the new resource/datasource names. The `resources:` and `datasources:` maps must stay alphabetically ordered.

### Resource entry format (full CRUD)

```yaml
  <resource_name>:
    api_tag: <APITag>
    parent_params:
      - name: <parent_param_snake>
        description: "<Human readable description of the parent>"
        path_param: "{<parentParamCamel>}"
    create:
      method: post
      path: /api/v1/<path>
    read:
      method: get
      path: /api/v1/<path>/{id}
    update:
      method: put
      path: /api/v1/<path>/{id}
    delete:
      method: delete
      path: /api/v1/<path>/{id}
```

Omit `update:` if no PUT/PATCH endpoint exists.
Omit `parent_params:` if the resource is top-level (no path params before the ID).

**Special flags** (add only when needed):

| Flag | When to use |
|------|-------------|
| `create_id_field: <tf_attr>` | POST returns no body (204) AND auto-detection will not find a suitable field. The generator now auto-detects this in most cases — only set explicitly if the auto-detect picks the wrong field. |
| `read_noop: true` | No GET endpoint exists for the resource (e.g. bulk-upsert style). |
| `delete_noop: true` | No DELETE endpoint exists for the resource. |

### Data source entry format (singular)

```yaml
  <datasource_name>:
    api_tag: <APITag>
    parent_params:
      - name: <parent_param_snake>
        description: "<Human readable description of the parent>"
        path_param: "{<parentParamCamel>}"
    singular:
      method: get
      path: /api/v1/<path>/{id}
```

### Data source entry format (plural / list)

```yaml
  <datasource_name>:
    api_tag: <APITag>
    parent_params:
      - name: <parent_param_snake>
        description: "<Human readable description of the parent>"
        path_param: "{<parentParamCamel>}"
    plural:
      method: get
      path: /api/v1/<path>
```

---

## STEP 4 — Build the tf-generator binary

```bash
cd /Users/aditya.anand/okta-sdks/terraform-provider-okta/.generator/go-generator
GOTOOLCHAIN=auto go build -o tf-generator ./cmd/...
```

If build fails, check `internal/config/config.go` for missing struct fields — add any new yaml-tagged flags used in config_full.yaml.

---

## STEP 5 — Run the generator

```bash
cd /Users/aditya.anand/okta-sdks/terraform-provider-okta/.generator/go-generator

./tf-generator \
  -output /Users/aditya.anand/okta-sdks/terraform-provider-okta/.generator/go-generator/okta/fwprovider \
  -templates ./templates \
  /Users/aditya.anand/okta-sdks/okta-sdk-golang/.generator/okta-management-APIs-oasv3-noEnums-inheritance.yaml \
  /Users/aditya.anand/okta-sdks/terraform-provider-okta/.generator/config_full.yaml
```

Generated files will appear at:
- `okta/fwprovider/resource_okta_<name>_generated.go`
- `okta/fwprovider/resource_okta_<name>_generated_test.go`
- `okta/fwprovider/data_source_okta_<name>_generated.go`

Verify the expected files were created:
```bash
ls -la /Users/aditya.anand/okta-sdks/terraform-provider-okta/.generator/go-generator/okta/fwprovider/*<name>*
```

---

## STEP 6 — Run goimports on ALL generated files

Always run goimports on every generated file in the output directory — not just newly added ones:

```bash
~/go/bin/goimports -w \
  /Users/aditya.anand/okta-sdks/terraform-provider-okta/.generator/go-generator/okta/fwprovider/resource_okta_*_generated.go \
  /Users/aditya.anand/okta-sdks/terraform-provider-okta/.generator/go-generator/okta/fwprovider/resource_okta_*_generated_test.go \
  /Users/aditya.anand/okta-sdks/terraform-provider-okta/.generator/go-generator/okta/fwprovider/data_source_okta_*_generated.go
```

---

## STEP 7 — Generate example .tf files for ALL resources

For every resource and data source in config_full.yaml, create or overwrite the example file. Check which entries exist in config_full.yaml (read the file), then for each one create the corresponding example if it doesn't already exist (or overwrite if regenerating).

Create example files under `.generator/go-generator/examples/`:

### Resource example

File: `.generator/go-generator/examples/resources/okta_<name>/resource.tf`

```hcl
# Example: okta_<name>

resource "okta_<name>" "example" {
  <parent_param> = "<parent-id>"    # Replace with actual parent ID

  # Required fields — replace with real values
  <required_field_1> = "<value>"
  <required_field_2> = "<value>"

  # Optional fields
  # <optional_field> = "<value>"
}
```

Use the actual schema fields from the generated `.go` file (read it to find Required/Optional/Computed attributes). Use placeholder strings like `"<identity-source-id>"` for IDs, `"user@example.com"` for emails, etc.

### Data source example

File: `.generator/go-generator/examples/data-sources/okta_<name>/data-source.tf`

```hcl
# Example: data.okta_<name>

data "okta_<name>" "example" {
  <parent_param> = "<parent-id>"    # Replace with actual parent ID
  id             = "<resource-id>"  # Replace with actual resource ID
}

output "<name>_external_id" {
  value = data.okta_<name>.example.<key_field>
}
```

---

## STEP 8 — Generate documentation for ALL resources

For every resource and data source in config_full.yaml, create or overwrite the corresponding doc file. Read the generated `.go` file for each entry to populate accurate schema sections.

Create docs under `.generator/go-generator/docs/`:

### Resource doc

File: `.generator/go-generator/docs/resources/<name>.md`

````markdown
---
page_title: "okta_<name> Resource - terraform-provider-okta"
description: |-
  Manages an Okta <HumanName> resource.
---

# okta_<name> (Resource)

Manages an Okta <HumanName>.

## Example Usage

```hcl
resource "okta_<name>" "example" {
  <parent_param> = "<parent-id>"

  <required_field_1> = "<value>"
}
```

## Schema

### Required

- `<parent_param>` (String) <description from schema>
- `<required_field>` (String) <description from schema>

### Optional

- `<optional_field>` (String) <description from schema>
- `<nested_block>` (Block, Optional) <description>
  - `<nested_field>` (String) <description>

### Read-Only

- `id` (String) The unique identifier for the resource.
- `<computed_field>` (String) <description from schema>

## Import

Import using `{<parent_param>}/{id}`:

```shell
terraform import okta_<name>.example <parent-id>/<resource-id>
```
````

### Data source doc

File: `.generator/go-generator/docs/data-sources/<name>.md`

````markdown
---
page_title: "okta_<name> Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta <HumanName>.
---

# okta_<name> (Data Source)

Use this data source to retrieve an Okta <HumanName>.

## Example Usage

```hcl
data "okta_<name>" "example" {
  <parent_param> = "<parent-id>"
  id             = "<resource-id>"
}
```

## Schema

### Required

- `<parent_param>` (String) <description>
- `id` (String) ID of the resource to look up.

### Read-Only

- `<computed_field>` (String) <description from schema>
````

**To populate the schema sections**, read the generated `.go` file and extract attribute names, types, Required/Optional/Computed status, and descriptions from the `Schema()` function.

---

## Key constants and paths

| Item | Value |
|------|-------|
| OAS spec | `/Users/aditya.anand/okta-sdks/okta-sdk-golang/.generator/okta-management-APIs-oasv3-noEnums-inheritance.yaml` |
| Config file | `/Users/aditya.anand/okta-sdks/terraform-provider-okta/.generator/config_full.yaml` |
| Generator binary | `.generator/go-generator/tf-generator` |
| Generator templates | `.generator/go-generator/templates/` |
| Generated output | `.generator/go-generator/okta/fwprovider/` |
| goimports | `~/go/bin/goimports` |
| Build env | `GOTOOLCHAIN=auto` (required — go.mod specifies go 1.25) |
| Examples dir | `.generator/go-generator/examples/` |
| Docs dir | `.generator/go-generator/docs/` |

---

## Naming conventions

| API path pattern | Resource name | Data source name |
|-----------------|---------------|------------------|
| `.../resources/{id}` (top-level) | `resource_name` | `resource_names` (plural) |
| `.../parents/{parentId}/children/{childId}` | `parent_child` | `parent_children` |
| `.../identity-sources/{id}/users/{userId}` | `identity_source_user` | `identity_source_users` |

- Use **singular** for resources (one item managed per TF resource block)
- Use **plural** for list data sources; use **singular** form for get-by-id data sources
- Strip `/api/v1/` prefix and path parameter segments; join remaining path segments with `_`; convert kebab-case segments to snake_case

---

## Checklist

**Adding new paths (mode A):**
- [ ] Resource name derived correctly (singular, snake_case)
- [ ] API tag found in OAS spec
- [ ] Parent params identified and added to config entry
- [ ] config_full.yaml entry inserted in alphabetical order

**Both modes (always):**
- [ ] Generator builds without errors (`GOTOOLCHAIN=auto go build`)
- [ ] Generator runs and ALL `*_generated.go` files refreshed
- [ ] goimports run on ALL `*_generated.go` files (glob, not per-file)
- [ ] Example `.tf` files exist for every resource/datasource in config_full.yaml
- [ ] Doc `.md` files exist for every resource/datasource in config_full.yaml
