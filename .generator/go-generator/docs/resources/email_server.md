---
page_title: "okta_email_server Resource - terraform-provider-okta"
description: |-
  Manages an Okta Email Server resource.
---

# okta_email_server (Resource)

Manages an Okta Email Server.

## Example Usage

```hcl
resource "okta_email_server" "example" {
  alias = "<alias>"
  auth_type = "<auth_type>"
  enabled = true
  host = "<host>"
  port = 0
  username = "Example Username"
}
```

## Schema

### Required

- `alias` (String) Human-readable name for your SMTP server
- `auth_type` (String) <x-lifecycle-container><x-lifecycle class='ea'></x-lifecycle> <x-lifecycle class='oie'></x-lifecycle></x-lifecycle-container>The authentication type that's used by your SMTP server
- `enabled` (Boolean) If `true`, all email traffic is routed through your SMTP server
- `host` (String) Hostname or IP address of your SMTP server
- `port` (Number) Port number of your SMTP server
- `username` (String) Username that's used to access your SMTP server

### Read-Only

- `id` (String) The unique identifier for the resource.
