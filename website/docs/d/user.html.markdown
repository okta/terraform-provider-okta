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
- `delay_read_seconds` - (Optional) Force delay of the user read by N seconds. Useful when eventual consistency of user information needs to be allowed for.

## Attributes Reference

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
- `honorific_prefix` - Honorific prefix(es) of the user, or title in most Western languages.
- `honorific_suffix` - Honorific suffix(es) of the user.
- `id` - User ID.
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
- `roles` - All roles assigned to user.
- `second_email` - Secondary email address of user typically used for account recovery.
- `state` - State or region component of user's address (region).
- `status` - Current status of user.
- `street_address` - Full street address component of user's address.
- `timezone` - User's time zone.
- `title` - User's title, such as "Vice President".
- `user_type` - Used to describe the organization to user relationship such as "Employee" or "Contractor".
- `zip_code` - Zipcode or postal code component of user's address (postalCode)
