---
layout: "okta"
page_title: "Okta: okta_user_base_schema"
sidebar_current: "docs-okta-resource-user-base-schema"
description: |-
  Manages a User Base Schema property.
---

# okta_user_base_schema

Manages a User Base Schema property.

This resource allows you to configure a base user schema property.

## Example Usage

```hcl
resource "okta_user_base_schema" "example" {
  index       = "customPropertyName"
  title       = "customPropertyName"
  type        = "string"
  master      = "OKTA"
}
```

## Argument Reference

The following arguments are supported:

* `index` - (Required) The property name.

* `title` - (Required) The property display name.

* `type` - (Required) The type of the schema property. It can be `"string"`, `"boolean"`, `"number"`, `"integer"`, `"array"`, or `"object"`.

* `required` - (Optional) Whether the property is required for this application's users.

* `pattern` - (Optional) Restrict `login` schema property to a particular regex.

* `permissions` - (Optional) Access control permissions for the property. It can be set to `"READ_WRITE"`, `"READ_ONLY"`, `"HIDE"`.

* `master` - (Optional) Master priority for the user schema property. It can be set to `"PROFILE_MASTER"` or `"OKTA"`.

## Attributes Reference

* `index` - ID of the user schema property.

## Import

User base schema property can be imported via the property index.

```
$ terraform import okta_user_base_schema.example <property name>
```
