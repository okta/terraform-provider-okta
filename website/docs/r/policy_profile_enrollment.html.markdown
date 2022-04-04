---
layout: 'okta'
page_title: 'Okta: okta_policy_profile_enrollment'
sidebar_current: 'docs-okta-resource-policy-profile-enrollment'
description: |-
  Creates a Profile Enrollment Policy
---

# okta_policy_profile_enrollment

~> **WARNING:** This feature is only available as a part of the Identity Engine. [Contact support](mailto:dev-inquiries@okta.com) for further information.

This resource allows you to create and configure a Profile Enrollment Policy.

## Example Usage

```hcl
resource "okta_policy_profile_enrollment" "example" {
  name            = "example"
  status          = "ACTIVE"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) Policy Name.

- `status` - (Optional) Status of the policy.

## Attributes Reference

- `id` - ID of the Policy.

## Import

A Profile Enrollment Policy can be imported via the Okta ID.

```
$ terraform import okta_policy_profile_enrollment.example <policy id>
```
