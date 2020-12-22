---
layout: 'okta'
page_title: 'Okta: okta_users'
sidebar_current: 'docs-okta-datasource-users'
description: |-
  Get a list of users from Okta.
---

# okta_users

Use this data source to retrieve a list of users from Okta.

## Example Usage

```hcl
data "okta_users" "example" {
  search {
    name       = "profile.company"
    value      = "Articulate"
    comparison = "sw"
  }
}
```

## Arguments Reference

- `search` - (Required) Map of search criteria to use to find users. It supports the following properties.
  - `name` - (Required) Name of property to search against.
  - `comparison` - (Required) Comparison to use.
  - `value` - (Required) Value to compare with.

## Attributes Reference

- `users` - collection of users retrieved from Okta with the following properties.
  - `admin_roles` - Administrator roles assigned to user.
  - `city` - user profile property.
  - `cost_center` - user profile property.
  - `country_code` - user profile property.
  - `custom_profile_attributes` - raw JSON containing all custom profile attributes.
  - `department` - user profile property.
  - `display_name` - user profile property.
  - `division` - user profile property.
  - `email` - user profile property.
  - `employee_number` - user profile property.
  - `first_name` - user profile property.
  - `group_memberships` - user profile property.
  - `honorific_prefix` - user profile property.
  - `honorific_suffix` - user profile property.
  - `last_name` - user profile property.
  - `locale` - user profile property.
  - `login` - user profile property.
  - `manager` - user profile property.
  - `manager_id` - user profile property.
  - `middle_name` - user profile property.
  - `mobile_phone` - user profile property.
  - `nick_name` - user profile property.
  - `organization` - user profile property.
  - `postal_address` - user profile property.
  - `preferred_language` - user profile property.
  - `primary_phone` - user profile property.
  - `profile_url` - user profile property.
  - `second_email` - user profile property.
  - `state` - user profile property.
  - `status` - user profile property.
  - `street_address` - user profile property.
  - `timezone` - user profile property.
  - `title` - user profile property.
  - `user_type` - user profile property.
  - `zip_code` - user profile property.
