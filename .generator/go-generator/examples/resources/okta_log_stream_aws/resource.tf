resource "okta_log_stream_aws" "example" {
  type = "<type>"
  name = "Example Name"
  account_id = "<account-id>"
  event_source_name = "Example Event Source Name"
  region = "<region>"
}
