---
page_title: "okta_behavior_anomalous_ip Resource - terraform-provider-okta"
description: |-
  Manages an Okta Behavior Anomalous Ip resource.
---

# okta_behavior_anomalous_ip (Resource)

Manages an Okta Behavior Anomalous Ip.

## Example Usage

```hcl
resource "okta_behavior_anomalous_ip" "example" {
  type = "<type>"
  name = "Example Name"

  # Optional fields
  # max_events_used_for_evaluation = 0
  # min_events_needed_for_evaluation = 0
  # status = "ACTIVE"
}
```

## Schema

### Required

- `type` (String) Discriminator field identifying the variant type. Must be set to \
- `name` (String) Name of the Behavior Detection Rule

### Optional

- `max_events_used_for_evaluation` (Number) MaxEventsUsedForEvaluation
- `min_events_needed_for_evaluation` (Number) MinEventsNeededForEvaluation
- `status` (String) Status

### Read-Only

- `id` (String) The unique identifier for the resource.
- `link` (String) Specifies link relations (see [Web Linking](https://www.rfc-editor.org/rfc/rfc8288)) available using the [JSON Hypertext Application Language](https://datatracker.ietf.org/doc/html/draft-kelly-json...
- `created` (String) Timestamp when the Behavior Detection Rule was created
- `last_updated` (String) Timestamp when the Behavior Detection Rule was last modified
