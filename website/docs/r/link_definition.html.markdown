---
layout: 'okta'
page_title: 'Okta: okta_link_definition'
sidebar_current: 'docs-okta-resource-link-definition'
description: |-
    Manages the creation and removal of the link definitions.
---

# okta_link_definition

Link definition operations allow you to manage the creation and removal of the link definitions. If you remove a link 
definition, links based on that definition are unavailable. Note that this resource is immutable, thus can not be modified.

~> **NOTE:** Links reappear if you recreate the definition. However, Okta is likely to change this behavior so that links don't reappear. Don't rely on this behavior in production environments.

## Example Usage

```hcl
resource "okta_link_definition" "example" {
  primary_name           = "emperor"
  primary_title          = "Emperor"
  primary_description    = "Hereditary ruler of the Imperium and the Known Universe"
  associated_name        = "sardaukar"
  associated_title       = "Sardaukar"
  associated_description = "Elite military force member"
}
```

## Argument Reference

- `primary_name` - (Required) API name of the primary link.

- `primary_title` - (Required) Display name of the primary link.

- `primary_description` - (Required) Description of the primary relationship.

- `associated_name` - (Required) API name of the associated link.

- `associated_title` - (Required) Display name of the associated link.

- `associated_description` - (Required) Description of the associated relationship.

## Attributes Reference

- `id` - Name of the primary link.

## Import

Okta Link Definition can be imported via the Okta Primary Link Name.

```
$ terraform import okta_link_definition.example &#60;primary_name&#62;
```
