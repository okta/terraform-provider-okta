---
layout: 'okta'
page_title: 'Okta: okta_resource_set'
sidebar_current: 'docs-okta-resource-okta-resource-set'
description: |-
  Manages Resource Sets as custom collections of resources.
---

# okta_resource_set

This resource allows the creation and manipulation of Okta Resource Sets as custom collections of Okta resources. You can use 
Okta Resource Sets to assign Custom Roles to administrators who are scoped to the designated resources.
The `resources` field supports the following:
 - Apps
 - Groups
 - All Users within a Group
 - All Users within the org
 - All Groups within the org
 - All Apps within the org
 - All Apps of the same type

~> **NOTE:** This an `Early Access` feature.

## Example Usage

```hcl
locals {
  org_url = "https://mycompany.okta.com"
}

resource "okta_resource_set" "test" {
  label       = "UsersAppsAndGroups"
  description = "All the users, app and groups"
  resources   = [
    format("%s/api/v1/users", local.org_url),
    format("%s/api/v1/apps", local.org_url),
    format("%s/api/v1/groups", local.org_url)
  ]
}
```

## Argument Reference

- `label` - (Required) Unique name given to the Resource Set.

- `description` - (Required) A description of the Resource Set.

- `resources` - (Optional) The endpoints that reference the resources to be included in the new Resource Set. At least one
  endpoint must be specified when creating resource set.

## Attributes Reference

- `id` - ID of the resource set.

## Import

Okta Resource Set can be imported via the Okta ID.

```
$ terraform import okta_resource_set.example &#60;resource_set_id&#62;
```
