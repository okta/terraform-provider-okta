---
layout: 'okta' 
page_title: 'Okta: okta_event_hook_verification' 
sidebar_current: 'docs-okta-resource-event-hook-verification'
description: |-
  Verifies the Event Hook.
---

# okta_event_hook_verification

Verifies the Event Hook. The resource won't be created unless the URI provided in the event hook returns a valid
JSON object with verification. See [Event Hooks](https://developer.okta.com/docs/concepts/event-hooks/#one-time-verification-request)
documentation for details.

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

resource "okta_event_hook_verification" "example" {
  event_hook_id = okta_event_hook.example.id
}
```

## Argument Reference

The following arguments are supported:

- `event_hook_id` - (Required) Event Hook ID.

## Import

This resource does not support importing.
