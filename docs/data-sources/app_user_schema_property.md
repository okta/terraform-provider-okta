---
page_title: "okta_app_user_schema_property Data Source - terraform-provider-okta"
subcategory: ""
description: |-
  Gets an application user schema property.
---

# okta_app_user_schema_property (Data Source)

Gets an application user schema property.

This data source allows you to retrieve information about an existing application user schema property without managing it in Terraform.

~> **Note:** App user schema properties may be automatically created by Okta when provisioning features are enabled on an application (e.g., `PUSH_NEW_USERS`, `PUSH_PROFILE_UPDATES`). Common auto-created properties include `userName`, `email`, `givenName`, `familyName`, `displayName`, `title`, `department`, and `manager`. This data source is useful for referencing these properties without bringing them under Terraform management.

## Example Usage

### Query a Custom Property

```terraform
data "okta_app_saml" "example" {
  label = "My SAML App"
}

data "okta_app_user_schema_property" "department" {
  app_id = data.okta_app_saml.example.id
  index  = "department"
}

output "department_type" {
  value = data.okta_app_user_schema_property.department.type
}
```

### Query an Auto-Created Property

```terraform
# When provisioning is enabled, Okta creates standard properties
data "okta_app_user_schema_property" "username" {
  app_id = okta_app_saml.example.id
  index  = "userName"  # Auto-created when PUSH_NEW_USERS is enabled
}

# Reference the property in other resources
resource "okta_profile_mapping" "example" {
  source_id = okta_app_saml.example.id
  target_id = data.okta_user_type.default.id
  
  mappings {
    id         = data.okta_app_user_schema_property.username.id
    expression = "user.login"
  }
}
```

### Conditional Logic Based on Property Existence

```terraform
# Use try() to handle properties that may not exist
locals {
  department_exists = try(data.okta_app_user_schema_property.department.id, null) != null
}

resource "okta_profile_mapping" "conditional" {
  count = local.department_exists ? 1 : 0
  # ... configuration
}
```

## Argument Reference

The following arguments are supported:

* `app_id` - (Required) The Application's ID the user schema property is associated with.
* `index` - (Required) Subschema unique string identifier. This is the property name/key.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier for the schema property in format `<app_id>/<index>`.
* `title` - The property's display title.
* `type` - The property's data type. Possible values: `string`, `boolean`, `number`, `integer`, `array`, or `object`.
* `description` - The description of the user schema property.
* `required` - Whether the property is required.
* `min_length` - The minimum length of the user property value. Only applies to type `string`.
* `max_length` - The maximum length of the user property value. Only applies to type `string`.
* `permissions` - Access control permissions for the property. Possible values: `READ_WRITE`, `READ_ONLY`, `HIDE`.
* `master` - Master priority for the user schema property. Possible values: `PROFILE_MASTER`, `OKTA`.
* `scope` - Determines whether an app user attribute can be set at the Personal `SELF` or Group `NONE` level.
* `enum` - Array of values that a primitive property can be set to.
* `one_of` - Array of maps containing a mapping for display name to enum value.
* `external_name` - External name of the user schema property.
* `external_namespace` - External namespace of the user schema property.
* `unique` - Whether the property should be unique. Possible values: `UNIQUE_VALIDATED`, `NOT_UNIQUE`.
* `union` - (For array types) Whether attribute value is determined by group priority (false) or combines values across groups (true).
* `array_type` - The type of the array elements if `type` is `array`.
* `array_enum` - Array of values that an array property's items can be set to.
* `array_one_of` - Array of maps containing display names and values for array enum items.

## Common Auto-Created Properties

When provisioning features are enabled on an application, Okta typically creates these properties:

| Property | Type | Description |
|----------|------|-------------|
| `userName` | string | User's username |
| `email` | string | User's email address |
| `givenName` | string | User's first name |
| `familyName` | string | User's last name |
| `middleName` | string | User's middle name |
| `displayName` | string | User's display name |
| `title` | string | User's job title |
| `department` | string | User's department |
| `manager` | string | User's manager |
| `employeeNumber` | string | User's employee number |

The exact properties created depend on the application type and provisioning configuration.
