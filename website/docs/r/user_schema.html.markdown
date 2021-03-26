---
layout: 'okta'
page_title: 'Okta: okta_user_schema'
sidebar_current: 'docs-okta-resource-user-schema'
description: |-
  Creates a User Schema property.
---

# okta_user_schema

Creates a User Schema property.

This resource allows you to create and configure a custom user schema property.

## Example Usage

```hcl
resource "okta_user_schema" "example" {
  index       = "customPropertyName"
  title       = "customPropertyName"
  type        = "string"
  description = "My custom property name"
  master      = "OKTA"
  scope       = "SELF"
  user_type   = "${data.okta_user_type.example.id}"
}
```

## Argument Reference

The following arguments are supported:

- `index` - (Required) The property name.

- `title` - (Required) The display name.

- `type` - (Required) The type of the schema property. It can be `"string"`, `"boolean"`, `"number"`, `"integer"`, `"array"`, or `"object"`.

- `enum` - (Optional) Array of values a primitive property can be set to. See `array_enum` for arrays.

- `one_of` - (Optional) Array of maps containing a mapping for display name to enum value.

  - `const` - (Required) value mapping to member of `enum`.
  - `title` - (Required) display name for the enum value.

- `description` - (Optional) The description of the user schema property.

- `required` - (Optional) Whether the property is required for this application's users.

- `min_length` - (Optional) The minimum length of the user property value. Only applies to type `"string"`.

- `max_length` - (Optional) The maximum length of the user property value. Only applies to type `"string"`.

- `scope` - (Optional) determines whether an app user attribute can be set at the Individual or Group Level.

- `array_type` - (Optional) The type of the array elements if `type` is set to `"array"`.

- `array_enum` - (Optional) Array of values that an array property's items can be set to.

- `array_one_of` - (Optional) Display name and value an enum array can be set to.

  - `const` - (Required) value mapping to member of `enum`.
  - `title` - (Required) display name for the enum value.

- `permissions` - (Optional) Access control permissions for the property. It can be set to `"READ_WRITE"`, `"READ_ONLY"`, `"HIDE"`.

- `master` - (Optional) Master priority for the user schema property. It can be set to `"PROFILE_MASTER"`, `"OVERRIDE"` or `"OKTA"`.

- `master_override_priority` - (Optional) Prioritized list of profile sources (required when `master` is `"OVERRIDE"`).
  - `type` - (Optional) - Type of profile source.
  - `value` - (Required) - ID of profile source.

- `external_name` - (Optional) External name of the user schema property.

- `external_namespace` - (Optional) External name of the user schema property.

- `unique` - (Optional) Whether the property should be unique. It can be set to `"UNIQUE_VALIDATED"` or `"NOT_UNIQUE"`.

- `user_type` - (Optional) User type ID

## Attributes Reference

- `index` - ID of the user schema property.

## Import

User schema property of default user type can be imported via the property index.

```
$ terraform import okta_user_schema.example <index>
```

User schema property of custom user type can be imported via user type id and property index

```
$ terraform import okta_user_schema.example <user type id>.<index>
```
