---
page_title: "Data Source: okta_event_hook"
description: |-
  Retrieves an Okta event hook by ID.
---

# Data Source: okta_event_hook

Retrieves an Okta [event hook](https://developer.okta.com/docs/api/openapi/okta-management/management/tag/EventHook/) by ID.

## Example Usage

```terraform
data "okta_event_hook" "example" {
  id = "<event_hook_id>"
}

output "event_hook_uri" {
  value = data.okta_event_hook.example.channel.config.uri
}
```

## Schema

### Required

- `id` (String) The ID of the event hook.

### Read-Only

- `name` (String) Display name for the event hook.
- `description` (String) Description of the event hook.
- `created` (String) Timestamp of the event hook's creation (RFC 3339).
- `created_by` (String) The ID of the user who created the event hook.
- `last_updated` (String) Timestamp of the last event hook update (RFC 3339).
- `status` (String) Status of the event hook: `ACTIVE` or `INACTIVE`.
- `verification_status` (String) Verification status of the event hook: `UNVERIFIED` or `VERIFIED`.
- `channel` (Block) The channel that the event hook sends requests through. (see [below for nested schema](#nestedblock--channel))
- `events` (Block) The event types the hook subscribes to. (see [below for nested schema](#nestedblock--events))

<a id="nestedblock--channel"></a>
### Nested Schema for `channel`

Read-Only:

- `type` (String) The channel type, for example `HTTP`.
- `version` (String) Version of the channel, for example `1.0.0`.
- `config` (Block) Endpoint configuration. (see [below for nested schema](#nestedblock--channel--config))

<a id="nestedblock--channel--config"></a>
### Nested Schema for `channel.config`

Read-Only:

- `uri` (String) The external service endpoint called to execute the event hook handler.
- `method` (String) The HTTP method of the Okta event hook request.
- `headers` (List of String) Headers sent with the request to the external service. Each element is a JSON-encoded `{"key": ..., "value": ...}` object.
- `auth_scheme` (Block) The authentication scheme used for this request. (see [below for nested schema](#nestedblock--channel--config--auth_scheme))

<a id="nestedblock--channel--config--auth_scheme"></a>
### Nested Schema for `channel.config.auth_scheme`

Read-Only:

- `key` (String) The name for the authorization header.
- `type` (String) The authentication scheme type.
- `value` (String) The header value. Okta doesn't return this from the API, so it's always empty.

<a id="nestedblock--events"></a>
### Nested Schema for `events`

Read-Only:

- `type` (String) The events object type, for example `EVENT_TYPE`.
- `items` (List of String) The subscribed event types that trigger the event hook.
- `filter` (Block) The optional filter defined on specific event types. (see [below for nested schema](#nestedblock--events--filter))

<a id="nestedblock--events--filter"></a>
### Nested Schema for `events.filter`

Read-Only:

- `type` (String) The type of filter, for example `EXPRESSION_LANGUAGE`.
- `event_filter_map` (List of String) Maps filters to event types. Each element is a JSON-encoded object with `event` and `condition` fields.
