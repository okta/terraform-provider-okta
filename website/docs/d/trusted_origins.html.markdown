---
layout: 'okta'
page_title: 'Okta: okta_trusted_origins'
sidebar_current: 'docs-okta-datasource-trusted-origins'
description: |-
  Get List of Trusted Origins using filters.
---

# okta_trusted_origins

This resource allows you to retrieve a list of trusted origins from Okta.

## Example Usage

```hcl
data "okta_trusted_origins" "all" {
}
```

## Argument Reference

The following arguments are supported:

- `filter` - (Optional) Filter criteria (will be URL-encoded by the provider). See [Filtering](https://developer.okta.com/docs/reference/core-okta-api/#filter) for more information on the expressions used in filtering.

## Attributes Reference

- `id` - The ID of the Trusted Origin.

- `trusted_origins`
  - `active` - Whether the Trusted Origin is active or not - can only be issued post-creation
  - `name` - Unique name for this trusted origin.
  - `origin` - Unique origin URL for this trusted origin.
  - `scopes` - Scopes of the Trusted Origin
