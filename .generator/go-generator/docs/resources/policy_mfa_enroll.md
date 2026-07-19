---
page_title: "okta_policy_mfa_enroll Resource - terraform-provider-okta"
description: |-
  Manages an Okta Policy Mfa Enroll resource.
---

# okta_policy_mfa_enroll (Resource)

Manages an Okta Policy Mfa Enroll.

## Example Usage

```hcl
resource "okta_policy_mfa_enroll" "example" {
  type = "<type>"
  name = "Example Name"

  # Optional fields
  # include = "<include>"
  # description = "Example description"
  # priority = 0
  # aaguid_groups = "<aaguid_groups>"
  # type = "<type>"
}
```

## Schema

### Required

- `type` (String) Discriminator field identifying the variant type. Must be set to \
- `name` (String) Name of the policy

### Optional

- `include` (List) Groups to be included
- `description` (String) Description of the policy
- `priority` (Number) Specifies the order in which this policy is evaluated in relation to the other policies
- `aaguid_groups` (List) The list of FIDO2 WebAuthn authenticator groups allowed for enrollment. The authenticators in the group are based on FIDO Alliance Metadata Service that's identified by name or the Authenticator At...
- `type` (String) Grace period type  * `BY_DATE_TIME`: The grace period is defined by a specific date and time. * <x-lifecycle class='ea'></x-lifecycle>`BY_SKIP_COUNT`: The grace period is defined by the number of t...
- `self` (String) Requirements for the user-initiated enrollment
- `key` (String) A label that identifies the authenticator
- `type` (String) Type of policy configuration object  <x-lifecycle class='oie'></x-lifecycle> The `type` property in the policy `settings` is only applicable to the authenticator enrollment policy available in Iden...
- `status` (String) Status
- `system` (Boolean) Specifies whether Okta created the policy

### Read-Only

- `id` (String) The unique identifier for the resource.
- `embedded` (String) Embedded
- `created` (String) Timestamp when the policy was created
- `last_updated` (String) Timestamp when the policy was last modified
