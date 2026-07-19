---
page_title: "okta_application_secure_password_store Resource - terraform-provider-okta"
description: |-
  Manages an Okta Application Secure Password Store resource.
---

# okta_application_secure_password_store (Resource)

Manages an Okta Application Secure Password Store.

## Example Usage

```hcl
resource "okta_application_secure_password_store" "example" {
  sign_on_mode = "<sign_on_mode>"
  label = "Example Label"
  name = "Example Name"
  password_field = "<password_field>"
  url = "https://example.com"
  username_field = "Example Username Field"

  # Optional fields
  # error_redirect_url = "https://example.com"
  # login_redirect_url = "https://example.com"
  # self_service = true
  # algorithm = "<algorithm>"
  # digest_algorithm = "<digest_algorithm>"
}
```

## Schema

### Required

- `sign_on_mode` (String) Discriminator field identifying the variant type. Must be set to \
- `label` (String) User-defined display name for app
- `name` (String) `template_sps` is the key name for a SWA app instance that uses HTTP POST and doesn't require a browser plugin
- `password_field` (String) CSS selector for the **Password** field in the sign-in form
- `url` (String) The URL of the sign-in page for this app
- `username_field` (String) CSS selector for the **Username** field in the sign-in form

### Optional

- `error_redirect_url` (String) Custom error page URL for the app
- `login_redirect_url` (String) Custom login page URL for the app > **Note:** The `loginRedirectUrl` property is deprecated in Identity Engine. This property is used with the custom app login feature. Orgs that actively use this ...
- `self_service` (Boolean) Represents whether the app can be self-assignable by users
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
- `reveal_password` (Boolean) Allow users to securely see their password
- `scheme` (String) Apps with `BASIC_AUTH`, `BROWSER_PLUGIN`, or `SECURE_PASSWORD_STORE` sign-on modes have credentials vaulted by Okta and can be configured with the following schemes.
- `kid` (String) Key identifier used for signing assertions > **Note:** Currently, only the X.509 JWK format is supported for apps with SAML_2_0 `signOnMode`.
- `rotation_mode` (String) The mode of key rotation
- `use` (String) Specifies the intended use of the key
- `user_name` (String) Shared username for the app
- `push_status` (String) Determines if the username is pushed to the app on updates for CUSTOM `type`
- `template` (String) Mapping expression used to generate usernames.  The following are supported mapping expressions that are used with the `BUILT_IN` template type:  | Name                            | Template Expres...
- `type` (String) Type of mapping expression. Empty string is allowed.
- `user_suffix` (String) An optional suffix appended to usernames for `BUILT_IN` mapping expressions
- `seat_count` (Number) Number of licenses purchased for the app
- `profile` (String) Contains any valid JSON schema for specifying properties that can be referenced from a request (only available to OAuth 2.0 client apps). For example, add an app manager contact email address or de...
- `optional_field1` (String) Name of the optional parameter in the sign-in form
- `optional_field1_value` (String) Name of the optional value in the sign-in form
- `optional_field2` (String) Name of the optional parameter in the sign-in form
- `optional_field2_value` (String) Name of the optional value in the sign-in form
- `optional_field3` (String) Name of the optional parameter in the sign-in form
- `optional_field3_value` (String) Name of the optional value in the sign-in form
- `identity_store_id` (String) Identifies an additional identity store app, if your app supports it. The `identityStoreId` value must be a valid identity store app ID. This identity store app must be created in the same org as y...
- `implicit_assignment` (Boolean) Controls whether Okta automatically assigns users to the app based on the user's role or group membership.
- `inline_hook_id` (String) Identifier of an inline hook. Inline hooks are outbound calls from Okta to your own custom code, triggered at specific points in Okta process flows. They allow you to integrate custom functionality...
- `admin` (String) An app message that's visible to admins
- `enduser` (String) A message that's visible in the End-User Dashboard
- `help_url` (String) An optional URL to a help page to assist your end users in signing in to your company VPN
- `message` (String) A VPN requirement message that's displayed to users
- `connection` (String) Specifies the VPN connection details required to access the app
- `exclude` (List) Defines the IP addresses or network ranges that are excluded from the VPN requirement
- `include` (List) Defines the IP addresses or network ranges that are required to use the VPN
- `app_links` (String) Links or icons that appear on the End-User Dashboard if they're set to `true`.
- `auto_launch` (Boolean) Automatically signs in to the app when user signs into Okta
- `auto_submit_toolbar` (Boolean) Automatically sign in when user lands on the sign-in page
- `i_os` (Boolean) Okta Mobile for iOS or Android (pre-dates Android)
- `web` (Boolean) Okta End-User Dashboard on a web browser

### Read-Only

- `id` (String) The unique identifier for the resource.
- `embedded` (String) Embedded resources related to the app using the [JSON Hypertext Application Language](https://datatracker.ietf.org/doc/html/draft-kelly-json-hal-06) specification. If the `expand=user/{userId}` que...
- `created` (String) Timestamp when the application object was created
- `last_rotated` (String) Timestamp when the signing key was last rotated
- `next_rotation` (String) The scheduled time for the next signing key rotation
- `express_configuration` (String) <div class='x-lifecycle-container'><x-lifecycle class='oie'></x-lifecycle></div> Indicates which Express Configuration capabilities the app supports and has enabled
- `features` (List) Enabled app features > **Note:** See [Application Features](/openapi/okta-management/management/tags/applicationfeatures/) for app provisioning features.
- `last_updated` (String) Timestamp when the application object was last updated
- `orn` (String) The Okta resource name (ORN) for the current app instance
- `em_opt_in_status` (String) The entitlement management opt-in status for the app
- `status` (String) App instance status
- `universal_logout` (String) <div class='x-lifecycle-container'><x-lifecycle class='oie'></x-lifecycle></div> Universal Logout properties for the app. These properties are only returned and can't be updated.
