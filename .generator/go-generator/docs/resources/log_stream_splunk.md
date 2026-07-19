---
page_title: "okta_log_stream_splunk Resource - terraform-provider-okta"
description: |-
  Manages an Okta Log Stream Splunk resource.
---

# okta_log_stream_splunk (Resource)

Manages an Okta Log Stream Splunk.

## Example Usage

```hcl
resource "okta_log_stream_splunk" "example" {
  type = "<type>"
  name = "Example Name"
  edition = "<edition>"
  host = "<host>"
  token = "<token>"
}
```

## Schema

### Required

- `type` (String) Discriminator field identifying the variant type. Must be set to \
- `name` (String) Unique name for the log stream object
- `edition` (String) Edition of the Splunk Cloud instance
- `host` (String) The domain name for your Splunk Cloud instance. Don't include `http` or `https` in the string. For example: `acme.splunkcloud.com`
- `token` (String) The HEC token for your Splunk Cloud HTTP Event Collector. The token value is set at object creation, but isn't returned.

### Read-Only

- `id` (String) The unique identifier for the resource.
