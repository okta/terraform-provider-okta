---
page_title: "okta_behavior_velocity Resource - terraform-provider-okta"
description: |-
  Manages an Okta Behavior Velocity resource.
---

# okta_behavior_velocity (Resource)

Manages an Okta Behavior Velocity.

## Example Usage

```hcl
resource "okta_behavior_velocity" "example" {
  type = "<type>"
  name = "Example Name"
  velocity_kph = 0

  # Optional fields
  # status = "ACTIVE"
}
```

## Schema

### Required

- `type` (String) Discriminator field identifying the variant type. Must be set to \
- `name` (String) Name of the Behavior Detection Rule
- `velocity_kph` (Number) VelocityKph

### Optional

- `status` (String) Status

### Read-Only

- `id` (String) The unique identifier for the resource.
- `link` (String) Specifies link relations (see [Web Linking](https://www.rfc-editor.org/rfc/rfc8288)) available using the [JSON Hypertext Application Language](https://datatracker.ietf.org/doc/html/draft-kelly-json...
- `created` (String) Timestamp when the Behavior Detection Rule was created
- `last_updated` (String) Timestamp when the Behavior Detection Rule was last modified
