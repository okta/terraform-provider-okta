---
page_title: "okta_application_ws_federation Resource - terraform-provider-okta"
description: |-
  Manages an Okta Application Ws Federation resource.
---

# okta_application_ws_federation (Resource)

Manages an Okta Application Ws Federation.

## Example Usage

```hcl
resource "okta_application_ws_federation" "example" {
  sign_on_mode = "<sign_on_mode>"
  label = "Example Label"
  name = "Example Name"
  audience_restriction = "<audience_restriction>"
  authn_context_class_ref = "<authn_context_class_ref>"
  group_value_format = "<group_value_format>"
  name_id_format = "Example Name Id Format"
  site_url = "https://example.com"
  username_attribute = "Example Username Attribute"
  w_reply_url = "https://example.com"

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
- `name` (String) `template_wsfed` is the key name for a WS-Federated app instance with a SAML 2.0 token
- `audience_restriction` (String) The entity ID of the SP. Use the entity ID value exactly as provided by the SP.
- `authn_context_class_ref` (String) Identifies the SAML authentication context class for the assertion's authentication statement
- `group_value_format` (String) Specifies the WS-Fed assertion attribute value for filtered groups. This attribute is only applied to Active Directory groups.
- `name_id_format` (String) The username format that you send in the WS-Fed response
- `site_url` (String) Launch URL for the web app
- `username_attribute` (String) Specifies additional username attribute statements to include in the WS-Fed assertion
- `w_reply_url` (String) The WS-Fed SP endpoint where your users sign in

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
- `attribute_statements` (String) You can federate user attributes such as Okta profile fields, LDAP, Active Directory, and Workday values. The SP uses the federated WS-Fed attribute values accordingly.
- `group_filter` (String) A regular expression that filters for the User Groups you want included with the `groupName` attribute. If the matching User Group has a corresponding AD group, then the attribute statement include...
- `group_name` (String) The group name to include in the WS-Fed response attribute statement. This property is used in conjunction with the `groupFilter` property.  Groups that are filtered through the `groupFilter` expre...
- `realm` (String) The uniform resource identifier (URI) of the WS-Fed app that's used to share resources securely within a domain. It's the identity that's sent to the Okta IdP when signing in. See [Realm name](http...
- `w_reply_override` (Boolean) Enables a web app to override the `wReplyURL` URL with a reply parameter.
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
