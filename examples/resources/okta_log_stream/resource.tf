### AWS EventBridge
resource "okta_log_stream" "example" {
  name   = "EventBridge Log Stream"
  type   = "aws_eventbridge"
  status = "ACTIVE"
  settings {
    account_id        = "123456789012"
    region            = "us-north-1"
    event_source_name = "okta_log_stream"
  }
}

### Splunk Event Collector
resource "okta_log_stream" "example" {
  name   = "Splunk log Stream"
  type   = "splunk_cloud_logstreaming"
  status = "ACTIVE"
  settings {
    host    = "acme.splunkcloud.com"
    edition = "gcp"
    token   = "YOUR_HEC_TOKEN"
  }
}
