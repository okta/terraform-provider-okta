---
layout: 'okta'
page_title: 'Okta: okta_policy_profile_enrollment_apps'
sidebar_current: 'docs-okta-resource-policy-profile-enrollment-apps'
description: |-
  Manages Profile Enrollment Policy Apps
---

# okta_policy_profile_enrollment_apps

~> **WARNING:** This feature is only available as a part of the Identity Engine. [Contact support](mailto:dev-inquiries@okta.com) for further information.

This resource allows you to manage the apps in the Profile Enrollment Policy. 

**Important Notes:** 
 - Default Enrollment Policy can not be used in this resource since it is used as a policy to re-assign apps to when they are unassigned from this one.
 - When re-assigning the app to another policy, please use `depends_on` in the policy to which the app will be assigned. This is necessary to avoid 
  unexpected behavior, since if the app is unassigned from the policy it is just assigned to the `Default` one.

## Example Usage

```hcl
data "okta_policy" "example" {
  name = "My Policy"
  type = "PROFILE_ENROLLMENT"
}

data "okta_app" "test" {
  label = "My App"
}

resource "okta_policy_profile_enrollment_apps" "example" {
  policy_id = okta_policy.example.id
  apps      = [data.okta_app.id]
}
```

## Argument Reference

The following arguments are supported:

- `policy_id` - (Required) ID of the enrollment policy.

- `apps` - (Optional) List of app IDs to be added to this policy.

## Attributes Reference

- `id` - ID of the enrollment policy.

- `default_policy_id` - ID of the default enrollment policy.

## Import

A Profile Enrollment Policy Apps can be imported via the Okta ID.

```
$ terraform import okta_policy_profile_enrollment_apps.example &#60;policy id&#62;
```
