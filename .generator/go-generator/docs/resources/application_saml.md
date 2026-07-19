---
page_title: "okta_application_saml Resource - terraform-provider-okta"
description: |-
  Manages an Okta Application Saml resource.
---

# okta_application_saml (Resource)

Manages an Okta Application Saml.

## Example Usage

```hcl
resource "okta_application_saml" "example" {
  sign_on_mode = "<sign_on_mode>"
  label = "Example Label"
  index = 0
  url = "https://example.com"
  allow_multiple_acs_endpoints = "https://example.com"
  assertion_signed = true
  audience = "<audience>"
  authn_context_class_ref = "<authn_context_class_ref>"
  destination = "<destination>"
  digest_algorithm = "<digest_algorithm>"
  honor_force_authn = true
  idp_issuer = "<idp_issuer>"
  recipient = "<recipient>"
  request_compressed = true
  response_signed = true
  signature_algorithm = "<signature_algorithm>"
  sso_acs_url = "https://example.com"
  subject_name_id_format = "Example Subject Name Id Format"
  subject_name_id_template = "Example Subject Name Id Template"

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
- `index` (Number) Index of the URL in the array of ACS endpoints
- `url` (String) URL of the ACS
- `allow_multiple_acs_endpoints` (Boolean) Determines whether the app allows you to configure multiple ACS URIs
- `assertion_signed` (Boolean) Determines whether the SAML assertion is digitally signed
- `audience` (String) The entity ID of the SP. Use the entity ID value exactly as provided by the SP.
- `authn_context_class_ref` (String) Identifies the SAML authentication context class for the assertion's authentication statement
- `destination` (String) Identifies the location inside the SAML assertion where the SAML response should be sent
- `digest_algorithm` (String) Determines the digest algorithm used to digitally sign the SAML assertion and response
- `honor_force_authn` (Boolean) Set to `true` to prompt users for their credentials when a SAML request has the `ForceAuthn` attribute set to `true`
- `idp_issuer` (String) SAML Issuer ID
- `recipient` (String) The location where the app may present the SAML assertion
- `request_compressed` (Boolean) Determines whether the SAML request is expected to be compressed
- `response_signed` (Boolean) Determines whether the SAML authentication response message is digitally signed by the IdP > **Note:** Either (or both) `responseSigned` or `assertionSigned` must be `TRUE`.
- `signature_algorithm` (String) Determines the signing algorithm used to digitally sign the SAML assertion and response
- `sso_acs_url` (String) Single Sign-On Assertion Consumer Service (ACS) URL
- `subject_name_id_format` (String) Identifies the SAML processing rules. Supported values:
- `subject_name_id_template` (String) Template for app user's username when a user is assigned to the app

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
- `enabled` (Boolean) Indicates whether Okta encrypts the assertions that it sends to the Service Provider
- `encryption_algorithm` (String) The encryption algorithm used to encrypt the SAML assertion
- `key_transport_algorithm` (String) The key transport algorithm used to encrypt the SAML assertion
- `x5c` (List) A list that contains exactly one x509 encoded certificate which Okta uses to encrypt the SAML assertion with
- `type` (String) The type of attribute statements object
- `audience_override` (String) Audience override for CASB configuration. See [CASB config guide](https://help.okta.com/en-us/Content/Topics/Apps/CASB-config-guide.htm).
- `type` (String) The type of attribute statements object
- `default_relay_state` (String) Identifies a specific application resource in an IdP-initiated SSO scenario
- `destination_override` (String) Destination override for CASB configuration. See [CASB config guide](https://help.okta.com/en-us/Content/Topics/Apps/CASB-config-guide.htm).
- `inline_hooks` (List) Associates the app with SAML inline hooks. See [the SAML assertion inline hook reference](https://developer.okta.com/docs/reference/saml-hook/).
- `binding_type` (String) Request binding type
- `enabled` (Boolean) Indicates whether the app is allowed to participate in front-channel SLO
- `logout_request_url` (String) URL where Okta sends the logout request
- `session_index_required` (Boolean) Determines whether Okta sends the `SessionIndex` elements in the logout request
- `recipient_override` (String) Recipient override for CASB configuration. See [CASB config guide](https://help.okta.com/en-us/Content/Topics/Apps/CASB-config-guide.htm).
- `saml_assertion_lifetime_seconds` (Number) Determines the SAML app session lifetimes with Okta
- `enabled` (Boolean) Whether the application supports SLO
- `issuer` (String) The issuer of the Service Provider that generates the SLO request
- `logout_url` (String) The location where the logout response is sent
- `x5c` (List) A list that contains exactly one x509 encoded certificate
- `sp_issuer` (String) The issuer ID for the Service Provider. This property appears when SLO is enabled.
- `sso_acs_url_override` (String) Assertion Consumer Service (ACS) URL override for CASB configuration. See [CASB config guide](https://help.okta.com/en-us/Content/Topics/Apps/CASB-config-guide.htm).
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
- `name` (String) A unique key is generated for the custom app instance when you use SAML_2_0 `signOnMode`.
- `orn` (String) The Okta resource name (ORN) for the current app instance
- `em_opt_in_status` (String) The entitlement management opt-in status for the app
- `status` (String) App instance status
- `universal_logout` (String) <div class='x-lifecycle-container'><x-lifecycle class='oie'></x-lifecycle></div> Universal Logout properties for the app. These properties are only returned and can't be updated.
