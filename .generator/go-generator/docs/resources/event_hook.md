---
page_title: "okta_event_hook Resource - terraform-provider-okta"
description: |-
  Manages an Okta Event Hook resource.
---

# okta_event_hook (Resource)

Manages an Okta Event Hook.

## Example Usage

```hcl
resource "okta_event_hook" "example" {
  uri = "https://example.com"
  type = "<type>"
  version = "<version>"
  items = "<items>"
  type = "<type>"
  name = "Example Name"

  # Optional fields
  # key = "<key>"
  # type = "<type>"
  # value = "<value>"
  # key = "<key>"
  # value = "<value>"
}
```

## Schema

### Required

- `uri` (String) The external service endpoint called to execute the event hook handler
- `type` (String) The channel type. Currently supports `HTTP`.
- `version` (String) Version of the channel. Currently the only supported version is `1.0.0`.
- `items` (List) The subscribed event types that trigger the event hook. When you register an event hook you need to specify which events you want to subscribe to. To see the list of event types currently eligible ...
- `type` (String) The events object type. Currently supports `EVENT_TYPE`.
- `name` (String) Display name for the event hook

### Optional

- `key` (String) The name for the authorization header
- `type` (String) The authentication scheme type. Currently only supports `HEADER`.
- `value` (String) The header value. This secret key is passed to your external service endpoint for security verification. This property is not returned in the response.
- `key` (String) The optional field or header name
- `value` (String) The value for the key
- `description` (String) Description of the event hook
- `expression` (String) The Okta Expression language statement that filters the event type
- `event` (String) The filtered event type

### Read-Only

- `id` (String) The unique identifier for the resource.
- `method` (String) The method of the Okta event hook request
- `created` (String) Timestamp of the event hook creation
- `created_by` (String) The ID of the user who created the event hook
- `version` (String) Internal field
- `type` (String) The type of filter. Currently only supports `EXPRESSION_LANGUAGE`
- `last_updated` (String) Date of the last event hook update
- `status` (String) Status of the event hook
- `verification_status` (String) Verification status of the event hook. `UNVERIFIED` event hooks won't receive any events.
