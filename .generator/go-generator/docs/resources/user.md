---
page_title: "okta_user Resource - terraform-provider-okta"
description: |-
  Manages an Okta User resource.
---

# okta_user (Resource)

Manages an Okta User.

## Example Usage

```hcl
resource "okta_user" "example" {

  # Optional fields
  # algorithm = "<algorithm>"
  # digest_algorithm = "<digest_algorithm>"
  # iteration_count = 0
  # key_size = 0
  # salt = "<salt>"
}
```

## Schema

### Optional

- `algorithm` (String) The algorithm used to generate the hash using the password (and salt, when applicable).
- `digest_algorithm` (String) Algorithm used to generate the key. Only required for the PBKDF2 algorithm.
- `iteration_count` (Number) The number of iterations used when hashing passwords using PBKDF2. Must be >= 4096. Only required for PBKDF2 algorithm.
- `key_size` (Number) Size of the derived key in bytes. Only required for PBKDF2 algorithm.
- `salt` (String) Only required for salted hashes. For BCRYPT, this specifies Radix-64 as the encoded salt used to generate the hash, which must be 22 characters long. For other salted hashes, this specifies the Bas...
- `salt_order` (String) Specifies whether salt was pre- or postfixed to the password before hashing. Only required for salted algorithms.
- `value` (String) For SHA-512, SHA-256, SHA-1, MD5, and PBKDF2, this is the actual base64-encoded hash of the password (and salt, if used). This is the Base64-encoded `value` of the SHA-512/SHA-256/SHA-1/MD5/PBKDF2 ...
- `work_factor` (Number) Governs the strength of the hash and the time required to compute it. Only required for BCRYPT algorithm.
- `type` (String) The type of password inline hook. Currently, must be set to default.
- `value` (String) Specifies the password for a user. The password policy validates this password.
- `name` (String) The name of the authentication provider
- `type` (String) The type of authentication provider
- `answer` (String) The answer to the recovery question
- `question` (String) The recovery question
- `city` (String) The city or locality of the user's address (`locality`)
- `cost_center` (String) Name of the cost center assigned to a user
- `country_code` (String) The country name component of the user's address (`country`). For validation, see [ISO 3166-1 alpha 2 'short' code format](https://datatracker.ietf.org/doc/html/draft-ietf-scim-core-schema-22#ref-I...
- `department` (String) Name of the user's department
- `display_name` (String) Name of the user suitable for display to end users
- `division` (String) Name of the user's division
- `email` (String) The primary email address of the user. For validation, see [RFC 5322 Section 3.2.3](https://datatracker.ietf.org/doc/html/rfc5322#section-3.2.3).
- `employee_number` (String) The organization or company assigned unique identifier for the user
- `first_name` (String) Given name of the user (`givenName`)
- `honorific_prefix` (String) Honorific prefix(es) of the user, or title in most Western languages
- `honorific_suffix` (String) Honorific suffix(es) of the user
- `last_name` (String) The family name of the user (`familyName`)
- `locale` (String) The user's default location for purposes of localizing items such as currency, date time format, numerical representations, and so on. A locale value is a concatenation of the ISO 639-1 two-letter ...
- `login` (String) The unique identifier for the user (`username`). For validation, see [Login pattern validation](https://developer.okta.com/docs/reference/api/schemas/#login-pattern-validation).  Every user within ...
- `manager` (String) The `displayName` of the user's manager
- `manager_id` (String) The `id` of the user's manager
- `middle_name` (String) The middle name of the user
- `mobile_phone` (String) The mobile phone number of the user
- `nick_name` (String) The casual way to address the user in real life
- `organization` (String) Name of the the user's organization
- `postal_address` (String) Mailing address component of the user's address
- `preferred_language` (String) The user's preferred written or spoken language. For validation, see [RFC 7231 Section 5.3.5](https://datatracker.ietf.org/doc/html/rfc7231#section-5.3.5).
- `primary_phone` (String) The primary phone number of the user such as a home number
- `profile_url` (String) The URL of the user's online profile. For example, a web page. See [URL](https://datatracker.ietf.org/doc/html/rfc1808).
- `second_email` (String) The secondary email address of the user typically used for account recovery. For validation, see [RFC 5322 Section 3.2.3](https://datatracker.ietf.org/doc/html/rfc5322#section-3.2.3).
- `state` (String) The state or region component of the user's address (`region`)
- `street_address` (String) The full street address component of the user's address
- `timezone` (String) The user's time zone
- `title` (String) The user's title, such as Vice President
- `user_type` (String) The property used to describe the organization-to-user relationship, such as employee or contractor  > **Note:** The `userType` property is a standard string attribute and should be treated as a de...
- `zip_code` (String) The ZIP code or postal code component of the user's address (`postalCode`)
- `type` (String) The user type that determines the schema for the user's profile. The `type` property is a map that identifies the [User Types](https://developer.okta.com/docs/api/openapi/okta-management/management...
- `group_ids` (List) The list of group IDs of groups that the user is added to at the time of creation

### Read-Only

- `id` (String) The unique identifier for the resource.
- `embedded` (String) Embedded resources related to the user using the [JSON Hypertext Application Language](https://datatracker.ietf.org/doc/html/draft-kelly-json-hal-06) specification
- `activated` (String) The timestamp when the user status transitioned to `ACTIVE`
- `created` (String) The timestamp when the user was created
- `last_login` (String) The timestamp of the last login
- `last_updated` (String) The timestamp when the user was last updated
- `password_changed` (String) The timestamp when the user's password was last updated
- `realm_id` (String) The ID of the realm in which the user is residing. See [Realms](/openapi/okta-management/management/tags/realm).
- `status` (String) The current status of the user.  The status of a user changes in response to explicit events, such as admin-driven lifecycle changes, user login, or self-service password recovery. Okta doesn't asy...
- `status_changed` (String) The timestamp when the status of the user last changed
- `transitioning_to_status` (String) The target status of an in-progress asynchronous status transition. This property is only returned if the user's state is transitioning.
