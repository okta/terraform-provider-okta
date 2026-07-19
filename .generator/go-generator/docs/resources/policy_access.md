---
page_title: "okta_policy_access Resource - terraform-provider-okta"
description: |-
  Manages an Okta Policy Access resource.
---

# okta_policy_access (Resource)

Manages an Okta Policy Access.

## Example Usage

```hcl
resource "okta_policy_access" "example" {
  type = "<type>"
  name = "Example Name"

  # Optional fields
  # conditions = "<conditions>"
  # description = "Example description"
  # priority = 0
  # status = "ACTIVE"
  # system = true
}
```

## Schema

### Required

- `type` (String) Discriminator field identifying the variant type. Must be set to \
- `name` (String) Name of the policy

### Optional

- `conditions` (String) Policy conditions aren't supported. Conditions are applied at the rule level for this policy type.
- `description` (String) Description of the policy
- `priority` (Number) Specifies the order in which this policy is evaluated in relation to the other policies
- `status` (String) Status
- `system` (Boolean) Specifies whether Okta created the policy

### Read-Only

- `id` (String) The unique identifier for the resource.
- `embedded` (String) Embedded
- `created` (String) Timestamp when the policy was created
- `last_updated` (String) Timestamp when the policy was last modified
