---
page_title: "Resource: okta_rate_limiting"
description: |-
  Manages rate limiting settings.
  This resource allows you to configure the per-client rate limit settings for your Okta organization.
---

# Resource: okta_rate_limiting

Manages per-client rate limiting settings for your Okta organization.

This resource configures how Okta handles rate limiting on a per-client basis, allowing you to set a default mode and override settings for specific use cases.

## Example Usage

### Basic Usage

```terraform
resource "okta_rate_limiting" "example" {
  default_mode = "ENFORCE"
}
```

### With Use Case Overrides

```terraform
resource "okta_rate_limiting" "example" {
  default_mode = "ENFORCE"

  use_case_mode_overrides {
    login_page       = "PREVIEW"
    oauth2_authorize = "ENFORCE"
    oie_app_intent   = "DISABLE"
  }
}
```

## Argument Reference

### Required

- `default_mode` (String) - The default per-client rate limit mode. Valid values:
  - `ENFORCE` - Enforce limit and log per client (recommended)
  - `DISABLE` - Do nothing (not recommended)
  - `PREVIEW` - Log per client without enforcing limits

### Optional

- `use_case_mode_overrides` (Block) - A map of per-client rate limit use cases to the applicable mode. Overrides the `default_mode` property for the specified use cases. Supports the following attributes:
  - `login_page` (String) - Rate limit mode for the Okta hosted login page. Valid values: `ENFORCE`, `DISABLE`, `PREVIEW`.
  - `oauth2_authorize` (String) - Rate limit mode for OAuth2 authorization requests. Valid values: `ENFORCE`, `DISABLE`, `PREVIEW`.
  - `oie_app_intent` (String) - Rate limit mode for OIE app intent. Valid values: `ENFORCE`, `DISABLE`, `PREVIEW`.

## Attributes Reference

- `id` (String) - The ID of this resource (always `rate_limiting`).

## Import

Import is supported using the following syntax:

```shell
terraform import okta_rate_limiting.example .
```

~> **Note:** The import ID is a literal dot (`.`) since there is only one rate limiting configuration per Okta organization.
