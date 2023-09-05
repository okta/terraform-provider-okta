---
layout: 'okta'
page_title: 'Okta: okta_rate_limiting'
sidebar_current: 'docs-okta-resource-rate-limit'
description: |-
  Manages rate limiting.
---

# okta_rate_limiting

This resource allows you to configure the client-based rate limit and rate limiting communications settings.

~> **WARNING:** This resource is available only when using a SSWS API token in the provider config, it is incompatible with OAuth 2.0 authentication.

~> **WARNING:** This resource makes use of an internal/private Okta API endpoint that could change without notice rendering this resource inoperable. 

## Example Usage

```hcl
resource "okta_rate_limiting" "example" {
  login                  = "ENFORCE"
  authorize              = "ENFORCE"
  communications_enabled = true
}
```

## Argument Reference

- `login` - (Required) Called when accessing the Okta hosted login page. Valid values: `"ENFORCE"` _(Enforce limit and 
log per client (recommended))_, `"DISABLE"` _(Do nothing (not recommended))_, `"PREVIEW"` _(Log per client)_.

- `authorize` - (Required) Called during authentication. Valid values: `"ENFORCE"` _(Enforce limit and
log per client (recommended))_, `"DISABLE"` _(Do nothing (not recommended))_, `"PREVIEW"` _(Log per client)_.

- `communications_enabled` - (Optional) Enable or disable rate limiting communications. By default, it is `true`.

## Import

Rate limit settings can be imported without any parameters.

```
$ terraform import okta_rate_limiting.example .
```
