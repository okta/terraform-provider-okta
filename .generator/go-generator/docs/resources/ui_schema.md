---
page_title: "okta_ui_schema Resource - terraform-provider-okta"
description: |-
  Manages an Okta Ui Schema resource.
---

# okta_ui_schema (Resource)

Manages an Okta Ui Schema.

## Example Usage

```hcl
resource "okta_ui_schema" "example" {

  # Optional fields
  # button_label = "Example Button Label"
  # label = "Example Label"
  # format = "<format>"
  # scope = "<scope>"
  # type = "<type>"
}
```

## Schema

### Optional

- `button_label` (String) Specifies the button label for the `Submit` button at the bottom of the enrollment form
- `label` (String) Label name for the UI element
- `format` (String) Specifies how the input appears
- `scope` (String) Specifies the property bound to the input field. It must follow the format `#/properties/PROPERTY_NAME` where `PROPERTY_NAME` is a variable name for an attribute in `profile editor`.
- `type` (String) Specifies the relationship between this input element and `scope`. The `Control` value specifies that this input controls the value represented by `scope`.
- `label` (String) Specifies the label at the top of the enrollment form under the logo
- `type` (String) Specifies the type of layout

### Read-Only

- `id` (String) The unique identifier for the resource.
