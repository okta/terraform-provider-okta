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

# Search for multiple users based on a raw search expression string
data "okta_users" "example" {
  search {
    expression = "profile.department eq \"Engineering\" and (created lt \"2014-01-01T00:00:00.000Z\" or status eq \"ACTIVE\")"
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

  # optionally include each user's administrator roles
  include_roles = true
}
```

## Arguments Reference

- `search` - (Optional) Map of search criteria. It supports the following properties.
  - `name` - (Required w/ comparison and value) Name of property to search against.
  - `comparison` - (Required w/ name and value) Comparison to use. Comparitors for strings: [`eq`, `ge`, `gt`, `le`, `lt`, `ne`, `pr`, `sw`](https://developer.okta.com/docs/reference/core-okta-api/#operators).
  - `value` - (Required w/ comparison and name) Value to compare with.
  - `expression` - (Optional, but overrides name/comparison/value) A raw search expression string. If present it will override name/comparison/value.
- `compound_search_operator` - (Optional) Given multiple search elements they will be compounded together with the op. Default is `and`, `or` is also valid.
- `group_id` - (Optional) Id of group used to find users based on membership.
- `include_groups` - (Optional) Fetch each user's group memberships. Defaults to `false`, in which case the `group_memberships` user attribute will be empty.
- `include_roles` - (Optional) Fetch each user's administrator roles. Defaults to `false`, in which case the `admin_roles` user attribute will be empty.
- `delay_read_seconds` - (Optional) Force delay of the users read by N seconds. Useful when eventual consistency of users information needs to be allowed for; for instance, when administrator roles are known to have been applied.

## Attributes Reference

- `users` - collection of users retrieved from Okta with the following properties.
  - `admin_roles` - Administrator roles assigned to user.
  - `city` - City or locality component of user's address.
  - `cost_center` - Name of a cost center assigned to user.
  - `country_code` - Country name component of user's address.
  - `custom_profile_attributes` - Raw JSON containing all custom profile attributes.
  - `department` - Name of user's department.
  - `display_name` - Name of the user, suitable for display to end users.
  - `division` - Name of user's division.
  - `email` - Primary email address of user.
  - `employee_number` - Organization or company assigned unique identifier for the user.
  - `first_name` - Given name of the user.
  - `group_memberships` - Groups user belongs to.
  - `honorific_prefix` - Honorific prefix(es) of the user, or title in most Western languages.
  - `honorific_suffix` - Honorific suffix(es) of the user.
  - `last_name` - Family name of the user.
  - `locale` - User's default location for purposes of localizing items such as currency, date time format, numerical representations, etc.
  - `login` - Unique identifier for the user.
  - `manager_id` - `id` of a user's manager.
  - `manager` - Display name of the user's manager.
  - `middle_name` - Middle name(s) of the user.
  - `mobile_phone` - Mobile phone number of user.
  - `nick_name` - Casual way to address the user in real life.
  - `organization` - Name of user's organization.
  - `postal_address` - Mailing address component of user's address.
  - `preferred_language` - User's preferred written or spoken languages.
  - `primary_phone` - Primary phone number of user such as home number.
  - `profile_url` - URL of user's online profile (e.g. a web page).
  - `second_email` - Secondary email address of user typically used for account recovery.
  - `state` - State or region component of user's address (region).
  - `status` - Current status of user.
  - `street_address` - Full street address component of user's address.
  - `timezone` - User's time zone.
  - `title` - User's title, such as "Vice President".
  - `user_type` - Used to describe the organization to user relationship such as "Employee" or "Contractor".
  - `zip_code` - Zipcode or postal code component of user's address (postalCode)