---
layout: 'okta'
page_title: 'Okta: okta_app_oauth_api_scope'
sidebar_current: 'docs-okta-resource-okta-app-oauth-api-scope'
description: |-
  Manages API scopes for OAuth applications.
---

# okta_app_oauth_api_scope

Manages API scopes for OAuth applications.

This resource allows you to grant or revoke API scopes for OAuth2 applications within your organization.

```
Note: you have to create an application before using this resource.
```

## Example Usage

```hcl
resource "okta_app_oauth_api_scope" "example" {
  app_id = "<application_id>"
  issuer = "<your org domain>"
  scopes = ["okta.users.read", "okta.users.manage"]
}
```

## Argument Reference

The following arguments are supported:

- `app_id` - (Required) ID of the application.

- `issuer` - (Required) The issuer of your Org Authorization Server, your Org URL.

- `scopes` - (Required) List of scopes for which consent is granted.

## Import

OAuth API scopes can be imported via the Okta Application ID.

```
$ terraform import okta_app_oauth_api_scope.example <app id>
```
