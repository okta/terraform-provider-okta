---
layout: 'okta'
page_title: 'Okta: okta_policy_rule_profile_enrollment'
sidebar_current: 'docs-okta-resource-policy-rule-profile-enrollment'
description: |-
  Creates a Profile Enrollment Policy Rule.
---

# okta_policy_rule_profile_enrollment

~> **WARNING:** This feature is only available as a part of the Identity Engine. [Contact support](mailto:dev-inquiries@okta.com) for further information.

This resource allows you to create and configure a Profile Enrollment Policy Rule.

## Example Usage

```hcl
resource "okta_policy_profile_enrollment" "example" {
  name = "My Enrollment Policy"
}

resource "okta_inline_hook" "example" {
  name    = "My Inline Hook"
  status  = "ACTIVE"
  type    = "com.okta.user.pre-registration"
  version = "1.0.3"

  channel = {
    type    = "HTTP"
    version = "1.0.0"
    uri     = "https://example.com/test2"
    method  = "POST"
  }
}

resource "okta_group" "example" {
  name        = "My Group"
  description = "Group of some users"
}

resource "okta_policy_rule_profile_enrollment" "example" {
  policy_id           = okta_policy_profile_enrollment.example.id
  inline_hook_id      = okta_inline_hook.example.id
  target_group_id     = okta_group.example.id
  unknown_user_action = "REGISTER"
  email_verification  = true
  access              = "ALLOW"
  profile_attributes {
    name     = "email"
    label    = "Email"
    required = true
  }
  profile_attributes {
    name     = "name"
    label    = "Name"
    required = true
  }
  profile_attributes {
    name     = "t-shirt"
    label    = "T-Shirt Size"
    required = false
  }
}
```

## Argument Reference

The following arguments are supported:

- `policy_id` - (Required) Policy ID.

- `inline_hook_id` - (Optional) ID of a Registration Inline Hook.

- `target_group_id` - (Optional) The ID of a Group that this User should be added to.

- `unknown_user_action` - (Required) Which action should be taken if this User is new. Valid values are: `"DENY"`, `"REGISTER"`.

- `email_verification` - (Optional) Indicates whether email verification should occur before access is granted. Default is `true`.

- `access` - (Optional) Allow or deny access based on the rule conditions. Valid values are: `"ALLOW"`, `"DENY"`. Default is `"ALLOW"`.

- `profile_attributes` - (Optional) A list of attributes to prompt the user during registration or progressive profiling.

    - `label` - (Required) A display-friendly label for this property
    - `name` - (Required) The name of a User Profile property
    - `required` - (Optional) Indicates if this property is required for enrollment. Default is `false`.
    
## Attributes Reference

- `id` - ID of the Rule.

- `name` - Name of the Rule.

- `status` - Status of the Rule.

## Import

A Policy Rule can be imported via the Policy and Rule ID.

```
$ terraform import okta_policy_rule_profile_enrollment.example <policy id>/<rule id>
```
