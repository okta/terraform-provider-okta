---
page_title: "Data Source: okta_app_connection"
description: |-
  Retrieves the default provisioning connection for an app.
---

# Data Source: okta_app_connection

Retrieves the default provisioning connection for an app.

## Example Usage

### Basic Usage

```terraform
data "okta_app_connection" "example" {
  id = "0oa1234567890abcdef"
}

# Use the retrieved data
output "connection_status" {
  value = data.okta_app_connection.example.status
}
```

## Schema

### Required

- `id` (String) The application ID for which to retrieve the provisioning connection information.

### Read-Only

- `status` (String) Provisioning connection status.
- `auth_scheme` (String) A token is used to authenticate with the app. This property is only returned for the TOKEN authentication scheme.
- `base_url` (String) The base URL for the provisioning connection.
- `profile` (Block) Profile information for the app connection. (see [below for nested schema](#nestedblock--profile))

<a id="nestedblock--profile"></a>
### Nested Schema for `profile`

#### Read-Only

- `auth_scheme` (String) Defines the method of authentication.