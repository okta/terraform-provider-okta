---
page_title: "Resource: okta_event_hook"
description: |-
  Manages an okta_event_hook resource.
---

# Resource: okta_event_hook

Manages an okta_event_hook resource.

## Example Usage

```terraform
resource "okta_event_hook" "example" {
  name = "example"

  channel {
    type    = "HTTP"
    version = "1.0.0"
    config {
      uri = "https://example.com/test"
    }
  }

  events {
    type  = "EVENT_TYPE"
    items = ["user.lifecycle.create", "user.lifecycle.delete.initiated"]
  }
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) Display name for the event hook
- `description` - (Optional) Description of the event hook
- `channel` - (Required) Channel
  - `type` - (Required) The channel type.
  - `version` - (Required) Version of the channel.
  - `config` - (Required) Config
    - `uri` - (Required) The external service endpoint called to execute the event hook handler
    - `method` - (Computed) The method of the Okta event hook request
    - `auth_scheme` - (Optional) The authentication scheme used for this request.
      - `key` - (Optional) The name for the authorization header
      - `type` - (Optional) The authentication scheme type
      - `value` - (Optional, Sensitive) The authentication secret
    - `headers` - (Optional) Additional request headers
      - `key` - (Optional) Header key
      - `value` - (Optional) Header value
- `events` - (Required) Events
  - `type` - (Required) The events object type. Currently, the only supported value is `EVENT_TYPE`.
  - `items` - (Required) The list of event types to deliver to this hook.
  - `filter` - (Optional) The optional filter defined on a specific event type
    - `type` - (Optional) The type of filter.
    - `event_filter_map` - (Optional) The object that maps the filter to the event type
      - `event` - (Optional) The filtered event type
      - `condition` - (Optional)
        - `expression` - (Optional) The Okta Expression language statement that filters the event type
        - `version` - (Optional)

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

- `id` - The ID of the resource.
- `created` - Timestamp of the event hook creation
- `created_by` - The ID of the user who created the event hook
- `last_updated` - Date of the last event hook update
- `status` - Status of the event hook
- `verification_status` - Verification status of the event hook.

## Import

Import is supported using the following syntax:

```shell
terraform import okta_event_hook.example <hook_id>
```
