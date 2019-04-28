# Configuration Examples

Here lies the examples that will aid you on your Okta Terraform journey.

## Example Stacks

* [Okta and Cognito](./oidc-cognito-stack.tf) Example of using Okta OIDC application with a Cognito ID Provider to provide a serverless SPA access to AWS resources.

## Test Fixture Examples

Anything that lies underneath a resource directory is config we use as fixtures to our tests. We use the special trigger `replace_with_uuid` to trigger the fixture manager to insert a unique id. Eventually we will stand up a wiki for the provider but we don't quite have the manpower yet to do so.

## Resources & Data Sources

* [okta_saml_app](./okta_saml_app) Supports the management of Okta SAML Applications.
* [okta_oauth_app](./okta_oauth_app) Supports the management of Okta OIDC Applications.
* [okta_bookmark_app](./okta_bookmark_app) Supports the management Okta Bookmark Application.
* [okta_app](./okta_app) Generic Application data source.
* [okta_user](./okta_user) Supports the management of Okta Users.
* [okta_group](./okta_group) Supports the management of Okta Groups.
* [okta_group_rule](./okta_group_rule) Supports the management of Okta Group Rules.
* [okta_trusted_origin](./okta_trusted_origin) Supports the management of Okta Trusted Sources and Origins.
* [okta_user_schemas](./okta_user_schemas) Supports the management of Okta User Profile Attribute Schemas.
* [okta_auth_server](./okta_auth_server) Supports the management of Okta Authorization servers.
* [okta_auth_server_policy](./okta_auth_server_policy) Supports the management of Okta Authorization servers policies.
* [okta_auth_server_policy_rule](./okta_auth_server_policy_rule) Supports the management of Okta Authorization servers policy rules.
* [okta_auth_server_scope](./okta_auth_server_scope) Supports the management of Okta Authorization servers scopes.
* [okta_auth_server_claim](./okta_auth_server_claim) Supports the management of Okta Authorization servers claims.
* [okta_inline_hook](./okta_inline_hook) Supports the management of Okta Inline Hooks EA feature.
* [okta_idp](./okta_idp) Supports the management of Okta OIDC Identity Providers.
* [okta_social_idp](./okta_social_idp) Supports the management of Okta Social Identity Providers. Such as Google, Facebook, Microsoft, and LinkedIn.
* [okta_saml_idp](./okta_idp) Supports the management of Okta SAML Identity Providers.
* [okta_oauth_app_redirect_uri](./okta_oauth_app_redirect_uri) Supports decentralizing redirect uri config. Due to Okta's API not allowing this field to be null, you must set a redirect uri in your app, and ignore changes to this attribute. We follow TF best practices and detect config drift. The best case scenario is Okta makes this field nullable and we can not detect config drift when this attr is not present.

## Deprecated Resources

* okta_identity_provider -- See okta_idp, okta_social_idp, and okta_saml_idp.
* okta_user_schemas -- See okta_user_schema.

## Notes

As resource fixtures are added, please be sure to only put VALID config in each resource sub directory. Intentionally invalid config for testing should stay in the test file.
