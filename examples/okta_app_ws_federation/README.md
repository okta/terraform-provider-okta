# okta_app_ws_federation $\pi$

This resource represents an Okta WS Federation Application in various configuration states. For more information see
the [API docs](https://developer.okta.com/docs/reference/api/apps/#add-ws-federation-application)

- Example of a custom WS-Federation app [can be found here](./custom.tf)
- Example of a custom WS-Federation app with attribute statements [can be found here](./custom_updated.tf)
- Example of an AWS pre-configured WS Fed app [can be found here](./preconfigured.tf)
- Example of WS Fed App data source [can be found here](./datasource.tf)

## Pre-configured Applications

There are some configuration options that cannot be configured on certain "pre-configured" OAuth applications due to
limitations in the Okta API.