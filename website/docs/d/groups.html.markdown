---
layout: "okta"
page_title: "Okta: okta_groups"
sidebar_current: "docs-okta-datasource-groups"
description: |-
  Get a list of groups from Okta.
---

# okta_groups

Use this data source to retrieve a list of groups from Okta.

## Example Usage

```hcl
data "okta_groups" "example" {
  q = "Engineering - "
}
```

## Arguments Reference

* `q` - (Required) Searches the name property of groups for matching value

## Attributes Reference

* `groups` - colletion of groups retrieved from Okta with the following properties.
  * `name` - Group name.
  * `description` - Group description.
