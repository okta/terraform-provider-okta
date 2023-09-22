---
layout: 'okta'
page_title: 'Okta: okta_app_oauth_role_assignment'
sidebar_current: 'docs-okta-resource-okta-app-oauth-role-assignment'
description: |-
  Manages assignment of an admin role to an OAuth application
---

# okta_app_oauth_role_assignment

Manages assignment of an admin role to an OAuth application.

This resource allows you to assign an Okta admin role to a OAuth service application. This requires the Okta tenant feature flag for this function to be enabled.

## Example Usage

Standard Role:

```hcl
resource "okta_app_oauth" "test" {
  label          = "test"
  type           = "service"
  response_types = ["token"]
  grant_types    = ["client_credentials"]
  jwks_uri       = "https://example.com"
}

resource "okta_app_oauth_role_assignment" "test" {
  client_id = okta_app_oauth.test.client_id
  type      = "HELP_DESK_ADMIN"
}
```

Custom Role:

```hcl
resource "okta_app_oauth" "test" {
  label          = "test"
  type           = "service"
  response_types = ["token"]
  grant_types    = ["client_credentials"]
  jwks_uri       = "https://example.com"
}

resource "okta_admin_role_custom" "test" {
  label       = "test"
  description = "testing, testing"
  permissions = ["okta.apps.assignment.manage", "okta.users.manage", "okta.apps.manage"]
}

resource "okta_resource_set" "test" {
  label       = "test"
  description = "testing, testing"
  resources = [
    format("%s/api/v1/users", "https://example.okta.com"),
    format("%s/api/v1/apps", "https://example.okta.com")
  ]
}

resource "okta_app_oauth_role_assignment" "test" {
  client_id    = okta_app_oauth.test.client_id
  type         = "CUSTOM"
  role         = okta_admin_role_custom.test.id
  resource_set = okta_resource_set.test.id
}
```

## Argument Reference

The following arguments are supported:

- `client_id` - (Required) Client ID for the role to be assigned to

- `type` - (Required) Role type to assign. This can be one of the standard Okta roles, such as `HELP_DESK_ADMIN` or `CUSTOM`. Using custom requires the `resource_set` and `role` attributes to be set.

- `resource_set` - (Optional) Resource set for the custom role to assign, must be the ID of the created resource set.

- `role` - (Optional) Custom Role ID

## Attribute Reference

- `id` - Role Assignment ID

- `status` - Status of the role assignment

- `label` - Label of the role assignment

## Import

Not implemented
