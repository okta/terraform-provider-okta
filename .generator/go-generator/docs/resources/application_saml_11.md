---
page_title: "okta_application_saml_11 Resource - terraform-provider-okta"
description: |-
  Manages an Okta Application Saml 11 resource.
---

# okta_application_saml_11 (Resource)

Manages an Okta Application Saml 11.

## Example Usage

```hcl
resource "okta_application_saml_11" "example" {
  sign_on_mode = "<sign_on_mode>"
  label = "Example Label"
  name = "Example Name"

  # Optional fields
  # error_redirect_url = "https://example.com"
  # login_redirect_url = "https://example.com"
  # self_service = true
  # kid = "<kid>"
  # rotation_mode = "<rotation_mode>"
}
```

## Schema

### Required

- `sign_on_mode` (String) Discriminator field identifying the variant type. Must be set to \
- `label` (String) User-defined display name for app
- `name` (String) The key name for the SAML 1.1 app definition. You can't create a custom SAML 1.1 app integration instance. Only existing OIN SAML 1.1 app integrations are supported.

### Optional

- `error_redirect_url` (String) Custom error page URL for the app
- `login_redirect_url` (String) Custom login page URL for the app > **Note:** The `loginRedirectUrl` property is deprecated in Identity Engine. This property is used with the custom app login feature. Orgs that actively use this ...
- `self_service` (Boolean) Represents whether the app can be self-assignable by users
- `kid` (String) Key identifier used for signing assertions > **Note:** Currently, only the X.509 JWK format is supported for apps with SAML_2_0 `signOnMode`.
- `rotation_mode` (String) The mode of key rotation
- `use` (String) Specifies the intended use of the key
- `push_status` (String) Determines if the username is pushed to the app on updates for CUSTOM `type`
- `template` (String) Mapping expression used to generate usernames.  The following are supported mapping expressions that are used with the `BUILT_IN` template type:  | Name                            | Template Expres...
- `type` (String) Type of mapping expression. Empty string is allowed.
- `user_suffix` (String) An optional suffix appended to usernames for `BUILT_IN` mapping expressions
- `seat_count` (Number) Number of licenses purchased for the app
- `profile` (String) Contains any valid JSON schema for specifying properties that can be referenced from a request (only available to OAuth 2.0 client apps). For example, add an app manager contact email address or de...
- `app` (String) App
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
- `audience_override` (String) The intended audience of the SAML assertion. This is usually the Entity ID of your application.
- `default_relay_state` (String) The URL of the resource to direct users after they successfully sign in to the SP using SAML. See the SP documentation to check if you need to specify a RelayState. In most instances, you can leave...
- `recipient_override` (String) The location where the application can present the SAML assertion. This is usually the Single Sign-On (SSO) URL.
- `sso_acs_url_override` (String) Assertion Consumer Services (ACS) URL value for the Service Provider (SP). This URL is always used for Identity Provider (IdP) initiated sign-on requests.
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
