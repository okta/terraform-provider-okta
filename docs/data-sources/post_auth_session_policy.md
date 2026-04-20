---
page_title: "Data Source: okta_post_auth_session_policy"
description: |-
  Retrieves the Post Auth Session Policy.
---

# Data Source: okta_post_auth_session_policy

Retrieves the Post Auth Session Policy. This is a system policy that is automatically created when Identity Threat Protection (ITP) with Okta AI is enabled. There is exactly one Post Auth Session Policy per organization.

## Example Usage

```terraform
data "okta_post_auth_session_policy" "example" {
}
```

## Attributes Reference

- `id` - ID of the Post Auth Session Policy.
- `name` - Name of the Post Auth Session Policy.
- `status` - Status of the policy: `ACTIVE` or `INACTIVE`.
- `rule_id` - ID of the modifiable policy rule (non-default). Use this for importing the policy rule resource.
