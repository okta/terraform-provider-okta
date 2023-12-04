resource "okta_log_stream" "eventbridge_log_stream_example" {
  name     = "testAcc_replace_with_uuid EventBridge"
  type     = "aws_eventbridge"
  status   = "ACTIVE"
  settings {
    account_id = "123456789012"
    region     = "eu-west-3"
    event_source_name = "testAcc_replace_with_uuid"
  }
}

resource "okta_log_stream" "splunk_log_stream_example" {
  name              = "testAcc_replace_with_uuid Splunk"
  type              = "splunk_cloud_logstreaming"
  status            = "ACTIVE"
  settings {
    host = "acme.splunkcloud.com"
    edition = "aws"
    token = "58A7C8D6-4E2F-4C3B-8F5B-D4E2F3A4B5C6"
  }
}
