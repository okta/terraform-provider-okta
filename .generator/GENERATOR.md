# Terraform Provider Code Generator

The Go-based generator at `.generator/go-generator/` reads an OpenAPI 3.x spec and a config
file to produce fully functional Terraform Plugin Framework resource, data source, and test files.

---

## Table of Contents

1. [Architecture Overview](#1-architecture-overview)
2. [How to Run](#2-how-to-run)
3. [What the Generator Reads from the OAS](#3-what-the-generator-reads-from-the-oas)
4. [The Config File (`config_full.yaml`)](#4-the-config-file-config_fullyaml)
5. [Non-Polymorphic APIs](#5-non-polymorphic-apis)
6. [Polymorphic APIs (oneOf + discriminator)](#6-polymorphic-apis-oneof--discriminator)
7. [Nested / Sub-Resources (parent_params)](#7-nested--sub-resources-parent_params)
8. [What Gets Generated](#8-what-gets-generated)
9. [Schema & Field Resolution](#9-schema--field-resolution)
10. [Request Body Field Derivation](#10-request-body-field-derivation)
11. [Body Setter Method Derivation](#11-body-setter-method-derivation)
12. [Write-Only Fields](#12-write-only-fields)
13. [Global Field Exclusions](#13-global-field-exclusions)
14. [date-time Fields](#14-date-time-fields)
15. [Debug Mode](#15-debug-mode)
16. [Post-Processing (goimports)](#16-post-processing-goimports)
17. [Adding a New Resource](#17-adding-a-new-resource)
18. [Known Limitations](#18-known-limitations)
19. [Why Not Magic Modules or tfplugingen-framework?](#19-why-not-magic-modules-or-tfplugingen-framework)

---

## 1. Architecture Overview

```
 ┌──────────────────────────────────────────────────────────────┐
 │  Inputs                                                       │
 │  ① okta-management-APIs-oasv3-noEnums-inheritance.yaml       │
 │     (OpenAPI 3.x spec — 457 paths, 116 tags)                 │
 │  ② config_full.yaml                                          │
 │     (maps resource names → API paths + variants + parents)   │
 └─────────────────────────┬────────────────────────────────────┘
                           │
          ┌────────────────▼───────────────┐
          │  tf-generator (Go binary)       │
          │                                 │
          │  internal/openapi   — spec parse│
          │  internal/config    — YAML load │
          │  internal/generator — orchestr. │
          │  internal/formatter — naming    │
          │  internal/types     — type map  │
          │  templates/*.tmpl   — Go tmpl   │
          └────────────────┬───────────────┘
                           │
 ┌─────────────────────────▼────────────────────────────────────┐
 │  Output (tmp-go-api-calls/ or okta/fwprovider/)              │
 │  resource_okta_{name}_generated.go                           │
 │  resource_okta_{name}_generated_test.go                      │
 │  data_source_okta_{name}_generated.go                        │
 └──────────────────────────────────────────────────────────────┘
```

### Internal packages

| Package | Responsibility |
|---|---|
| `internal/openapi` | Load YAML, resolve `$ref`, walk `allOf`/`oneOf`, extract properties, derive union type names, extract request body schema refs |
| `internal/config` | Unmarshal `config_full.yaml` into Go structs |
| `internal/generator` | Build `TemplateData`, expand variants, render templates, apply field exclusions, merge write-only fields |
| `internal/formatter` | `snake_case` ↔ `CamelCase`, derive SDK method names, sanitize descriptions |
| `internal/types` | Map OAS `type`+`format` → `types.String/Int64/Bool/List`, `schema.*Attribute` |
| `templates/` | `text/template` files for resource, data source, and test |

---

## 2. How to Run

```bash
cd .generator/go-generator

# Build the binary after code changes
go build -o tf-generator ./cmd/...

# Normal generation
./tf-generator \
  -output /path/to/terraform-provider-okta/tmp-go-api-calls \
  -templates ./templates \
  /path/to/okta-management-APIs-oasv3-noEnums-inheritance.yaml \
  ../config_full.yaml

# Full debug output (see every spec lookup, schema resolution, property accepted/skipped)
./tf-generator -debug \
  -output ./tmp-go-api-calls \
  -templates ./templates \
  /path/to/spec.yaml \
  ../config_full.yaml 2>&1 | grep -A 40 "resource: policy_access"
```

> **Flags must come before positional arguments.**

**Flags**

| Flag | Default | Description |
|---|---|---|
| `-output` | required | Directory to write generated files into |
| `-templates` | required | Directory containing `*.tmpl` files |
| `-debug` | false | Emit verbose lines to stderr |

**Post-processing (run after generation):**

```bash
# Remove unused imports and format all generated Go files
find ./tmp-go-api-calls -name "*.go" | xargs ~/go/bin/goimports -w

# Format Terraform HCL example files
make tf-fmt
```

---

## 3. What the Generator Reads from the OAS

### 3.1 `paths` — Required

Every path listed in `config_full.yaml` must exist in the spec's `paths` map.

```yaml
paths:
  /api/v1/policies/{policyId}:
    get:      # needed for Read
    post:     # needed for Create
    put:      # needed for Update
    delete:   # needed for Delete
```

### 3.2 Response schemas — Drive property list

The generator derives Terraform schema attributes from the **response body schema** of the Read
operation (GET). See §9 for the full fallback chain.

### 3.3 Request body schemas — Drive request body fields

The generator reads the **Create op request body schema** to build `RequestBodyFields` and the
**Update op request body schema** to build `UpdateRequestBodyFields`. These may be different
schemas (e.g. `AppServiceAccount` for POST vs `AppServiceAccountForUpdate` for PATCH). See §10.

### 3.4 OAS fields used by the generator

| OAS field | Effect on generated code |
|---|---|
| `properties.<name>.type` | Maps to `types.String / Int64 / Bool / List` |
| `properties.<name>.format: date-time` | Sets `IsDateTime: true` → `.Format(time.RFC3339)` in Read |
| `properties.<name>.readOnly: true` | Sets `Computed: true`; excluded from request bodies |
| `properties.<name>.description` | Populates `Description:` in schema |
| `required: [field]` | Sets `Required: true` on the attribute |
| `allOf: [$ref]` | Followed recursively — base + subclass properties are merged |
| `$ref` on a property | Resolved to get concrete `type` / `properties` |
| `oneOf: [...]` + `discriminator` | Triggers per-variant generation (see §6) |
| `x-codegen-request-body-name` | Overrides the SDK body setter method name (see §11) |
| `operationId` | Drives SDK API method name derivation |

### 3.5 Fields NOT used by the generator

| OAS field | Status |
|---|---|
| `parameters` | Not used — path params come from `parent_params` in config |
| `security` | Not used |
| `x-okta-*` extensions (other than `x-codegen-request-body-name`) | Not used |

---

## 4. The Config File (`config_full.yaml`)

The config maps logical resource names to API operations. It has two top-level keys:
`resources` and `datasources`.

### 4.1 Simple resource (no parent, no variants)

```yaml
resources:
  group:
    api_tag: Group          # → client.GroupAPI in generated code
    read:
      method: get
      path: /api/v1/groups/{groupId}
    create:
      method: post
      path: /api/v1/groups
    update:
      method: put
      path: /api/v1/groups/{groupId}
    delete:
      method: delete
      path: /api/v1/groups/{groupId}
```

- `api_tag` — the Okta SDK client field name, e.g. `GroupAPI`
- `read` / `create` / `update` / `delete` — each optional; omitting `update` or `delete` emits
  a "not supported" stub in the template

### 4.2 Data source

```yaml
datasources:
  group:
    api_tag: Group
    singular:
      method: get
      path: /api/v1/groups/{groupId}      # fetch-by-id
    plural:
      method: get
      path: /api/v1/groups                # list-all
```

### 4.3 Nested / sub-resource (`parent_params`)

See §7 for details.

### 4.4 Polymorphic resource (`variants`)

See §6 for details.

---

## 5. Non-Polymorphic APIs

A non-polymorphic API has a single concrete schema for its response — no `oneOf`.

**What the generator does:**

1. Looks up the Read operation in the spec
2. Resolves `$ref` → reads `properties` + walks `allOf` recursively
3. Builds `[]PropData` — one entry per property (skipping `id`, `links`, `_links`)
4. Merges any request-body-only fields from Create/Update schemas (see §12)
5. Renders a full resource file with `Read`, `Create`, `Update`, `Delete`

**Generated schema (excerpt):**

```go
"created": schema.StringAttribute{
    Description: "Created",
    Computed:    true,   // readOnly: true in OAS
},
"name": schema.StringAttribute{
    Description: "Name",
    Required:    true,
},
```

**Generated Read (excerpt):**

```go
result, httpResp, err := client.GroupAPI.GetGroup(ctx, id).Execute()
// ...
state.Created = types.StringValue(result.GetCreated().Format(time.RFC3339))  // date-time field
state.Name = types.StringValue(result.GetName())
```

---

## 6. Polymorphic APIs (oneOf + discriminator)

Some Okta APIs return different schema shapes depending on a **discriminator field**. For example,
`GET /api/v1/logStreams/{id}` returns `LogStreamAws` or `LogStreamSplunk` depending on `type`.

**Config:**

```yaml
resources:
  log_stream:
    api_tag: StreamingAPI
    variants:
      - suffix: aws
        schema_ref: LogStreamAws
        discriminator_field: type
        discriminator_value: aws_eventbridge
      - suffix: splunk
        schema_ref: LogStreamSplunk
        discriminator_field: type
        discriminator_value: splunk_cloud_logstreaming
    create:
      method: post
      path: /api/v1/logStreams
    read:
      method: get
      path: /api/v1/logStreams/{logStreamId}
    update:
      method: put
      path: /api/v1/logStreams/{logStreamId}
    delete:
      method: delete
      path: /api/v1/logStreams/{logStreamId}
```

**Generated resources:**

```
resource_okta_log_stream_aws_generated.go     → type = "aws_eventbridge"
resource_okta_log_stream_splunk_generated.go  → type = "splunk_cloud_logstreaming"
```

### 6.1 Union type unwrapping

When the list/get operation returns a union wrapper type (e.g. `ListLogStreams200ResponseInner`),
the generator:

1. Auto-derives `UnionTypeName` from the list op's `operationId`
   (`listLogStreams` → `"ListLogStreams200ResponseInner"`)
2. **Read**: unwraps `result.LogStreamAws` before accessing fields
3. **Create/Update**: wraps the body: `okta.LogStreamAwsAsListLogStreams200ResponseInner(body)`

```go
// Generated Read for log_stream_aws:
variantObj := result.LogStreamAws
if variantObj == nil {
    resp.Diagnostics.AddError(...)
    return
}
state.Name = types.StringValue(variantObj.GetName())
```

### 6.2 Variant request body fields

For variants, `RequestBodyFields` is built from the variant schema's props filtered to
`!p.Computed && isScalarGoType`. The `Computed: true` flag always wins — fields marked
`readOnly: true` in OAS are excluded from the body even if also listed in `required`.

### 6.3 Currently configured polymorphic resources

| Base resource | Discriminator field | Variants |
|---|---|---|
| `application` | `signOnMode` | `auto_login`, `basic_auth`, `bookmark`, `browser_plugin`, `oidc`, `saml_11`, `saml`, `secure_password_store`, `ws_federation` |
| `behavior` | `type` | `anomalous_location`, `anomalous_ip`, `anomalous_device`, `velocity`, `anomalous_asn` |
| `log_stream` | `type` | `aws`, `splunk` |
| `network_zone` | `type` | `ip`, `dynamic`, `dynamic_v2` |
| `policy` | `type` | `access`, `idp_discovery`, `mfa_enroll`, `okta_sign_on`, `password`, `profile_enrollment`, `post_auth_session`, `entity_risk`, `device_signal_collection` |

---

## 7. Nested / Sub-Resources (parent_params)

Sub-resources are resources whose API paths contain a parent resource's ID segment, e.g.
`/api/v1/apps/{appId}/federated-claims/{claimId}`.

> **Note:** The generator does NOT read path `parameters` from the OAS. Parent IDs are declared
> explicitly in `config_full.yaml` via `parent_params`.

**Config:**

```yaml
resources:
  application_federated_claim:
    api_tag: ApplicationSSOFederatedClaimsAPI
    parent_params:
      - name: app_id
        description: "ID of the parent application"
        path_param: "{appId}"
    read:
      method: get
      path: /api/v1/apps/{appId}/federated-claims/{claimId}
    create:
      method: post
      path: /api/v1/apps/{appId}/federated-claims
    update:
      method: put
      path: /api/v1/apps/{appId}/federated-claims/{claimId}
    delete:
      method: delete
      path: /api/v1/apps/{appId}/federated-claims/{claimId}
```

**Generated schema (excerpt):**

```go
"app_id": schema.StringAttribute{
    Description: "ID of the parent application",
    Required:    true,
    PlanModifiers: []planmodifier.String{
        stringplanmodifier.RequiresReplace(),
    },
},
```

**Generated API call:**

```go
result, httpResp, err := client.ApplicationSSOFederatedClaimsAPI.
    GetFederatedClaim(ctx, app_id, id).Execute()
```

---

## 8. What Gets Generated

For each **resource** entry, two files are emitted:

### `resource_okta_{name}_generated.go`

```
package fwprovider

type {Name}Resource struct { Config *config.Config }
type {Name}Model struct {
    ID               types.String  // always present
    {ParentParam}    types.String  // one per parent_param (Required + RequiresReplace)
    {Discriminator}  types.String  // variant resources only (Required + RequiresReplace)
    {Property}       types.{X}     // one per OAS schema property (String/Int64/Bool/List)
    {WriteOnlyField} types.String  // request-body-only fields absent from response
}

func New{Name}Resource() resource.Resource
func (r *{Name}Resource) Metadata(...)    // TypeName = "okta_{tf_name}"
func (r *{Name}Resource) Configure(...)  // injects *config.Config
func (r *{Name}Resource) Schema(...)     // full Terraform attribute schema
func (r *{Name}Resource) ImportState(...) // ImportStatePassthroughID
func (r *{Name}Resource) Read(...)       // full SDK call + state mapping
func (r *{Name}Resource) Create(...)     // SDK call + request body construction + computed field mapping
func (r *{Name}Resource) Update(...)     // only if update op is configured
func (r *{Name}Resource) Delete(...)     // only if delete op is configured
```

### `resource_okta_{name}_generated_test.go`

Acceptance test scaffold with `TestAccOkta{Name}_basic`, `testAccOkta{Name}Config_basic`,
`testAccCheckOkta{Name}Exists`, `testAccCheckOkta{Name}Destroy`.

### `data_source_okta_{name}_generated.go`

```
type {Name}DataSource struct { Config *config.Config }
type {Name}DataSourceModel struct {
    ID         types.String  // Optional + Computed
    {Property} types.{X}    // all Computed
}

func (d *{Name}DataSource) Read(...)
  // if state.ID set → singular GET by ID
  // else            → list + filter
```

---

## 9. Schema & Field Resolution

### 9.1 Primary schema (drives Properties / model struct)

The generator tries schema sources in priority order. First non-nil wins:

```
1. GET  read.path    → responses[200/201].content["application/json"].schema
2. POST create.path  → responses[200/201].content["application/json"].schema
3. PUT  update.path  → responses[200/201].content["application/json"].schema
4. POST create.path  → requestBody.content["application/json"].schema
5. PUT  update.path  → requestBody.content["application/json"].schema

→ if all nil: WARNING logged, resource generated with 0 properties
```

### 9.2 Request-body-only field merging

After the primary schema is resolved, the generator **merges in fields from the Create/Update
request body schemas that don't already exist in the response schema**. These are marked
`WriteOnly: true` (see §12).

Example: `EmailDomain` POST body has `brandId`, but `EmailDomainResponseWithEmbedded` (GET
response) does not → `brandId` is merged in so it appears in the schema and Create body, but
`state.BrandId = result.GetBrandId()` is never emitted in Read (getter doesn't exist).

### 9.3 Per-property resolution

For each property the generator:

1. Resolves `$ref` chains recursively
2. Walks `allOf` entries (base class inheritance)
3. Sets `Computed: true` if `readOnly: true` in OAS
4. Sets `Required: true` if listed in the schema's `required` array AND `allComputed` is false
5. Sets `IsDateTime: true` if `type: string` and `format: date-time`
6. Skips `id` (always surfaced as the resource's `ID` field)
7. Skips globally excluded fields (`links`, `_links`)

---

## 10. Request Body Field Derivation

The generator derives **separate** field lists for Create and Update because the two operations
may use entirely different schemas.

### `RequestBodyFields` (Create body)

Derived from the Create op's request body schema:
- Filter: `!p.Computed && isScalarGoType(p.GoType)`
- Excludes readOnly fields and non-scalar types (lists, nested objects)
- Falls back to Update schema if there is no Create op

### `UpdateRequestBodyFields` (Update body)

Derived independently from the Update op's request body schema:
- Same filter: `!p.Computed && isScalarGoType(p.GoType)`
- Falls back to `RequestBodyFields` if Update has no request body schema

**Example — `application_federated_claim`:**

| Op | Schema | Fields sent |
|---|---|---|
| POST (Create) | `FederatedClaimRequestBody` | `expression`, `name` |
| PUT (Update) | `FederatedClaim` | `expression`, `name` (readOnly `id`, `created`, `lastUpdated` excluded) |

**Example — `service_account`:**

| Op | Schema | Fields sent |
|---|---|---|
| POST (Create) | `AppServiceAccount` | `containerOrn`, `description`, `name`, `password`, `username` |
| PATCH (Update) | `AppServiceAccountForUpdate` | `description`, `name` only |

### SDK type names

| Field | Purpose |
|---|---|
| `SDKTypeName` | `okta.New<SDKTypeName>WithDefaults()` for Create body |
| `UpdateSDKTypeName` | `okta.New<UpdateSDKTypeName>WithDefaults()` for Update body |

Both are derived from the request body `$ref` name of each op's schema, falling back to
`TitleName` if no `$ref` is present.

---

## 11. Body Setter Method Derivation

The Okta Go SDK generates a builder method to set the request body. Its name comes from the
request body schema `$ref` name when `x-codegen-request-body-name` is absent.

**Priority (evaluated independently per op):**

```
1. x-codegen-request-body-name (OAS extension)  → GoFieldName(value)
2. Request body $ref name                        → used as-is (already CamelCase)
3. "Body"                                        → default fallback
```

- **`BodySetterMethod`** — derived from the **Create** op (falling back to Update)
- **`UpdateBodySetterMethod`** — derived independently from the **Update** op

**Example — `application_federated_claim`:**

```go
// POST uses FederatedClaimRequestBody schema:
createReq = createReq.FederatedClaimRequestBody(*body)

// PUT uses FederatedClaim schema:
updateReq = updateReq.FederatedClaim(*updateBody)
```

---

## 12. Write-Only Fields

Some OAS schemas include fields in the **request body** that are absent from the **response** —
the API never echoes them back (e.g. `brandId` on `EmailDomain`).

These fields are merged into `Properties` (model struct + TF schema) but marked `WriteOnly: true`.

**Behaviour:**

| Location | `WriteOnly = false` | `WriteOnly = true` |
|---|---|---|
| Model struct | ✅ present | ✅ present |
| TF schema attribute | ✅ emitted | ✅ emitted |
| Create body (`body.SetX(...)`) | if in `RequestBodyFields` | ✅ included |
| Read state mapping (`state.X = result.GetX()`) | ✅ emitted | ❌ **skipped** — getter doesn't exist on response type |

This matches the pattern used by the hand-written `resource_okta_email_domain.go` — `brand_id`
stays in state from the user's plan between refreshes.

---

## 13. Global Field Exclusions

The `excludedTFAttrs` map in `generator.go` lists TF attribute names that are never surfaced,
regardless of which resource they appear on:

| Excluded field | Reason |
|---|---|
| `links` | HAL hypermedia links — purely navigational, meaningless in TF state |
| `_links` | Same |

These are filtered in `schemaToPropsDepth` at parse time, so they never appear in any generated
model, schema, or CRUD method across all 289 resources and data sources.

---

## 14. date-time Fields

OAS properties with `type: string` and `format: date-time` map to `types.String` in Terraform,
but the Okta SDK getter returns `time.Time`. The generator sets `IsDateTime: true` on `PropData`
and the template emits `.Format(time.RFC3339)`:

```go
// Generated Read:
state.Created = types.StringValue(result.GetCreated().Format(time.RFC3339))

// Compare with a regular string field:
state.Name = types.StringValue(result.GetName())
```

This applies in Read, and in post-Create / post-Update computed field mapping.

---

## 15. Debug Mode

Run with `-debug` to see every step on stderr:

```
[DEBUG] Spec paths loaded: 457
[DEBUG] Resources in config: 85

[DEBUG] ─── resource: application_federated_claim ───────────────────
[DEBUG]   APITag               : ApplicationSSOFederatedClaimsAPI
[DEBUG]   HasParent            : true
[DEBUG]   HasUpdate            : true   HasDelete: true
[DEBUG]   SDKTypeName          : FederatedClaimRequestBody
[DEBUG]   UpdateSDKTypeName    : FederatedClaim
[DEBUG]   BodySetterMethod     : FederatedClaimRequestBody
[DEBUG]   UpdateBodySetterMethod: FederatedClaim
[DEBUG]   [schemaToProps]   SKIP "id" (collides with id)
[DEBUG]   [schemaToProps]   SKIP "links" (globally excluded)
[DEBUG]   [schemaToProps]   ACCEPT expression → computed=false readOnly=false
[DEBUG]   [schemaToProps]   ACCEPT created    → computed=true  readOnly=true
[DEBUG]   [buildResourceData] merging request-body-only field "brand_id" (WriteOnly)
[DEBUG]   requestBodyFields=2 updateRequestBodyFields=2
```

**Useful debug filters:**

```bash
# See only warnings
./tf-generator -debug ... 2>&1 | grep WARNING

# Trace one specific resource end-to-end
./tf-generator -debug ... 2>&1 | grep -A 60 "resource: email_domain"

# See which source each resource's schema came from
./tf-generator -debug ... 2>&1 | grep "using schema from\|WARNING: no schema"

# See all merged write-only fields
./tf-generator -debug ... 2>&1 | grep "request-body-only field"

# Count resources with 0 properties
./tf-generator -debug ... 2>&1 | grep "WARNING: no schema found" | wc -l
```

---

## 16. Post-Processing (goimports)

The generator does not run `goimports` automatically. Always run it after regenerating to remove
unused imports (`time`, `net/http`, `okta`, etc.) and format code:

```bash
find ./tmp-go-api-calls -name "*.go" | xargs ~/go/bin/goimports -w

# Install a compatible version (Go 1.24 requires tools <= v0.29)
go install golang.org/x/tools/cmd/goimports@v0.29.0
```

---

## 17. Adding a New Resource

### Step 1 — Verify the OAS has what's needed

```bash
python3 -c "
import yaml
with open('spec.yaml') as f:
    spec = yaml.safe_load(f)
path = '/api/v1/my-resource/{id}'
for method, op in spec['paths'].get(path, {}).items():
    rb = op.get('requestBody',{}).get('content',{}).get('application/json',{}).get('schema',{})
    rs = op.get('responses',{}).get('200',{}).get('content',{}).get('application/json',{}).get('schema',{})
    print(f'{method}: requestBody={rb}, response={rs}')
"
```

### Step 2 — Determine the resource type

| Scenario | Config pattern |
|---|---|
| Single concrete response schema | Simple resource (§5) |
| `oneOf` + `discriminator` in response | Polymorphic with `variants` (§6) |
| Path contains `{parentId}` segment | Nested with `parent_params` (§7) |
| POST and PUT use different body schemas | Auto-derived — no manual config needed |
| Combination of nested + polymorphic | Use both `parent_params` and `variants` |

### Step 3 — Add entry to `config_full.yaml`

Minimum viable entry:

```yaml
resources:
  my_resource:
    api_tag: MyResourceAPI   # must match the Okta SDK client field exactly
    read:
      method: get
      path: /api/v1/my-resources/{resourceId}
    create:
      method: post
      path: /api/v1/my-resources
    update:
      method: put
      path: /api/v1/my-resources/{resourceId}
    delete:
      method: delete
      path: /api/v1/my-resources/{resourceId}
```

Key rules:
- `api_tag` must match the Okta SDK client field name exactly
- Paths must match the OAS exactly including `{camelCase}` path param names
- For nested resources, `parent_params[].name` becomes the Terraform attribute name (snake_case)
- `SDKTypeName`, `UpdateSDKTypeName`, `BodySetterMethod`, and `UpdateBodySetterMethod` are all
  auto-derived from the OAS — no manual annotation needed in most cases

### Step 4 — Regenerate and post-process

```bash
cd .generator/go-generator
./tf-generator \
  -output ../../tmp-go-api-calls \
  -templates ./templates \
  /path/to/spec.yaml \
  ../config_full.yaml

find ../../tmp-go-api-calls -name "*.go" | xargs ~/go/bin/goimports -w
```

### Step 5 — Verify the output

```bash
# Check model struct has expected fields
grep -A 20 "type myResourceModel struct" ../../tmp-go-api-calls/resource_okta_my_resource_generated.go

# Check correct SDK types and setter methods
grep "body\.\|createReq\.\|updateReq\." ../../tmp-go-api-calls/resource_okta_my_resource_generated.go

# Confirm no readOnly fields leaked into request body
grep "body\.Set" ../../tmp-go-api-calls/resource_okta_my_resource_generated.go
```

---

## 18. Known Limitations

| Issue | Impact | Workaround |
|---|---|---|
| `anyOf` | Not supported — treated as unknown | Use `oneOf` equivalent or handle manually |
| `oneOf` without `discriminator.mapping` | 0 properties generated | Add `variants` to config with concrete schema names |
| `additionalProperties` | Not mapped | Parent becomes `types.Object`; implement manually if needed |
| Nested `object` (SingleNestedAttribute) | Schema and model struct generated; body setters for nested fields not generated | Add nested body setters manually |
| `array` items type | `schema.ListAttribute` emitted; SDK setters for list fields not generated | Add list body setters manually |
| Import state for nested resources | `ImportStatePassthroughID` only works for simple IDs | Implement composite ID import manually (e.g. `app_id/claim_id`) |
| `204 No Content` on all ops | 0 properties (no response schema) | Point `read` at a list endpoint |
| Response `$ref` pointing to a response object (not a schema) | Schema not resolved | Add inline schema to spec |
| `x-codegen-request-body-name` absent AND no `$ref` on request body | Setter defaults to `.Body(...)` — may be wrong | Add `x-codegen-request-body-name` to the OAS operation |

---

## 19. Why Not Magic Modules or tfplugingen-framework?

### Magic Modules (GoogleCloudPlatform/magic-modules)

Magic Modules is built exclusively for the **GCP Terraform provider**. Its input is a
**hand-authored `ResourceName.yaml`** per resource — there is no OAS parser. It generates raw
HTTP calls (no typed Go SDK), targets SDKv2 (`helper/schema` not Plugin Framework), and has no
concept of oneOf discriminators, typed request body builders, or per-op schema differences.

For 289 Okta resources it would require 289 hand-maintained YAML files — defeating the purpose
of automation. Our generator re-runs against the OAS and regenerates everything when the spec
changes with a single command.

### Hashicorp `terraform-plugin-codegen-framework`

The HashiCorp Framework Code Generator (tech preview, last release v0.4.1 Sep 2024) generates
only:
- `Schema()` function
- Model struct

**It generates zero CRUD logic.** No `Read()`, `Create()`, `Update()`, `Delete()`. The
documentation explicitly states: *"Over time, it is anticipated that the generator will be further
enhanced to support CRUD logic."*

Our core value is generating **correct API calls** against the typed Okta Go SDK:

| Capability | HashiCorp generator | Our generator |
|---|---|---|
| `Schema()` + model struct | ✅ | ✅ |
| `Read/Create/Update/Delete` methods | ❌ | ✅ |
| `okta.New<Type>WithDefaults()` SDK call | ❌ | ✅ |
| `.FederatedClaimRequestBody(body)` setter | ❌ | ✅ derived from OAS `$ref` |
| POST vs PUT use different body types | ❌ Update op "does not affect mapping" | ✅ `UpdateSDKTypeName` / `UpdateBodySetterMethod` |
| `readOnly` → `Computed`, excluded from body | Partial | ✅ `!p.Computed` filter |
| Request-body-only fields (write-only) | Manual `ignores` per resource | ✅ `WriteOnly: true` auto-detected |
| oneOf / discriminator / union unwrap | ❌ | ✅ full variant expansion + union wrap/unwrap |
| `date-time` → `.Format(time.RFC3339)` | ❌ | ✅ `IsDateTime` flag |
| `links`/`_links` globally excluded | Manual `ignores` per resource | ✅ `excludedTFAttrs` global map |
| Production-ready | ❌ "NOT INTENDED FOR PRODUCTION USE" | ✅ |
