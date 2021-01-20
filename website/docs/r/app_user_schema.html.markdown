---
layout: 'okta'
page_title: 'Okta: okta_app_user_schema'
sidebar_current: 'docs-okta-resource-app-user-schema'
description: |-
  Creates an Application User Schema property.
---

# okta_app_user_schema

Creates an Application User Schema property.

This resource allows you to create and configure a custom user schema property and associate it with an application.

## Example Usage

```hcl
resource "okta_app_user_schema" "example" {
  app_id      = "<app id>"
  index       = "customPropertyName"
  title       = "customPropertyName"
  type        = "string"
  description = "My custom property name"
  master      = "OKTA"
  scope       = "SELF"
}
```

## Argument Reference

The following arguments are supported:

- `app_id` - (Required) The Application's ID the user custom schema property should be assigned to.

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

- `master` - (Optional) Master priority for the user schema property. It can be set to `"PROFILE_MASTER"` or `"OKTA"`.

- `external_name` - (Optional) External name of the user schema property.

- `external_namespace` - (Optional) External namespace of the user schema property.

- `union` - (Optional) Used to assign attribute group priority. Can not be set to 'true' if `scope` is set to Individual level.

## Attributes Reference

- `app_id` - ID of the application the user property is associated with.

- `index` - ID of the user schema property.

## Import

App user schema property can be imported via the property index and app id.

```
$ terraform import okta_app_user_schema.example <app id>/<property name>
```
