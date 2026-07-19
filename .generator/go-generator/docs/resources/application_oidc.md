---
page_title: "okta_application_oidc Resource - terraform-provider-okta"
description: |-
  Manages an Okta Application Oidc resource.
---

# okta_application_oidc (Resource)

Manages an Okta Application Oidc.

## Example Usage

```hcl
resource "okta_application_oidc" "example" {
  sign_on_mode = "<sign_on_mode>"
  label = "Example Label"
  name = "Example Name"
  grant_types = "<grant_types>"
  mode = "<mode>"
  connection = "<connection>"
  rotation_type = "<rotation_type>"

  # Optional fields
  # error_redirect_url = "https://example.com"
  # login_redirect_url = "https://example.com"
  # self_service = true
  # auto_key_rotation = true
  # client_id = "<client-id>"
}
```

## Schema

### Required

- `sign_on_mode` (String) Discriminator field identifying the variant type. Must be set to \
- `label` (String) User-defined display name for app
- `name` (String) `oidc_client` is the key name for an OAuth 2.0 client app instance
- `grant_types` (List) GrantTypes
- `mode` (String) The mode to use for the IdP-initiated sign-in flow. For `OKTA` or `SPEC` modes, the client must have an `initiate_login_uri` registered. > **Note:** For web and SPA apps, if the mode is `SPEC` or `...
- `connection` (String) The connection type of the network. Can be `ANYWHERE` or `ZONE`.
- `rotation_type` (String) The refresh token rotation mode for the OAuth 2.0 client

### Optional

- `error_redirect_url` (String) Custom error page URL for the app
- `login_redirect_url` (String) Custom login page URL for the app > **Note:** The `loginRedirectUrl` property is deprecated in Identity Engine. This property is used with the custom app login feature. Orgs that actively use this ...
- `self_service` (Boolean) Represents whether the app can be self-assignable by users
- `auto_key_rotation` (Boolean) Requested key rotation mode
- `client_id` (String) Unique identifier for the OAuth 2.0 client app  > **Notes:** > * If you don't specify the `client_id`, this immutable property is populated with the [Application instance ID](/openapi/okta-manageme...
- `client_secret` (String) OAuth 2.0 client secret string (used for confidential clients)  > **Notes:** If a `client_secret` isn't provided on creation, and the `token_endpoint_auth_method` requires one, Okta generates a ran...
- `pkce_required` (Boolean) Requires Proof Key for Code Exchange (PKCE) for additional verification. If `token_endpoint_auth_method` is `none`, then `pkce_required` must be `true`. The default is `true` for browser and native...
- `token_endpoint_auth_method` (String) Requested authentication method for the token endpoint
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
- `application_type` (String) The type of client app Specific `grant_types` are valid for each `application_type`. See [Create a Client Application](/openapi/okta-oauth/oauth/client/createclient).
- `backchannel_authentication_request_signing_alg` (String) The signing algorithm for Client-Initiated Backchannel Authentication (CIBA) signed requests using JWT. If this value isn't set and a JWT-signed request is sent, the request fails. > **Note:** This...
- `backchannel_custom_authenticator_id` (String) The ID of the custom authenticator that authenticates the user > **Note:** This property appears for clients with `urn:openid:params:grant-type:ciba` defined as one of the `grant_types`.
- `backchannel_token_delivery_mode` (String) The delivery mode for Client-Initiated Backchannel Authentication (CIBA).  Currently, only `poll` is supported. > **Note:** This property appears for clients with `urn:openid:params:grant-type:ciba...
- `client_uri` (String) URL string of a web page providing information about the client
- `consent_method` (String) Indicates whether user consent is required or implicit. A consent dialog appears for the end user depending on the values of three elements:  * [prompt](/openapi/okta-oauth/oauth/orgas/authorize#or...
- `dpop_bound_access_tokens` (Boolean) Indicates that the client application uses Demonstrating Proof-of-Possession (DPoP) for token requests. If `true`, the authorization server rejects token requests from this client that don't contai...
- `frontchannel_logout_session_required` (Boolean) <x-lifecycle-container><x-lifecycle class='ea'></x-lifecycle> <x-lifecycle class='oie'></x-lifecycle></x-lifecycle-container>Determines whether Okta sends `sid` and `iss` in the logout request
- `frontchannel_logout_uri` (String) <x-lifecycle-container><x-lifecycle class='ea'></x-lifecycle> <x-lifecycle class='oie'></x-lifecycle></x-lifecycle-container>URL where Okta sends the logout request
- `id_token_encrypted_response_alg` (String) JWE alg algorithm for encrypting the ID token issued to this client. If this is requested, the response is signed, and then encrypted with the result being a nested JWT. The default, if omitted, is...
- `default_scope` (List) The scopes to use for the request when `mode` is `OKTA`
- `initiate_login_uri` (String) URL string that a third party can use to initiate the sign-in flow by the client
- `issuer_mode` (String) Indicates whether the Okta authorization server uses the original Okta org domain URL or a custom domain URL as the issuer of the ID token for this client
- `keys` (List) Keys
- `jwks_uri` (String) URL string that references a JSON Web Key Set for validating JWTs presented to Okta or for encrypting ID tokens minted by Okta for the client
- `logo_uri` (String) The URL string that references a logo for the client. This logo appears on the client tile in the End-User Dashboard. It also appears on the client consent dialog during the client consent flow.
- `exclude` (List) If `ZONE` is specified as a connection, then specify the excluded IP network zones here. Value can be 'ALL_IP_ZONES' or an array of zone IDs.
- `include` (List) If `ZONE` is specified as a connection, then specify the included IP network zones here. Value can be 'ALL_IP_ZONES' or an array of zone IDs.
- `participate_slo` (Boolean) <x-lifecycle-container><x-lifecycle class='ea'></x-lifecycle> <x-lifecycle class='oie'></x-lifecycle></x-lifecycle-container>Allows the app to participate in front-channel Single Logout  > **Note:*...
- `policy_uri` (String) URL string of a web page providing the client's policy document
- `post_logout_redirect_uris` (List) Array of redirection URI strings for relying party-initiated logouts
- `redirect_uris` (List) Array of redirection URI strings for use in redirect-based flows. > **Note:** At least one `redirect_uris` and `response_types` are required for all client types, with exceptions: if the client use...
- `leeway` (Number) The leeway, in seconds, allowed for the OAuth 2.0 client. After the refresh token is rotated, the previous token remains valid for the specified period of time so clients can get the new token.  > ...
- `request_object_signing_alg` (String) The type of JSON Web Key Set (JWKS) algorithm that must be used for signing request objects
- `response_types` (List) Array of OAuth 2.0 response type strings
- `sector_identifier_uri` (String) The sector identifier used for pairwise `subject_type`. See [OIDC Pairwise Identifier Algorithm](https://openid.net/specs/openid-connect-messages-1_0-20.html#idtype.pairwise.alg)
- `subject_type` (String) Type of the subject
- `tos_uri` (String) URL string of a web page providing the client's terms of service document
- `wildcard_redirect` (String) Indicates if the client is allowed to use wildcard matching of `redirect_uris`
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
