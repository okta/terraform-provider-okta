# okta_app_saml

This resource represents an Okta SAML Application in various configuration states. For more information see the [API docs](https://developer.okta.com/docs/api/resources/apps#add-custom-saml-application)

- Example of a custom SAML app [can be found here](./custom_saml_app.tf)
- Example of a custom SAML app with attribute statements [can be found here](./custom_saml_app_all_fields.tf)
- Example of an AWS preconfigured SAML app [can be found here](./saml_app_with_groups_and_users.tf)
- Example of SAML App data source [can be found here](./datasource.tf)

## Preconfigured Applications

There are some configuration options that cannot be configured on certain "preconfigured" OAuth applications due to limitations in the Okta API.
