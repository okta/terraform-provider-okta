---
layout: "okta"
page_title: "Okta: okta_user"
sidebar_current: "docs-okta-resource-user"
description: |-
  Creates an Okta User.
---

# okta_user

Creates an Okta User.

This resource allows you to create and configure an Okta User.

## Example Usage

```hcl
resource "okta_user" "example" {
  index       = "customPropertyName"
  title       = "customPropertyName"
  type        = "string"
  description = "My custom property name"
  master      = "OKTA"
  scope       = "SELF"
}
```

## Argument Reference

The following arguments are supported:

* `email` - (Required) User profile property.

* `login` - (Required) User profile property.

* `first_name` - (Required) User's First Name, required by default.

* `last_name` - (Required) User's Last Name, required by default.

* `custom_profile_attributes` - (Optional) raw JSON containing all custom profile attributes.

* `admin_roles` - (Optional) Administrator roles assigned to User.

* `city` - (Optional) User profile property.

* `cost_center` - (Optional) User profile property.

* `country_code` - (Optional) User profile property.

* `department` - (Optional) User profile property.

* `display_name` - (Optional) User profile property.

* `division` - (Optional) User profile property.

* `employee_number` - (Optional) User profile property.

* `group_memberships` - (Optional) User profile property.

* `honorific_prefix` - (Optional) User profile property.

* `honorific_suffix` - (Optional) User profile property.

* `locale` - (Optional) User profile property.

* `manager` - (Optional) User profile property.

* `manager_id` - (Optional) User profile property.

* `middle_name` - (Optional) User profile property.

* `mobile_phone` - (Optional) User profile property.

* `nick_name` - (Optional) User profile property.

* `organization` - (Optional) User profile property.

* `postal_address` - (Optional) User profile property.

* `preferred_language` - (Optional) User profile property.

* `primary_phone` - (Optional) User profile property.

* `profile_url` - (Optional) User profile property.

* `second_email` - (Optional) User profile property.

* `state` - (Optional) User profile property.

* `status` - (Optional) User profile property.

* `street_address` - (Optional) User profile property.

* `timezone` - (Optional) User profile property.

* `title` - (Optional) User profile property.

* `user_type` - (Optional) User profile property.

* `zip_code` - (Optional) User profile property.

* `password` - (Optional) User password.

* `recovery_question` - (Optional) User password recovery question.

* `recovery_answer` - (Optional) User password recovery answer.

## Attributes Reference

* `index` - (Optional) ID of the User schema property.

## Import

An Okta User can be imported via the ID.

```
$ terraform import okta_user.example <user id>
```
