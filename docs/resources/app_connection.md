---
page_title: "Resource: okta_app_connection"
description: |-
  Manages Okta App Connection configurations for provisioning.
---

# Resource: okta_app_connection

Manages Okta App Connection configurations for provisioning. This resource allows you to configure and manage provisioning connections for applications, including authentication schemes and connection activation/deactivation.

## Example Usage

### TOKEN Authentication

```terraform
resource "okta_app_connection" "token_example" {
  id       = "0oa1234567890abcdef"
  base_url = "https://api.example.com/scim/v2"
  action   = "activate"
  
  profile {
    auth_scheme = "TOKEN"
    token       = "your-bearer-token-here"
  }
}
```

### OAUTH2 Authentication

```terraform
resource "okta_app_connection" "oauth2_example" {
  id       = "0oa1234567890abcdef"
  action   = "activate"
  
  profile {
    auth_scheme = "OAUTH2"
    client_id   = "oauth2-client-id"
    
    settings {
      admin_username = "admin@example.com"
      admin_password = "secure-password"
    }
    
    signing {
      rotation_mode = "MANUAL"
    }
  }
}
```

### Deactivated Connection

```terraform
resource "okta_app_connection" "deactivated_example" {
  id       = "0oa1234567890abcdef"
  base_url = "https://api.example.com/scim/v2"
  action   = "deactivate"
  
  profile {
    auth_scheme = "TOKEN"
    token       = "your-bearer-token-here"
  }
}
```

## Schema

### Required

- `id` (String) The application ID for which to configure the provisioning connection.
- `base_url` (String) The base URL for the provisioning connection (typically the SCIM endpoint).
- `action` (String) The action to perform on the connection. Valid values are `activate` or `deactivate`.
- `profile` (Block) Profile configuration for the app connection. (see [below for nested schema](#nestedblock--profile))

### Read-Only

- `status` (String) The current status of the app connection. Values can be `ENABLED`, `DISABLED`, or `UNKNOWN`.

<a id="nestedblock--profile"></a>
### Nested Schema for `profile`

#### Required

- `auth_scheme` (String) Authentication scheme for the provisioning connection. Valid values are `TOKEN` or `OAUTH2`.

#### Optional

- `token` (String, Sensitive) Authentication token. Required when `auth_scheme` is `TOKEN`.
- `client_id` (String) OAuth2 client ID. Required when `auth_scheme` is `OAUTH2`.
- `signing` (Block) The signing key rotation setting.. Only used for the Okta Org2Org (okta_org2org) app. (see [below for nested schema](#nestedblock--profile--signing))
- `settings` (Block) Settings required for the Microsoft Office 365 provisioning connection. (see [below for nested schema](#nestedblock--profile--settings))

<a id="nestedblock--profile--signing"></a>
### Nested Schema for `profile.signing`

#### Optional

- `rotation_mode` (String) The signing key rotation setting for the provisioning connection.
    - `AUTO` Okta manages key rotation for the provisioning connection.
    - `MANUAL` You need to rotate the keys for your provisioning connection manually based on your own schedule.

<a id="nestedblock--profile--settings"></a>
### Nested Schema for `profile.settings`

#### Optional

- `admin_username` (String) Microsoft Office 365 global administrator username.
- `admin_password` (String, Sensitive) Microsoft Office 365 global administrator password.

## Import

App connections can be imported using the application ID:

```shell
terraform import okta_app_connection.example <application_id>
```