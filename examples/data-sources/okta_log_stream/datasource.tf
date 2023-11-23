resource "okta_log_stream" "test" {
  name     = "testAcc_replace_with_uuid"
  type     = "aws_eventbridge"
  status   = "ACTIVE"
  settings {
    account_id = "123456789012"
    region     = "eu-west-3"
    event_source_name = "testAcc_replace_with_uuid"
  }
}

data "okta_log_stream" "test" {
  name     = okta_log_stream.test.name
}
