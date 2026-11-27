---
page_title: "Resource: okta_event_hook"
description: |-
  Manages an Okta event hook.
---

# Resource: okta_event_hook

Manages an Okta [event hook](https://developer.okta.com/docs/api/openapi/okta-management/management/tag/EventHook/). An event hook sends an HTTP POST to an external service when any of the subscribed event types occur in your org.

New event hooks are created `ACTIVE` and `UNVERIFIED`. An `UNVERIFIED` event hook does not receive any events until it is verified; see `okta_event_hook_verification`.

## Example Usage

```terraform
resource "okta_event_hook" "example" {
  name        = "example"
  description = "Notifies an external service of user lifecycle events"

  channel {
    type    = "HTTP"
    version = "1.0.0"
    config {
      uri = "https://example.com/test"

      auth_scheme {
        type  = "HEADER"
        key   = "Authorization"
        value = "my-shared-secret"
      }

      headers {
        key   = "X-Other-Header"
        value = "my-header-value"
      }
    }
  }

  events {
    type = "EVENT_TYPE"
    items = [
      "user.lifecycle.create",
      "user.lifecycle.delete.initiated",
    ]
  }
}
```

## Schema

### Required

- `name` (String) Display name for the event hook.
- `channel` (Block) The channel that the event hook sends requests through. (see [below for nested schema](#nestedblock--channel))
- `events` (Block) The event types the hook subscribes to. (see [below for nested schema](#nestedblock--events))

### Optional

- `description` (String) Description of the event hook.

### Read-Only

- `id` (String) The ID of the event hook.
- `created_by` (String) The ID of the user who created the event hook.
- `status` (String) Status of the event hook: `ACTIVE` or `INACTIVE`.
- `verification_status` (String) Verification status of the event hook: `UNVERIFIED` or `VERIFIED`. `UNVERIFIED` event hooks don't receive any events.

<a id="nestedblock--channel"></a>
### Nested Schema for `channel`

Required:

- `type` (String) The channel type. Currently supports `HTTP`.
- `version` (String) Version of the channel. Currently the only supported version is `1.0.0`.
- `config` (Block) Endpoint configuration. (see [below for nested schema](#nestedblock--channel--config))

<a id="nestedblock--channel--config"></a>
### Nested Schema for `channel.config`

Required:

- `uri` (String) The external service endpoint called to execute the event hook handler. Must use HTTPS.

Optional:

- `auth_scheme` (Block) The authentication scheme used for this request. (see [below for nested schema](#nestedblock--channel--config--auth_scheme))
- `headers` (Block List) Optional list of key/value pairs for headers that can be sent with the request to the external service. (see [below for nested schema](#nestedblock--channel--config--headers))

Read-Only:

- `method` (String) The HTTP method of the Okta event hook request (`POST`).

<a id="nestedblock--channel--config--auth_scheme"></a>
### Nested Schema for `channel.config.auth_scheme`

Optional:

- `key` (String) The name for the authorization header, for example `Authorization`.
- `type` (String) The authentication scheme type. Currently supports `HEADER`.
- `value` (String) The header value. Okta never returns this value from the API.

<a id="nestedblock--channel--config--headers"></a>
### Nested Schema for `channel.config.headers`

Optional:

- `key` (String) The header name.
- `value` (String) The header value.

<a id="nestedblock--events"></a>
### Nested Schema for `events`

Required:

- `type` (String) The events object type. Currently supports `EVENT_TYPE`.
- `items` (List of String) The subscribed event types that trigger the event hook. For the eligible types, search the [Event Types catalog](https://developer.okta.com/docs/reference/api/event-types/#catalog) for `event-hook-eligible`.

Optional:

- `filter` (Block) An optional filter on specific event types. Event hook filters are a self-service Early Access feature. (see [below for nested schema](#nestedblock--events--filter))

<a id="nestedblock--events--filter"></a>
### Nested Schema for `events.filter`

Optional:

- `event_filter_map` (Block List) Maps a filter to an event type. (see [below for nested schema](#nestedblock--events--filter--event_filter_map))

Read-Only:

- `type` (String) The type of filter. Currently only supports `EXPRESSION_LANGUAGE`.

<a id="nestedblock--events--filter--event_filter_map"></a>
### Nested Schema for `events.filter.event_filter_map`

Optional:

- `event` (String) The filtered event type.
- `condition` (Block) The filter condition. (see [below for nested schema](#nestedblock--events--filter--event_filter_map--condition))

<a id="nestedblock--events--filter--event_filter_map--condition"></a>
### Nested Schema for `events.filter.event_filter_map.condition`

Optional:

- `expression` (String) The Okta Expression Language statement that filters the event type.

Read-Only:

- `version` (String) Internal field.

## Import

Import is supported using the following syntax:

```shell
terraform import okta_event_hook.example <event_hook_id>
```
