# Configuration Examples

Here lies the examples that will aid you on your Okta Terraform journey.

## Example Stacks

* [Okta and Cognito](./oidc-cognito-stack.tf) Example of using Okta OIDC application with a Cognito ID Provider to provide a serverless SPA access to AWS resources.

## Test Fixture Examples

Anything that lies underneath a resource directory is config we use as fixtures to our tests. This means we pass our tests through a string formatter that replaces the format placeholders with the actual values but regardless seeing these might be helpful when trying to implement this provider. Just remember, it is not valid config until you replace or remove `%[1]d`.

## Resources

* [okta_saml_app](./okta_saml_app) Supports the management of Okta SAML Applications.
* [okta_oauth_app](./okta_oauth_app) Supports the management of Okta OIDC Applications.
* [okta_user](./okta_user) Supports the management of Okta Users.
* [okta_group](./okta_group) Supports the management of Okta Groups.
* [okta_trusted_origin](./okta_trusted_origin) Supports the management of Okta Trusted Sources and Origins.
* [okta_user_schemas](./okta_user_schemas) Supports the management of Okta User Profile Attribute Schemas.
* [okta_identity_provider](./okta_identity_provider) Supports the management of Okta Identity Provider.

## Notes

As resource fixtures are added, please be sure to only put VALID config in each resource sub directory. Intentionally invalid config for testing should stay in the test file.
