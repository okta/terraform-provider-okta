resource "okta_log_stream" "eventbridge" {
  name   = "testAcc_replace_with_uuid EventBridge Updated"
  type   = "aws_eventbridge"
  status = "INACTIVE"
  settings {
    account_id        = "123456789012"
    region            = "eu-west-3"
    event_source_name = "testAcc_replace_with_uuid"
  }
}

resource "okta_log_stream" "splunk" {
  name   = "testAcc_replace_with_uuid Splunk Updated"
  type   = "splunk_cloud_logstreaming"
  status = "ACTIVE"
  settings {
    host    = "acme.splunkcloud.com"
    edition = "aws"
    token   = "58A7C8D6-4E2F-4C3B-8F5B-D4E2F3A4B5C6"
  }
}