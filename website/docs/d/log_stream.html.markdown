---
layout: 'okta'
page_title: 'Okta: okta_log_stream'
sidebar_current: 'docs-okta-datasource-log-stream'
description: |-
  Gets Okta Log Stream.
---

# okta_log_stream

Use this data source to retrieve a log stream from Okta.

## Example Usage

```hcl
data "okta_log_stream" "example" {
  name = "Example Stream"
}
```

## Argument Reference

- `id` - (Optional) ID of the log stream to retrieve, conflicts with `name`.

- `name` - (Optional) Name of the log stream to retrieve, conflicts with `id`.

## Attributes Reference

- `id` - ID of the log stream.

- `name` - Name of the log stream.

- `type` - Type of the Log Stream.

- `status` - Log Stream Status - can either be ACTIVE or INACTIVE only.

- `settings` - Provider specific configuration.
