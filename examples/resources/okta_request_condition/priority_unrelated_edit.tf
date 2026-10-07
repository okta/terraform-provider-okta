resource "okta_request_condition" "test_priority_unrelated" {
  resource_id          = "0oa14jmvwvzsRjALs1d8"
  approval_sequence_id = "68d224058c0cff364ca377e8"
  name                 = "test-condition-unrelated-before"
  description          = "before"
  priority             = 5
  access_scope_settings {
    type = "RESOURCE_DEFAULT"
  }
  requester_settings {
    type = "EVERYONE"
  }
}
