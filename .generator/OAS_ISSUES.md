# OAS Spec Issues for Terraform Provider Generation

Issues discovered while building the Terraform code generator against the Okta Management API OAS.
Each entry documents the spec problem, its impact on code generation, and the workaround applied.

---

## Part 1: Design Guidelines for Terraform-Compatible APIs

These are the rules to follow **when designing new APIs** that need to be managed as Terraform
resources. Violating them creates friction ranging from workarounds in the generator to impossible
auto-generation that requires hand-maintained code forever.

---

### G1: Every resource must have a full CRUD surface

Terraform's resource model is built on four operations. Each missing operation degrades the user
experience and forces generator workarounds.

| HTTP operation | Terraform method | If missing… |
|---|---|---|
| `POST` | `Create` | Resource cannot be created |
| `GET /{id}` | `Read` | Drift cannot be detected; state goes stale silently |
| `PUT` or `PATCH /{id}` | `Update` | Every config change forces destroy + recreate |
| `DELETE /{id}` | `Delete` | Resource leaks in Okta after `terraform destroy` |

**Rule:** Every API resource that will be a Terraform resource **must** expose `GET`, `POST`,
`PUT`/`PATCH`, and `DELETE`. Fire-and-forget POST-only endpoints (no GET/DELETE) cannot be
modelled as proper Terraform resources.

---

### G2: Create must return the resource with its server-assigned ID

```yaml
# ❌ Bad — 202 Accepted with no body; Terraform can't know the new resource's ID
POST /sessions/{sessionId}/bulk-upsert:
  responses:
    '202':
      description: Accepted

# ✅ Good — 201 Created with the resource body
POST /resources:
  responses:
    '201':
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/MyResource'   # must include `id`
```

**Rule:** Create (`POST`) **must** return `201 Created` with the full resource body including the
server-assigned `id`. Without this, Terraform cannot populate `plan.ID` and subsequent Read/Update/
Delete calls have no ID to work with.

---

### G3: Read must return a single resource by ID, not a list

```yaml
# ❌ Bad — returns a list; Terraform can't diff a single resource
GET /groups/{groupId}/membership:
  responses:
    '200':
      schema:
        type: object
        properties:
          memberExternalIds:
            type: array
            items: { type: string }

# ✅ Good — single resource GET by ID
GET /groups/{groupId}/membership/{memberId}:
  responses:
    '200':
      schema:
        $ref: '#/components/schemas/GroupMembership'
```

**Rule:** Read (`GET`) must fetch **a single resource by its ID** in the path. List endpoints
(`GET /resources`) are for data sources, not resource reads. If you only have a list endpoint,
the provider must scan the whole list on every refresh — expensive and fragile.

---

### G4: Mark all server-set fields as `readOnly: true`

```yaml
# ❌ Bad — status is set by the server but spec says nothing
IdentitySourceSession:
  properties:
    status:
      type: string        # Terraform generates Optional → phantom diffs on every apply

# ✅ Good
IdentitySourceSession:
  properties:
    status:
      type: string
      readOnly: true      # Terraform generates Computed → no diffs
    createdAt:
      type: string
      format: date-time
      readOnly: true
    id:
      type: string
      readOnly: true
```

**Rule:** Any field the server populates (IDs, timestamps, status, computed values) **must** have
`readOnly: true`. Without it, generators emit these as `Optional` and every apply produces a
"Provider produced inconsistent result" error because the plan had `null` but the server returned
a value.

Fields that the user sets on Create but that can't be changed afterwards should also be annotated:
```yaml
name:
  type: string
  x-immutable: true   # signals RequiresReplace to the TF generator
```

---

### G5: Never define a named schema as `type: array`

```yaml
# ❌ Bad — named array schema; SDK codegen creates an "Inner" companion type
IdentitySourceGroupMembershipsDeleteProfile:
  type: array
  items:
    type: object
    properties:
      groupExternalId: { type: string }
      memberExternalIds:
        type: array
        items: { type: string }

# ✅ Good — named object for the item; array only at field level
IdentitySourceGroupMembershipEntry:
  type: object
  properties:
    groupExternalId: { type: string }
    memberExternalIds:
      type: array
      items: { type: string }

BulkGroupMembershipsDeleteRequestBody:
  properties:
    memberships:
      type: array
      items:
        $ref: '#/components/schemas/IdentitySourceGroupMembershipEntry'
```

**Rule:** Named schemas in `components/schemas` must always be `type: object`. Arrays belong at
field level only. A named array schema forces SDK generators to create a `<Name>Inner` type for
items, meaning no `New<Name>WithDefaults()` constructor exists and consumers must know the
implementation-detail `Inner` suffix.

---

### G6: Always use `$ref` for reusable nested objects — avoid inline anonymous objects

```yaml
# ❌ Bad — anonymous inline object inside array; SDK generates ugly synthesized name
BulkGroupUpsertRequestBody:
  properties:
    profiles:
      type: array
      items:
        type: object                  # no name → SDK calls it BulkGroupUpsertRequestBodyProfilesInner
        properties:
          externalId: { type: string }
          profile:
            $ref: '#/components/schemas/IdentitySourceGroupProfileForUpsert'

# ✅ Good — named item schema
BulkGroupUpsertItem:
  type: object
  properties:
    externalId: { type: string }
    profile:
      $ref: '#/components/schemas/IdentitySourceGroupProfileForUpsert'

BulkGroupUpsertRequestBody:
  properties:
    profiles:
      type: array
      items:
        $ref: '#/components/schemas/BulkGroupUpsertItem'
```

**Rule:** Any object used as array items or as a nested field that will appear in Terraform
config must be a named `$ref`. Anonymous inline objects produce synthesized, fragile SDK type
names that break if the parent schema is renamed.

---

### G7: Keep nesting shallow — maximum 2 levels deep

```yaml
# ❌ Bad — 3 levels; TF config becomes profile { profile { display_name = "..." } }
IdentitySourceGroup:
  properties:
    profile:
      $ref: '#/components/schemas/IdentitySourceGroupProfile'

IdentitySourceGroupProfile:
  properties:
    profile:                          # second "profile" — same name, different level
      $ref: '#/components/schemas/IdentitySourceGroupProfileForUpsert'

IdentitySourceGroupProfileForUpsert:
  properties:
    displayName: { type: string }
    description: { type: string }

# ✅ Good — 1 level; TF config is profile { display_name = "..." }
IdentitySourceGroup:
  properties:
    profile:
      type: object
      properties:
        displayName: { type: string }
        description: { type: string }
```

**Rule:** Terraform schema attributes nested more than 2 levels deep are confusing to users and
error-prone to generate. Flatten where semantically reasonable. Never reuse the same field name
at consecutive nesting levels (e.g. `profile.profile`).

---

### G8: Use symmetric request/response schemas with `readOnly` annotations

```yaml
# ❌ Bad — POST and GET use entirely different schema types; generator needs two code paths
POST /groups    body:     GroupWrite      (no id, no createdAt)
GET  /groups/{id} response: Group         (has id, createdAt — completely separate schema)

# ✅ Good — single schema, server-set fields annotated readOnly
Group:
  properties:
    id:
      type: string
      readOnly: true
    createdAt:
      type: string
      format: date-time
      readOnly: true
    name:
      type: string          # user-settable
    description:
      type: string          # user-settable
```

**Rule:** Use **one schema** for both request and response bodies. Mark server-set fields
`readOnly: true` — generators know to exclude them from the request body automatically. Separate
`Write` / `Read` schemas double the maintenance burden and require generators to track two SDK
type names per resource.

---

### G9: Use consistent `id` field naming across all resources

```yaml
# ❌ Bad — different field names mean different generator logic per resource
IdentitySourceGroup:   externalId   # user-provided external ID is the primary key
IdentitySourceSession: id           # server-assigned ID
IdentitySourceUser:    externalId   # same as group, but different from session

# ✅ Good — always use "id" as the primary resource identifier
IdentitySourceGroup:
  properties:
    id:
      type: string
      readOnly: true        # server-assigned primary key, always "id"
    externalId:
      type: string          # business key — separate, clearly named
```

**Rule:** The Terraform provider convention is that `id` is the primary resource identifier used
in `GET /{id}`, `PUT /{id}`, `DELETE /{id}`. If your business key is user-provided (like
`externalId`), still expose a server-assigned `id` for the Terraform resource identity. Using
`externalId` as the Terraform ID requires custom `create_id_field` config and breaks the standard
generator path.

---

### G10: Add `x-codegen-request-body-name` when the body setter name is ambiguous

```yaml
# ❌ Bad — generator guesses the setter method name from schema $ref, may be wrong
post:
  requestBody:
    content:
      application/json:
        schema:
          $ref: '#/components/schemas/FederatedClaimRequestBody'

# ✅ Good — explicitly name the setter method
post:
  x-codegen-request-body-name: FederatedClaimRequestBody
  requestBody:
    content:
      application/json:
        schema:
          $ref: '#/components/schemas/FederatedClaimRequestBody'
```

**Rule:** Always set `x-codegen-request-body-name` on `POST` and `PUT` operations. Without it,
the Go SDK generator derives the setter method name from heuristics that can produce incorrect
results (e.g. `.Body(...)` instead of `.FederatedClaimRequestBody(...)`), causing compile errors
in generated provider code.

---

### G11: Provide a stable, non-composite primary key

```yaml
# ❌ Bad — composite natural key as the only identifier
GET /identity-sources/{identitySourceId}/groups/{groupOrExternalId}/membership/{memberExternalId}
# The resource ID in Terraform would need to be "identitySourceId/groupExternalId/memberExternalId"

# ✅ Good — server assigns a unique ID usable in isolation
GET /memberships/{membershipId}
# Simple ID, or composite expressed as a single opaque token
```

**Rule:** Terraform import (`terraform import`) works best with a single opaque ID. Composite
IDs (e.g. `parentId/childId`) require custom `ImportState` implementations and documentation
explaining the format. If a composite key is unavoidable, document the separator and order
explicitly in the OAS description field.

---

## Part 2: Known Issues in the Current Okta Management API Spec

---

### Issue 1: Named schemas defined as `type: array` instead of `type: object`

**Affected schemas:**
- `IdentitySourceGroupMembershipsDeleteProfile`
- `IdentitySourceGroupMembershipsUpsertProfile`

**Spec as-is:**
```yaml
IdentitySourceGroupMembershipsDeleteProfile:
  type: array
  items:
    type: object
    properties:
      groupExternalId: { type: string }
      memberExternalIds:
        type: array
        items: { type: string }
```

**Impact:**
The Go SDK creates `IdentitySourceGroupMembershipsDeleteProfileInner` as the item struct. No
`NewIdentitySourceGroupMembershipsDeleteProfileWithDefaults()` constructor exists. Consumer code
must know the `Inner` suffix — a leaky SDK implementation detail.

**Workaround applied:**
Generator detects when `OriginalRef` points to a named `type: array` schema and automatically
appends `Inner` to derive the correct SDK item constructor name.

**Violates guideline:** G5

---

### Issue 2: Server-set fields missing `readOnly: true`

**Affected fields:**
- `status` on `IdentitySourceSession`

**Spec as-is:**
```yaml
IdentitySourceSession:
  properties:
    status:
      type: string     # missing readOnly: true
```

**Impact:**
Generator emits `Optional: true`. On `terraform apply`, Terraform detects `status` was `null`
in the plan but is `"CREATED"` in post-apply state:
```
Error: Provider produced inconsistent result after apply
was null, but now cty.StringVal("CREATED")
```

**Workaround applied:**
Manually added `Computed: true` in the hand-written resource file.

**Violates guideline:** G4

---

### Issue 3: Bulk upload endpoints have no GET or DELETE

**Affected endpoints:**
- `POST .../sessions/{sessionId}/bulk-delete`
- `POST .../sessions/{sessionId}/bulk-group-memberships-delete`
- `POST .../sessions/{sessionId}/bulk-group-memberships-upsert`
- `POST .../sessions/{sessionId}/bulk-groups-delete`
- `POST .../sessions/{sessionId}/bulk-groups-upsert`
- `POST .../sessions/{sessionId}/bulk-upsert`

**Impact:**
No GET means drift cannot be detected. No DELETE means `terraform destroy` silently does nothing.
Generator would emit a Read calling a non-existent SDK method — compile error.

**Workaround applied:**
Added `read_noop: true` and `delete_noop: true` config keys. Generated `Read` preserves state
unchanged; generated `Delete` is a no-op.

**Violates guideline:** G1

---

### Issue 4: Create returns 202/204 — no ID in response body

**Affected endpoints:**
- All 6 bulk upload endpoints (202 Accepted, no body)
- `POST .../groups/{groupOrExternalId}/membership` (204 No Content)

**Impact:**
Standard generated Create calls `result.GetId()`. With no response body, `result` doesn't exist
and the call is `_, err := req.Execute()`, causing a compile error if `result.GetId()` is emitted.

**Workaround applied:**
Added `create_id_field: <plan_field>` config key. Generated Create sets `plan.ID = plan.<Field>`
from an existing plan field (e.g. `session_id`) instead of `result.GetId()`.

**Violates guideline:** G2

---

### Issue 5: Inline anonymous object schemas inside array items

**Affected schemas:**
- `BulkGroupUpsertRequestBody.profiles` items
- `BulkUpsertRequestBody.profiles` items

**Impact:**
SDK generator synthesizes the name `BulkGroupUpsertRequestBodyProfilesInner`. Fragile — changes
if the parent field or schema is renamed.

**Workaround applied:**
Generator handles anonymous object items by recursing into the inline schema and deriving a
local model name from the parent field name.

**Violates guideline:** G6

---

### Issue 6: `profile.profile` double-nesting in `IdentitySourceGroup`

**Affected schema:** `IdentitySourceGroup` → `IdentitySourceGroupProfile` → `IdentitySourceGroupProfileForUpsert`

**Impact:**
Generated Terraform HCL would require:
```hcl
profile {
  profile {            # confusing double nesting
    display_name = "..."
  }
}
```
Model struct path becomes `plan.Profile.Profile.DisplayName` — error-prone.

**Workaround applied:**
Hand-written resource flattens to single `profile { display_name, description }`.

**Violates guideline:** G7

---

### Issue 7: Group membership Read endpoint returns a list, not a single item

**Affected endpoint:** `GET .../groups/{groupOrExternalId}/membership`

**Impact:**
Standard Read pattern calls `GET .../membership/{id}` for a single item. This endpoint only
returns a list of all member IDs. Generator would emit a call appending `/{id}` to the list
endpoint — a 404 on every refresh.

**Workaround applied:**
Added `read_list_id_field: member_external_ids` config key. Generated Read iterates the list and
removes the resource from state if the ID is absent.

**Violates guideline:** G3

---

## Summary

| # | Schema / Endpoint | Guideline violated | Severity | Fixed in generator? |
|---|---|---|---|---|
| 1 | `IdentitySourceGroupMembershipsDeleteProfile` | G5: Named schema as array | Medium | ✅ Auto-detects `Inner` suffix |
| 2 | `IdentitySourceSession.status` | G4: Missing `readOnly` | High | ⚠️ Manual fix in hand file |
| 3 | `bulk-*` endpoints | G1: Missing GET + DELETE | High | ✅ `read_noop` / `delete_noop` |
| 4 | `bulk-*` + membership create | G2: No ID in Create response | High | ✅ `create_id_field` |
| 5 | `BulkGroupUpsertRequestBody` | G6: Anonymous inline array items | Low | ✅ Generator handles inline items |
| 6 | `IdentitySourceGroup.profile` | G7: Deep nesting (`profile.profile`) | Medium | ⚠️ Manual flatten in hand file |
| 7 | Group membership read | G3: List endpoint, no single-item GET | Medium | ✅ `read_list_id_field` |


Issues discovered while building the Terraform code generator against the Okta Management API OAS.
Each entry documents the spec problem, its impact on code generation, and the workaround applied.

---

## Issue 1: Named schemas defined as `type: array` instead of `type: object`

**Affected schemas:**
- `IdentitySourceGroupMembershipsDeleteProfile`
- `IdentitySourceGroupMembershipsUpsertProfile`

**Spec as-is:**
```yaml
IdentitySourceGroupMembershipsDeleteProfile:
  type: array          # ← top-level type is array, not object
  items:
    type: object
    properties:
      groupExternalId: { type: string }
      memberExternalIds:
        type: array
        items: { type: string }
```

**Impact:**
The Go SDK generator (openapi-generator) creates two types:
- `IdentitySourceGroupMembershipsDeleteProfile` = `[]IdentitySourceGroupMembershipsDeleteProfileInner`
- `IdentitySourceGroupMembershipsDeleteProfileInner` = the actual item struct

There is **no `NewIdentitySourceGroupMembershipsDeleteProfileWithDefaults()` constructor** — only
`NewIdentitySourceGroupMembershipsDeleteProfileInnerWithDefaults()`. Any consumer building a body
must know about the `Inner` suffix, which is a leaky SDK implementation detail.

**Workaround applied:**
The Terraform generator detects when `OriginalRef` points to a named schema with `type: array`
and automatically appends `Inner` to derive the correct SDK item constructor name.

**Recommended fix:**
```yaml
# Define the item as a named object schema
IdentitySourceGroupMembershipEntry:
  type: object
  properties:
    groupExternalId: { type: string }
    memberExternalIds:
      type: array
      items: { type: string }

# Reference it as array items in the request body
BulkGroupMembershipsDeleteRequestBody:
  type: object
  properties:
    memberships:
      type: array
      items:
        $ref: '#/components/schemas/IdentitySourceGroupMembershipEntry'
```

---

## Issue 2: Server-set fields missing `readOnly: true`

**Affected fields:**
- `status` on `IdentitySourceSession`

**Spec as-is:**
```yaml
IdentitySourceSession:
  properties:
    status:
      type: string     # ← no readOnly: true
    id:
      type: string
      readOnly: true
```

**Impact:**
The generator emits `status` as `Optional: true` (user-configurable). But `status` is always set
by the server after creation — the user never provides it. On `terraform apply`, Terraform detects
`status` was `null` in the plan but is `"CREATED"` in the post-apply state, producing:

```
Error: Provider produced inconsistent result after apply
was null, but now cty.StringVal("CREATED")
```

**Workaround applied:**
Manually added `Computed: true` alongside `Optional: true` in the hand-written resource file
(`resource_okta_identity_source_session.go`).

**Recommended fix:**
```yaml
status:
  type: string
  readOnly: true    # ← mark all server-set fields
```

---

## Issue 3: Bulk upload endpoints have no corresponding GET or DELETE

**Affected endpoints:**
- `POST .../sessions/{sessionId}/bulk-delete`
- `POST .../sessions/{sessionId}/bulk-group-memberships-delete`
- `POST .../sessions/{sessionId}/bulk-group-memberships-upsert`
- `POST .../sessions/{sessionId}/bulk-groups-delete`
- `POST .../sessions/{sessionId}/bulk-groups-upsert`
- `POST .../sessions/{sessionId}/bulk-upsert`

**Spec as-is:**
```yaml
/api/v1/identity-sources/{identitySourceId}/sessions/{sessionId}/bulk-upsert:
  post:
    responses:
      '202':
        description: Accepted    # ← no response body
# No GET, no DELETE defined for this path
```

**Impact:**
Terraform resources require a `Read` method to detect drift and a `Delete` to clean up.
Without a GET endpoint, the generator would emit a `Read` that calls a non-existent
`GetIdentitySourceBulkUpsert` SDK method, causing a compile error.

**Workaround applied:**
Added `read_noop: true` and `delete_noop: true` config keys. When set, the generated `Read`
preserves existing state unchanged (no API call), and `Delete` is a no-op. These resources
behave as fire-and-forget actions wrapped in Terraform resource lifecycle.

**Recommended fix:**
Either:
1. Add a GET endpoint that returns the last submitted payload for the session, or
2. Add an OAS extension to signal the action-style semantics:
```yaml
post:
  x-terraform-resource-type: action   # no Read/Delete needed
```

---

## Issue 4: Create returns 204 No Content — no ID in response body

**Affected endpoints:**
- `POST .../sessions/{sessionId}/bulk-*` (all 6 bulk endpoints)
- `POST .../groups/{groupOrExternalId}/membership`

**Spec as-is:**
```yaml
responses:
  '202':
    description: Accepted    # ← no body, no id field
```

**Impact:**
The standard generated Create calls `result.GetId()` to set the Terraform resource ID. With a
204/202 response there is no `result` — the call signature is `_, err := req.Execute()` — so
`result.GetId()` causes a compile error.

**Workaround applied:**
Added `create_id_field: <plan_field>` config key. When set, the generated Create does not
attempt `result.GetId()` and instead sets `plan.ID = plan.<FieldName>` from a field already
present in the plan (e.g. `session_id`).

**Recommended fix:**
Return a `201 Created` response with the created resource body including its `id`:
```yaml
responses:
  '201':
    content:
      application/json:
        schema:
          $ref: '#/components/schemas/BulkUpsertResponse'
```

---

## Issue 5: Inline anonymous object schemas inside array items

**Affected schemas:**
- `BulkGroupUpsertRequestBody.profiles` items
- `BulkUpsertRequestBody.profiles` items

**Spec as-is:**
```yaml
BulkGroupUpsertRequestBody:
  properties:
    profiles:
      type: array
      items:
        type: object          # ← anonymous inline object, no $ref name
        properties:
          externalId: { type: string }
          profile:
            $ref: '#/components/schemas/IdentitySourceGroupProfileForUpsert'
```

**Impact:**
The SDK generator synthesizes an ugly type name for the anonymous object:
`BulkGroupUpsertRequestBodyProfilesInner`. This name is fragile — it changes if the field is
renamed or the schema is restructured. Consumer code must know this synthesized name.

**Recommended fix:**
Extract the inline object as a named schema:
```yaml
BulkGroupUpsertItem:
  type: object
  properties:
    externalId: { type: string }
    profile:
      $ref: '#/components/schemas/IdentitySourceGroupProfileForUpsert'

BulkGroupUpsertRequestBody:
  properties:
    profiles:
      type: array
      items:
        $ref: '#/components/schemas/BulkGroupUpsertItem'
```

---

## Issue 6: `profile.profile` double-nesting in `IdentitySourceGroup`

**Affected schema:** `IdentitySourceGroup`

**Spec as-is:**
```yaml
IdentitySourceGroup:
  properties:
    profile:
      $ref: '#/components/schemas/IdentitySourceGroupProfile'

IdentitySourceGroupProfile:
  properties:
    profile:                         # ← second level named "profile" again
      $ref: '#/components/schemas/IdentitySourceGroupProfileForUpsert'

IdentitySourceGroupProfileForUpsert:
  properties:
    description: { type: string }
    displayName: { type: string }
```

**Impact:**
The generated Terraform HCL config would require triple nesting:
```hcl
resource "okta_identity_source_group" "example" {
  profile {
    profile {               # ← confusing double "profile"
      display_name = "..."
    }
  }
}
```

The generator resolves the outermost `$ref` correctly but the schema depth produces a
`profile.profile.display_name` path in the model struct, which is unintuitive and error-prone.

**Workaround applied:**
The hand-written resource file flattens the schema to a single `profile { display_name, description }`
level, bypassing the double nesting entirely.

**Recommended fix:**
Flatten the nesting in the spec:
```yaml
IdentitySourceGroup:
  properties:
    profile:
      type: object
      properties:
        description: { type: string }
        displayName: { type: string }
```

---

## Issue 7: Membership read endpoint returns a flat list, not a single item

**Affected endpoint:** `GET /api/v1/identity-sources/{identitySourceId}/groups/{groupOrExternalId}/membership`

**Spec as-is:**
```yaml
responses:
  '200':
    content:
      application/json:
        schema:
          type: object
          properties:
            memberExternalIds:
              type: array
              items: { type: string }   # ← returns all member IDs, not a single one
```

**Impact:**
The standard Read pattern calls `GET .../membership/{memberId}` to fetch a specific member.
But the API only provides a list of all member IDs for the group — there is no single-item GET.
The generator would emit a read call that appends `/{id}` to a list endpoint, which doesn't exist.

**Workaround applied:**
Added `read_list_id_field: member_external_ids` config key. The generated `Read` calls the list
endpoint, iterates `memberExternalIds`, and removes the resource from state if the ID is not found.

**Recommended fix:**
Add a single-item GET endpoint:
```yaml
/api/v1/identity-sources/{identitySourceId}/groups/{groupOrExternalId}/membership/{memberExternalId}:
  get:
    summary: Get a specific group membership
    responses:
      '200':
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/IdentitySourceGroupMembership'
```

---

## Summary Table

| # | Schema / Endpoint | Issue | Severity | Fixed in generator? |
|---|---|---|---|---|
| 1 | `IdentitySourceGroupMembershipsDeleteProfile` | Named schema is `type: array` not `type: object` | Medium | ✅ Auto-detects `Inner` suffix |
| 2 | `IdentitySourceSession.status` | Missing `readOnly: true` | High | ⚠️ Manual fix in hand file |
| 3 | `bulk-*` endpoints | No GET or DELETE | High | ✅ `read_noop` / `delete_noop` config keys |
| 4 | `bulk-*` + membership create | `202`/`204` — no ID in response | High | ✅ `create_id_field` config key |
| 5 | `BulkGroupUpsertRequestBody` | Inline anonymous array item objects | Low | ✅ Generator handles anonymous items |
| 6 | `IdentitySourceGroup.profile` | Double `profile.profile` nesting | Medium | ⚠️ Manual flatten in hand file |
| 7 | Group membership read | List endpoint, no single-item GET | Medium | ✅ `read_list_id_field` config key |
