---
layout: "okta"
page_title: "Okta: okta_everyone_group"
sidebar_current: "docs-okta-datasource-everyone-group"
description: |-
  Get the Everyone group from Okta.
---

# okta_everyone_group

Use this data source to retrieve the collaborators for a given repository.

## Example Usage

```hcl
data "okta_everyone_group" "example" {}
```

## Attributes Reference

 * `id` - `id` of application.

 * `label` - `label` of application.

 * `description` - `description` of application.

 * `name` - `name` of application.

 * `status` - `status` of application.
