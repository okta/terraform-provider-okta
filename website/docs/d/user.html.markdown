---
layout: 'okta'
page_title: 'Okta: okta_user'
sidebar_current: 'docs-okta-datasource-user'
description: |-
  Get a single users from Okta.
---

# okta_user

Use this data source to retrieve a users from Okta.

## Example Usage

```hcl
# Get a single user by their id value
data "okta_user" "example" {
  user_id = "00u22mtxlrJ8YkzXQ357"
}

# Search for a single user based on supported profile properties
data "okta_user" "example" {
  search {
    name  = "profile.firstName"
    value = "John"
  }

  search {
    name  = "profile.lastName"
    value = "Doe"
  }
}

# Search for a single user based on a raw search expression string
data "okta_user" "example" {
  search {
    expression  = "profile.firstName eq \"John\""
  }
}
```

## Arguments Reference

- `user_id` - (Optional) String representing a specific user's id value

- `search` - (Optional) Map of search criteria. It supports the following properties.
  - `name` - (Required w/ comparison and value) Name of property to search against.
  - `comparison` - (Required w/ name and value) Comparison to use. Comparitors for strings: [`eq`, `ge`, `gt`, `le`, `lt`, `ne`, `pr`, `sw`](https://developer.okta.com/docs/reference/core-okta-api/#operators).
  - `value` - (Required w/ comparison and name) Value to compare with.
  - `expression` - (Optional, but overrides name/comparison/value) A raw search expression string. If present it will override name/comparison/value.
- `compound_search_operator` - (Optional) Given multiple search elements they will be compounded together with the op. Default is `and`, `or` is also valid.

- `skip_groups` - (Optional) Additional API call to collect user's groups will not be made.

- `skip_roles` - (Optional) Additional API call to collect user's roles will not be made.

## Attributes Reference

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
