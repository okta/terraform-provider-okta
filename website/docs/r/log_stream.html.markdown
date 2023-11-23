---
layout: 'okta'
page_title: 'Okta: okta_log_stream'
sidebar_current: 'docs-okta-resource-log-stream'
description: |-
  Creates an Okta Log Stream.
---

# okta_log_stream

Creates an Okta Log Stream.

This resource allows you to create and configure an Okta Log Stream.

## Example Usage - AWS EventBridge

```hcl
resource "okta_log_stream" "example" {
  name     = "EventBridge Log Stream"
  type     = "aws_eventbridge"
  status   = "ACTIVE"
  settings {
    account_id = "123456789012"
    region     = "us-north-1"
    event_source_name = "okta_log_stream"
  }
}
```

## Example Usage - Splunk Event Collector

```hcl
resource "okta_log_stream" "example" {
  name              = "Splunk log Stream"
  type              = "splunk_cloud_logstreaming"
  status            = "ACTIVE"
  settings {
    host = "acme.splunkcloud.com"
    edition = "gcp"
    token = "YOUR_HEC_TOKEN"
  }
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) Name of the Log Stream Resource.

- `type` - (Required) Type of the Log Stream - can either be `"aws_eventbridge"` or `"splunk_cloud_logstreaming"` only.

- `status` - (Optional) Log Stream Status - can either be ACTIVE or INACTIVE only. Default is ACTIVE.

- `settings` - (Required) Stream provider specific configuration.

  - `account_id` - (Required for `"aws_eventbridge"`) AWS account ID.

  - `event_source_name` - (Required for `"aws_eventbridge"`) An alphanumeric name (no spaces) to identify this event source in AWS EventBridge.`.

  - `region` - (Required for `"aws_eventbridge"`) The destination AWS region where event source is located.

  - `edition` - (Required for `"splunk_cloud_logstreaming"`) Edition of the Splunk Cloud instance. Could be one of: 'aws', 'aws_govcloud', 'gcp'.

  - `host` - (Required for `"splunk_cloud_logstreaming"`) The domain name for Splunk Cloud instance. Don't include http or https in the string. For example: 'acme.splunkcloud.com'.

  - `token` - (Required for `"splunk_cloud_logstreaming"`) The HEC token for your Splunk Cloud HTTP Event Collector.

## Attributes Reference

- `id` - Log Stream ID.

## Import

Okta Log Stream can be imported via the Okta ID.

```
$ terraform import okta_log_stream.example &#60;strema id&#62;
```
