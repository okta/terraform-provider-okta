---
layout: 'okta'
page_title: 'Okta: okta_org_support'
sidebar_current: 'docs-okta-resource-okta-admin-role-targets'
description: |-
  Manages Okta Support access your org
---

# okta_org_support

This resource allows you to temporarily allow Okta Support to access your org as an administrator. By default,
access will be granted for eight hours. Removing this resource will revoke Okta Support access to your org.

## Example Usage

```hcl
resource "okta_org_support" "example" {
  extend_by = 1
}
```

## Argument Reference

- `extend_by` - (Optional) Number of days the support should be extended by in addition to the standard eight hours.

## Attributes Reference

- `status` - Status of Okta Support

- `expiration` - Expiration of Okta Support

## Import

This resource does not support importing.
