---
layout: "okta"
page_title: "Okta: okta_user"
sidebar_current: "docs-okta-datasource-user"
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
```

## Arguments Reference

* `user_id` - (Optional) String representing a specific user's id value

* `search` - (Optional) Map of search criteria. It supports the following properties.
  * `name` - (Required) Name of property to search against.
  * `comparison` - (Optional) Comparison to use.
  * `value` - (Required) Value to compare with.

## Attributes Reference

* `admin_roles` - Administrator roles assigned to user.

* `city` - user profile property.

* `cost_center` - user profile property.

* `country_code` - user profile property.

* `custom_profile_attributes` - raw JSON containing all custom profile attributes.

* `department` - user profile property.

* `display_name` - user profile property.

* `division` - user profile property.

* `email` - user profile property.

* `employee_number` - user profile property.

* `first_name` - user profile property.

* `group_memberships` - user profile property.

* `honorific_prefix` - user profile property.

* `honorific_suffix` - user profile property.

* `last_name` - user profile property.

* `locale` - user profile property.

* `login` - user profile property.

* `manager` - user profile property.

* `manager_id` - user profile property.

* `middle_name` - user profile property.

* `mobile_phone` - user profile property.

* `nick_name` - user profile property.

* `organization` - user profile property.

* `postal_address` - user profile property.

* `preferred_language` - user profile property.

* `primary_phone` - user profile property.

* `profile_url` - user profile property.

* `second_email` - user profile property.

* `state` - user profile property.

* `status` - user profile property.

* `street_address` - user profile property.

* `timezone` - user profile property.

* `title` - user profile property.

* `user_type` - user profile property.

* `zip_code` - user profile property.
