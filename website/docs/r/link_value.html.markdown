---
layout: 'okta'
page_title: 'Okta: okta_link_value'
sidebar_current: 'docs-okta-resource-link-value'
description: |-
  Manages users relationships
---

# okta_link_value

Link value operations allow you to create relationships between primary and associated users.

## Example Usage

```hcl
resource "okta_link_definition" "padishah" {
  primary_name           = "emperor"
  primary_title          = "Emperor"
  primary_description    = "Hereditary ruler of the Imperium and the Known Universe"
  associated_name        = "sardaukar"
  associated_title       = "Sardaukar"
  associated_description = "Elite military force member"
}

resource "okta_user" "emperor" {
  first_name = "Shaddam"
  last_name  = "Corrino IV"
  login      = "shaddam.corrino.iv@salusa-secundus.planet"
  email      = "shaddam.corrino.iv@salusa-secundus.planet"
}

resource "okta_user" "sardaukars" {
  count      = 5
  first_name = "Amrit"
  last_name  = "Sardaukar_${count.index}"
  login      = "amritsardaukar_${count.index}@salusa-secundus.planet"
  email      = "amritsardaukar_${count.index}@salusa-secundus.planet"
}

resource "okta_link_value" "example" {
  primary_name        = okta_link_definition.padishah.primary_name
  primary_user_id     = okta_user.emperor.id
  associated_user_ids = [
    okta_user.sardaukars[0].id,
    okta_user.sardaukars[1].id,
    okta_user.sardaukars[2].id,
    okta_user.sardaukars[3].id,
    okta_user.sardaukars[4].id,
  ]
}
```

## Argument Reference

- `primary_name` - (Required) Name of the `primary` relationship being assigned.

- `primary_user_id` - (Required) User ID to be assigned to `primary` for the `associated` user in the specified relationship.

- `associated_user_ids` - (Optional) Set of User IDs or login values of the users to be assigned the 'associated' relationship.

## Attributes Reference

- `id` - ID of this resource in `primary_name/primary_user_id` format.

## Import

Okta Link Value can be imported via Primary Name and Primary User ID.

```
$ terraform import okta_link_value.example <primary_name>/<primary_user_id>
```
