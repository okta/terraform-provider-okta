# Configuration Examples

Here lies the examples that will aid you on your Okta Terraform journey. PLEASE NOTE not all resources are yet outlined
here. Just some really common ones! We will stand up a wiki soon.

## Example Stacks

- [Okta and Cognito](./oidc-cognito-stack.tf) Example of using Okta OIDC application with a Cognito ID Provider to
  provide a serverless SPA access to AWS resources.

## Test Fixture Examples

Anything that lies underneath a resource directory is config we use as fixtures to our tests. We use the special
trigger `replace_with_uuid` to trigger the fixture manager to insert a unique id. Eventually we will stand up a wiki for
the provider, but we don't quite have the manpower yet to do so.

## Resources 

- [okta_app_auto_login](./resources/okta_app_auto_login) Supports the management of Okta Auto Login Applications.
- [okta_app_bookmark](./resources/okta_app_bookmark) Supports the management Okta Bookmark Application.
- [okta_app_oauth](./resources/okta_app_oauth) Supports the management of Okta OIDC Applications.
- [okta_app_saml](./resources/okta_app_saml) Supports the management of Okta SAML Applications.
- [okta_app_secure_password_store](./resources/okta_app_secure_password_store) Supports the management of Okta Secure Password
  Store Applications.
- [okta_app_swa](./resources/okta_app_swa) Supports the management of Okta SWA Applications.
- [okta_app_three_field](./resources/okta_app_three_field) Supports the management of Okta Three Field Applications.
- [okta_auth_server_claim](./resources/okta_auth_server_claim) Supports the management of Okta Authorization servers claims.
- [okta_auth_server_policy_rule](./resources/okta_auth_server_policy_rule) Supports the management of Okta Authorization servers
  policy rules.
- [okta_auth_server_policy](./resources/okta_auth_server_policy) Supports the management of Okta Authorization servers policies.
- [okta_auth_server_scope](./resources/okta_auth_server_scope) Supports the management of Okta Authorization servers scopes.
- [okta_auth_server](./resources/okta_auth_server) Supports the management of Okta Authorization servers.
- [okta_group_rule](./resources/okta_group_rule) Supports the management of Okta Group Rules.
- [okta_group](./resources/okta_group) Supports the management of Okta Groups.
- [okta_event_hook](./resources/okta_event_hook) Supports the management of Okta Event Hooks.
- [okta_idp_saml](./resources/okta_idp_saml) Supports the management of Okta SAML Identity Providers.
- [okta_idp_social](./resources/okta_idp_social) Supports the management of Okta Social Identity Providers. Such as Google,
  Facebook, Microsoft, and LinkedIn.
- [okta_inline_hook](./resources/okta_inline_hook) Supports the management of Okta Inline Hooks EA feature.
- [okta_network_zone](./resources/okta_network_zone) Supports the management of Okta Network Zones for whitelisting IPs or
  countries dynamically.
- [okta_policy_mfa](./resources/okta_policy_mfa) Supports the management of MFA policies.
- [okta_policy_password](./resources/okta_policy_password) Supports the management of password policies.
- [okta_policy_rule_signon](./resources/okta_policy_rule_signon) Supports the management of sign-on policy rules.
- [okta_policy_signon](./resources/okta_policy_signon) Supports the management of sign-on policies.
- [okta_trusted_origin](./resources/okta_trusted_origin) Supports the management of Okta Trusted Sources and Origins.
- [okta_user_base_schema_property](./resources/okta_user_base_schema_property) Supports the management of Okta User Profile
  Attribute Schemas.
- [okta_user_schema_property](./resources/okta_user_schema_property) Supports the management of Okta defined User Profile
  Attribute Schemas.
- [okta_user](./resources/okta_user) Supports the management of Okta Users.
- [okta_app_oauth_post_logout_redirect_uri](./resources/okta_app_oauth_post_logout_redirect_uri) Supports decentralizing post logout redirect uri config. 
- [okta_app_oauth_redirect_uri](./resources/okta_app_oauth_redirect_uri) Supports decentralizing redirect uri config. Due to
  Okta's API not allowing this field to be null, you must set a redirect uri in your app, and ignore changes to this
  attribute. We follow TF best practices and detect config drift. The best case scenario is Okta makes this field
  nullable, and we can not detect config drift when this attr is not present.

##  Data Sources

- [okta_app_saml_metadata](./data-sources/okta_app_saml_metadata) Data source for SAML app metadata.
- [okta_app](./data-sources/okta_app) Generic Application data source.
- [okta_idp_metadata_saml](./data-sources/okta_app_metadata_saml) Data source for SAML IdP metadata.
- [okta_users](./data-sources/okta_users) Data source to retrieve a group of users.

## Deprecated Resources

- okta_identity_provider -- See okta_idp, okta_idp_social, and okta_idp_saml.
- okta_user_schema_propertys -- See okta_user_schema_property.

## Notes

As resource fixtures are added, please be sure to only put a VALID config in each resource sub-directory. Intentionally
invalid config for testing should stay in the test file.
