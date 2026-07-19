---
page_title: "okta_log_stream_aws Resource - terraform-provider-okta"
description: |-
  Manages an Okta Log Stream Aws resource.
---

# okta_log_stream_aws (Resource)

Manages an Okta Log Stream Aws.

## Example Usage

```hcl
resource "okta_log_stream_aws" "example" {
  type = "<type>"
  name = "Example Name"
  account_id = "<account-id>"
  event_source_name = "Example Event Source Name"
  region = "<region>"
}
```

## Schema

### Required

- `type` (String) Discriminator field identifying the variant type. Must be set to \
- `name` (String) Unique name for the log stream object
- `account_id` (String) Your AWS account ID
- `event_source_name` (String) An alphanumeric name (no spaces) to identify this event source in AWS EventBridge
- `region` (String) The destination AWS region where your event source is located

### Read-Only

- `id` (String) The unique identifier for the resource.
