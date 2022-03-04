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


### Lookup Users by Search Criteria

```hcl
data "okta_users" "example" {
  search {
    name       = "profile.company"
    value      = "Articulate"
    comparison = "sw"
  }
}
```

### Lookup Users by Group Membership
```hcl
resource "okta_group" "example" {
  name = "example-group"
}

data "okta_users" "example" {
  group_id = okta_group.example.id
  
  # optionally include each user's group membership
  include_groups = true
}
```

## Arguments Reference

- `search` - (Optional) Map of search criteria to find users. It supports the following properties. 
  - `name` - (Required) Name of property to search against.
  - `comparison` - (Required) Comparison to use.
  - `value` - (Required) Value to compare with.

- `group_id` - (Optional) Id of group used to find users based on membership.

- `include_groups` - (Optional) Fetch each user's group memberships. Defaults to `false`, in which case the `group_memberships` user attribute will be empty.

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
