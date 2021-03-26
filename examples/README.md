# Configuration Examples

Here lies the examples that will aid you on your Okta Terraform journey. PLEASE NOTE not all resources are yet outlined here. Just some really common ones! We will stand up a wiki soon.

## Example Stacks

- [Okta and Cognito](./oidc-cognito-stack.tf) Example of using Okta OIDC application with a Cognito ID Provider to provide a serverless SPA access to AWS resources.

## Test Fixture Examples

Anything that lies underneath a resource directory is config we use as fixtures to our tests. We use the special trigger `replace_with_uuid` to trigger the fixture manager to insert a unique id. Eventually we will stand up a wiki for the provider, but we don't quite have the manpower yet to do so.

## Resources & Data Sources

- [okta_app_auto_login](./okta_app_auto_login) Supports the management of Okta Auto Login Applications.
- [okta_app_bookmark](./okta_app_bookmark) Supports the management Okta Bookmark Application.
- [okta_app_metadata_saml](./okta_app_metadata_saml) Data source for SAML app metadata.
- [okta_app_oauth](./okta_app_oauth) Supports the management of Okta OIDC Applications.
- [okta_app_saml](./okta_app_saml) Supports the management of Okta SAML Applications.
- [okta_app_secure_password_store](./okta_app_secure_password_store) Supports the management of Okta Secure Password Store Applications.
- [okta_app_swa](./okta_app_swa) Supports the management of Okta SWA Applications.
- [okta_app_three_field](./okta_app_three_field) Supports the management of Okta Three Field Applications.
- [okta_app](./okta_app) Generic Application data source.
- [okta_auth_server_claim](./okta_auth_server_claim) Supports the management of Okta Authorization servers claims.
- [okta_auth_server_policy_rule](./okta_auth_server_policy_rule) Supports the management of Okta Authorization servers policy rules.
- [okta_auth_server_policy](./okta_auth_server_policy) Supports the management of Okta Authorization servers policies.
- [okta_auth_server_scope](./okta_auth_server_scope) Supports the management of Okta Authorization servers scopes.
- [okta_auth_server](./okta_auth_server) Supports the management of Okta Authorization servers.
- [okta_group_roles](./okta_group_roles) Supports the management of Okta Group Administrator Roles.
- [okta_group_rule](./okta_group_rule) Supports the management of Okta Group Rules.
- [okta_group](./okta_group) Supports the management of Okta Groups.
- [okta_event_hook](./okta_event_hook) Supports the management of Okta Event Hooks.
- [okta_idp_metadata_saml](./okta_app_metadata_saml) Data source for SAML IdP metadata.
- [okta_idp_saml](./okta_idp_saml) Supports the management of Okta SAML Identity Providers.
- [okta_idp_social](./okta_idp_social) Supports the management of Okta Social Identity Providers. Such as Google, Facebook, Microsoft, and LinkedIn.
- [okta_inline_hook](./okta_inline_hook) Supports the management of Okta Inline Hooks EA feature.
- [okta_network_zone](./okta_network_zone) Supports the management of Okta Network Zones for whitelisting IPs or countries dynamically.
- [okta_policy_mfa](./okta_policy_mfa) Supports the management of MFA policies.
- [okta_policy_password](./okta_policy_password) Supports the management of password policies.
- [okta_policy_rule_signon](./okta_policy_rule_signon) Supports the management of sign-on policy rules.
- [okta_policy_signon](./okta_policy_signon) Supports the management of sign-on policies.
- [okta_template_email](./okta_template_email) Supports the management of custom email templates.
- [okta_trusted_origin](./okta_trusted_origin) Supports the management of Okta Trusted Sources and Origins.
- [okta_user_base_schema](./okta_user_base_schema) Supports the management of Okta User Profile Attribute Schemas.
- [okta_user_schema](./okta_user_schema) Supports the management of Okta defined User Profile Attribute Schemas.
- [okta_user](./okta_user) Supports the management of Okta Users.
- [okta_users](./okta_users) Data source to retrieve a group of users.
- [okta_app_oauth_redirect_uri](./okta_app_oauth_redirect_uri) Supports decentralizing redirect uri config. Due to Okta's API not allowing this field to be null, you must set a redirect uri in your app, and ignore changes to this attribute. We follow TF best practices and detect config drift. The best case scenario is Okta makes this field nullable, and we can not detect config drift when this attr is not present.

## Deprecated Resources

- okta_identity_provider -- See okta_idp, okta_idp_social, and okta_idp_saml.
- okta_user_schemas -- See okta_user_schema.

## Notes

As resource fixtures are added, please be sure to only put a VALID config in each resource sub-directory. Intentionally invalid config for testing should stay in the test file.
