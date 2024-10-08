---
page_title: "Resource: okta_link_definition"
description: |-
  Manages the creation and removal of the link definitions.
  Link definition operations allow you to manage the creation and removal of the link definitions. If you remove a link
  definition, links based on that definition are unavailable. Note that this resource is immutable, thus can not be modified.
  ~> NOTE: Links reappear if you recreate the definition. However, Okta is likely to change this behavior so that links don't reappear. Don't rely on this behavior in production environments.
---

# Resource: okta_link_definition

Manages the creation and removal of the link definitions.
		
Link definition operations allow you to manage the creation and removal of the link definitions. If you remove a link 
definition, links based on that definition are unavailable. Note that this resource is immutable, thus can not be modified.
~> **NOTE:** Links reappear if you recreate the definition. However, Okta is likely to change this behavior so that links don't reappear. Don't rely on this behavior in production environments.

## Example Usage

```terraform
resource "okta_link_definition" "example" {
  primary_name           = "emperor"
  primary_title          = "Emperor"
  primary_description    = "Hereditary ruler of the Imperium and the Known Universe"
  associated_name        = "sardaukar"
  associated_title       = "Sardaukar"
  associated_description = "Elite military force member"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `associated_description` (String) Description of the associated relationship.
- `associated_name` (String) API name of the associated link.
- `associated_title` (String) Display name of the associated link.
- `primary_description` (String) Description of the primary relationship.
- `primary_name` (String) API name of the primary link.
- `primary_title` (String) Display name of the primary link.

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import okta_link_definition.example <primary_name>
```
