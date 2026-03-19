resource "okta_log_stream" "test_aws" {
  name   = "testAcc_replace_with_uuid AWS"
  type   = "aws_eventbridge"
  status = "ACTIVE"
  settings {
    account_id        = "123456789012"
    region            = "eu-west-3"
    event_source_name = "testAcc_replace_with_uuid_AWS"
  }
}

resource "okta_log_stream" "test_splunk" {
  name   = "testAcc_replace_with_uuid Splunk"
  type   = "splunk_cloud_logstreaming"
  status = "ACTIVE"
  settings {
    host    = "acme.splunkcloud.com"
    edition = "aws"
    token   = "58A7C8D6-4E2F-4C3B-8F5B-D4E2F3A4B5C6"
  }
}

data "okta_log_stream" "test_by_name" {
  name       = okta_log_stream.test_splunk.name
  depends_on = [okta_log_stream.test_splunk, okta_log_stream.test_aws]
}

data "okta_log_stream" "test_by_id" {
  id         = okta_log_stream.test_aws.id
  depends_on = [okta_log_stream.test_splunk, okta_log_stream.test_aws]
}