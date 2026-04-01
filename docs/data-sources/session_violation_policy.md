---
page_title: "Data Source: okta_session_violation_policy"
description: |-
  Retrieves the Session Violation Detection Policy.
---

# Data Source: okta_session_violation_policy

Retrieves the Session Violation Detection Policy. This is a system policy that is automatically created when the Session Violation Detection feature is enabled. There is exactly one Session Violation Detection Policy per organization.

## Example Usage

```terraform
data "okta_session_violation_policy" "example" {
}
```

## Attributes Reference

- `id` - ID of the Session Violation Detection Policy.
- `name` - Name of the Session Violation Detection Policy.
- `status` - Status of the policy: `ACTIVE` or `INACTIVE`.
- `rule_id` - ID of the modifiable policy rule (non-default). Use this for importing the policy rule resource.
