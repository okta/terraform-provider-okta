---
layout: "okta"
page_title: "Okta: okta_event_hook"
sidebar_current: "docs-okta-resource-event-hook"
description: |-
  Creates an event hook.
---

# okta_event_hook

Creates an event hook.

This resource allows you to create and configure an event hook.

## Example Usage

```hcl
resource "okta_event_hook" "example" {
  name    = "example"
  events  = [
    "user.lifecycle.create",
    "user.lifecycle.delete.initiated",
  ]

  channel = {
    type    = "HTTP"
    version = "1.0.0"
    uri     = "https://example.com/test"
  }

  auth = {
    type  = "HEADER"
    key   = "Authorization"
    value = "123"
  }
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The event hook display name.

- `events` - (Required) The events that will be delivered to this hook. [See here for a list of supported events](https://developer.okta.com/docs/reference/api/event-types/?q=event-hook-eligible).

- `headers` - (Optional) Map of headers to send along in event hook request.

- `auth` - (Optional) Authentication required for event hook request.

  - `key` - (Required) Key to use for authentication, usually the header name, for example `"Authorization"`.
  - `value` - (Required) Authentication secret.
  - `type` - (Optional) Auth type. Currently, the only supported type is `"HEADER"`.

- `channel` - (Required) Details of the endpoint the event hook will hit.
  - `version` - (Required) The version of the channel. The currently-supported version is `"1.0.0"`.
  - `uri` - (Required) The URI the hook will hit.
  - `type` - (Optional) The type of hook to trigger. Currently, the only supported type is `"HTTP"`.

## Attributes Reference

- `id` - The ID of the event hooks.

## Import

An event hook can be imported via the Okta ID.

```
$ terraform import okta_event_hook.example <hook id>
```
