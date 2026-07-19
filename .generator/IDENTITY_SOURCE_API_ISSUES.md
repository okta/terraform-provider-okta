# Identity Source API — Terraform Compatibility Issues

Issues encountered while implementing Terraform resources and data sources for the Okta Identity
Source API. Each entry describes the API design problem, its concrete impact on the provider, and
the recommended fix.

---

## Issue 1: Bulk upload endpoints have no GET or DELETE

**Affected endpoints:**
- `POST .../sessions/{sessionId}/bulk-upsert`
- `POST .../sessions/{sessionId}/bulk-delete`
- `POST .../sessions/{sessionId}/bulk-groups-upsert`
- `POST .../sessions/{sessionId}/bulk-groups-delete`
- `POST .../sessions/{sessionId}/bulk-group-memberships-upsert`
- `POST .../sessions/{sessionId}/bulk-group-memberships-delete`

**Problem:**
Terraform's resource model requires four operations: Create, Read, Update, Delete. These endpoints
only provide POST. Without GET, Terraform cannot detect drift (e.g. if the uploaded data was lost
or the session expired). Without DELETE, `terraform destroy` silently does nothing — the staged
data is never cleaned up.

**Impact on provider:**
- `Read` is a no-op that just preserves the last known state unchanged. Any out-of-band change
  is invisible to Terraform.
- `Delete` emits a warning and removes the resource from state without touching Okta.
- These resources behave as write-once actions, not managed resources.

**Workaround applied:**
`read_noop: true` and `delete_noop: true` flags in the generator config. The generated `Read`
preserves state unchanged; `Delete` is a no-op with a user-facing warning.

**Recommended fix:**
Return a stable upload ID from POST and expose GET/DELETE for it:
```
POST .../sessions/{sessionId}/bulk-upsert         → 201 { "id": "bul_abc123", ... }
GET  .../sessions/{sessionId}/bulk-upsert/bul_abc123
DELETE .../sessions/{sessionId}/bulk-upsert/bul_abc123
```

---

## Issue 2: Create returns 202 Accepted with no body — no ID available

**Affected endpoints:**
- All 6 bulk upload endpoints above (202 Accepted, no body)
- `POST .../groups/{groupOrExternalId}/membership` (204 No Content)

**Problem:**
The standard Terraform resource pattern after Create is to call `result.GetId()` from the API
response to set the resource ID. With a 202/204 response there is no body, so `result` does not
exist.

**Impact on provider:**
- Generator cannot emit `plan.ID = types.StringValue(result.GetId())` — compile error.
- Resource ID must be synthesized from an existing plan field (e.g. `session_id`) instead of
  being assigned by the server. This means Terraform can't distinguish two bulk uploads to the
  same session.

**Workaround applied:**
`create_id_field: session_id` generator config key. The generated Create sets
`plan.ID = plan.SessionId` instead of reading from the response.

**Recommended fix:**
Return `201 Created` with the created resource body including a server-assigned `id`:
```yaml
responses:
  '201':
    content:
      application/json:
        schema:
          $ref: '#/components/schemas/BulkUpsertResponse'  # must include `id`
```

---

## Issue 3: Workflow is imperative — 5 chained resources required for one logical sync

**Problem:**
The Identity Source sync workflow is inherently sequential:
1. Create a session
2. POST bulk upsert (users to add/update)
3. POST bulk delete (users to remove)
4. POST bulk groups upsert
5. POST bulk groups delete
6. POST start-import to apply everything

This requires 6 separate Terraform resources with explicit `depends_on` chains. Terraform is
declarative — users expect to describe desired state, not orchestrate a sequence of steps.

**Impact on provider:**
- Users must manage the ordering manually with `depends_on`.
- Any failure mid-sequence leaves orphaned sessions and staged-but-never-imported data in Okta.
- The `okta_identity_source_session_import` resource only makes sense as the last step, but
  nothing in the schema enforces or communicates this.

**Example of the boilerplate required:**
```hcl
resource "okta_identity_source_session" "s" { ... }

resource "okta_identity_source_bulk_upsert" "u" {
  session_id = okta_identity_source_session.s.id
  depends_on = [okta_identity_source_session.s]
}

resource "okta_identity_source_bulk_delete" "d" {
  session_id = okta_identity_source_session.s.id
  depends_on = [okta_identity_source_session.s]
}

resource "okta_identity_source_session_import" "i" {
  session_id = okta_identity_source_session.s.id
  depends_on = [
    okta_identity_source_bulk_upsert.u,
    okta_identity_source_bulk_delete.d,
  ]
}
```

**Recommended fix:**
A single atomic endpoint that accepts all data and triggers import in one call:
```
POST /api/v1/identity-sources/{identitySourceId}/sync
{
  "users":  { "upsert": [...], "delete": [...] },
  "groups": { "upsert": [...], "delete": [...] }
}
```
This collapses 6 Terraform resources into 1 and eliminates session lifecycle management entirely.

---

## Issue 4: Sessions cannot be deleted after import — DELETE returns 400

**Problem:**
Once `start-import` has been triggered on a session, calling `DELETE .../sessions/{sessionId}`
returns `400 Bad Request`. Terraform always tries to destroy resources it created during test
teardown and on `terraform destroy`.

**Impact on provider:**
- `TestAcc*` tests fail during the post-test destroy phase with:
  ```
  Error: Error deleting identity_source_session
  400 Bad Request
  ```
- The provider must special-case 400 responses on session Delete to silently remove from state,
  which masks legitimate bad-request errors.

**Workaround applied:**
Session `Delete` treats both `404 Not Found` and `400 Bad Request` as "already gone" and removes
the resource from state without erroring.

**Recommended fix:**
Either:
- Allow deletion of sessions in any state (DELETE is idempotent in REST convention), or
- Return `409 Conflict` with a clear error code (e.g. `SESSION_ALREADY_IMPORTED`) so the provider
  can distinguish this case from a genuine bad request.

---

## Issue 5: `start-import` is fire-and-forget — no way to check import completion

**Affected endpoint:** `POST .../sessions/{sessionId}/start-import`

**Problem:**
The endpoint triggers an async import and returns immediately with the session object. The
`status` field transitions through states (e.g. `CREATED` → `IN_PROGRESS` → `COMPLETED`), but
there is no webhook, no callback, and no documented polling pattern.

**Impact on provider:**
- After `Create`, `status` is likely still `CREATED` or `IN_PROGRESS`. Terraform reports success
  immediately even though the actual import hasn't completed.
- Users have no way to know from Terraform output whether the import succeeded or failed.
- Import errors surface only when users manually check the Okta Admin Console.

**Recommended fix:**
Either make the endpoint synchronous (wait for completion before responding), or document a
polling pattern and expose a dedicated status endpoint:
```
GET .../sessions/{sessionId}/import-status
→ { "status": "COMPLETED|FAILED|IN_PROGRESS", "errors": [...] }
```
The provider could then poll in the `Create` function until a terminal state is reached.

---

## Issue 6: Separate endpoints for users vs groups — unnecessary resource sprawl

**Problem:**
Users and groups have completely separate bulk operation endpoints despite having identical
operation semantics (upload a list of records to stage for upsert/delete). This results in 4
bulk resources (`bulk_upsert`, `bulk_delete`, `bulk_groups_upsert`, `bulk_groups_delete`) where
1 or 2 would suffice.

**Impact on provider:**
- 10 Terraform resources total (6 resources + 4 data sources) for a single API feature.
- Users must know which resource handles users vs groups — not obvious from names alone.
- More resources to register, test, document, and maintain.

**Recommended fix:**
Unify under a single bulk endpoint with an `entityType` discriminator:
```
POST .../sessions/{sessionId}/bulk-upsert
{ "entityType": "USERS" | "GROUPS", "profiles": [...] }

POST .../sessions/{sessionId}/bulk-delete
{ "entityType": "USERS" | "GROUPS", "externalIds": [...] }
```
Two resources instead of four, with `entity_type` as a required attribute.

---

## Issue 7: Server-set fields missing `readOnly: true`

**Affected field:** `status` on `IdentitySourceSession`

**Problem:**
`status` is set by the server after creation (always `"CREATED"` initially). The OAS spec does
not mark it `readOnly: true`, so the generator emits it as `Optional`. On `terraform apply`,
Terraform sees `status` was `null` in the plan but `"CREATED"` in the response and throws:
```
Error: Provider produced inconsistent result after apply
was null, but now cty.StringVal("CREATED")
```

**Workaround applied:**
Manually added `Computed: true` in the hand-written resource file.

**Recommended fix:**
```yaml
IdentitySourceSession:
  properties:
    status:
      type: string
      readOnly: true   # ← required for all server-set fields
```

---

## Issue 8: Named schemas defined as `type: array` instead of `type: object`

**Affected schemas:**
- `IdentitySourceGroupMembershipsDeleteProfile`
- `IdentitySourceGroupMembershipsUpsertProfile`

**Problem:**
These schemas are defined at the top level as `type: array`. The Go SDK generator creates a
type alias (`type IdentitySourceGroupMembershipsDeleteProfile = []IdentitySourceGroupMembershipsDeleteProfileInner`)
and a separate `Inner` struct for the actual item. There is no
`NewIdentitySourceGroupMembershipsDeleteProfileWithDefaults()` constructor — only the `Inner`
variant exists. Provider code must know this implementation-detail suffix.

**Workaround applied:**
Generator detects named `type: array` schemas and automatically appends `Inner` to derive the
correct constructor name.

**Recommended fix:**
```yaml
# Name the item object, not the array
IdentitySourceGroupMembershipEntry:
  type: object
  properties:
    groupExternalId: { type: string }
    memberExternalIds:
      type: array
      items: { type: string }
```

---

## Issue 9: `profile.profile` double-nesting in group schema

**Affected schema:** `IdentitySourceGroup` → `IdentitySourceGroupProfile` → `IdentitySourceGroupProfileForUpsert`

**Problem:**
The group schema nests a `profile` object that itself contains another field also named `profile`.
The generated Terraform HCL would require:
```hcl
resource "okta_identity_source_group" "example" {
  profile {
    profile {           # confusing double nesting
      display_name = "Engineering"
    }
  }
}
```
The model struct path becomes `plan.Profile.Profile.DisplayName` — error-prone and unintuitive.

**Workaround applied:**
Hand-written resource flattens to a single `profile { display_name, description }` level.

**Recommended fix:**
```yaml
IdentitySourceGroup:
  properties:
    profile:
      type: object
      properties:
        displayName: { type: string }
        description: { type: string }
```

---

## Issue 10: Group membership Read returns a list, not a single item

**Affected endpoint:** `GET .../groups/{groupOrExternalId}/membership`

**Problem:**
The endpoint returns `{ "memberExternalIds": ["id1", "id2", ...] }` — a flat list of all member
IDs for the group. There is no `GET .../membership/{memberId}` endpoint to fetch a single member.
The standard provider Read pattern appends `/{id}` to fetch a specific resource, which does not
work here.

**Workaround applied:**
`read_list_id_field: member_external_ids` generator config key. The generated Read calls the list
endpoint and scans for the target ID, removing the resource from state if absent.

**Recommended fix:**
Add a single-item GET:
```
GET .../groups/{groupOrExternalId}/membership/{memberExternalId}
→ { "memberExternalId": "...", "groupExternalId": "..." }
```

---

## Issue 11: Composite IDs required throughout due to non-globally-unique session IDs

**Problem:**
Sessions, bulk operations, and memberships can only be looked up with both the
`identitySourceId` and the `sessionId`. Session IDs are not globally unique — they are scoped to
an identity source. This forces every resource to use composite import IDs
(`identity_source_id/session_id`) and custom `ImportState` implementations.

**Impact on provider:**
- Every `terraform import` command requires the user to know and provide the composite key format.
- Custom `ImportState` functions required for every resource instead of the default
  `resource.ImportStatePassthroughID`.
- Documentation must explain the composite format for each resource.

**Recommended fix:**
Make session (and child resource) IDs globally unique so `GET /sessions/{sessionId}` works
without the parent path parameter. Standard REST practice: server-assigned IDs should be opaque
and globally unique.

---

## Issue 12: Inline anonymous object schemas inside array items

**Affected schemas:**
- `BulkGroupUpsertRequestBody.profiles` items
- `BulkUpsertRequestBody.profiles` items

**Problem:**
Array item schemas are defined as anonymous inline objects rather than named `$ref` schemas.
The Go SDK synthesizes the name `BulkGroupUpsertRequestBodyProfilesInner` for the item type.
This name is fragile — it changes if the parent field is renamed — and is not obvious to consumers.

**Workaround applied:**
Generator recurses into the inline schema and derives a local model name from the parent field.

**Recommended fix:**
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
        $ref: '#/components/schemas/BulkGroupUpsertItem'   # named ref
```

---

## Summary

| # | API area | Issue | Provider impact | Severity |
|---|---|---|---|---|
| 1 | Bulk upload endpoints | No GET or DELETE | Write-only resources, no drift detection | High |
| 2 | Bulk upload + membership create | 202/204 — no ID in response | ID must be synthesized from plan | High |
| 3 | Overall workflow | Imperative 6-step sequence | 6 chained resources, manual `depends_on` | High |
| 4 | Session delete | 400 after import — cannot destroy | Special-cased 400 handling in Delete | High |
| 5 | `start-import` | Async, no completion signal | Terraform reports success before import finishes | High |
| 6 | Users vs groups | Separate endpoints for identical semantics | 4 bulk resources instead of 2 | Medium |
| 7 | `IdentitySourceSession.status` | Missing `readOnly: true` | "Provider produced inconsistent result" error | High |
| 8 | `IdentitySourceGroupMemberships*Profile` | Named schema is `type: array` | Leaky `Inner` suffix in SDK | Medium |
| 9 | `IdentitySourceGroup.profile` | `profile.profile` double nesting | Confusing HCL, hand-written flatten required | Medium |
| 10 | Group membership read | List endpoint only, no single-item GET | Custom scan-and-search Read required | Medium |
| 11 | All resources | Non-globally-unique IDs | Composite import IDs, custom ImportState everywhere | Medium |
| 12 | Bulk upsert request bodies | Anonymous inline array item schemas | Fragile synthesized SDK type names | Low |
