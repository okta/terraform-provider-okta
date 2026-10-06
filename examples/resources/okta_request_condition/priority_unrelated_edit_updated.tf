resource "okta_request_condition" "test_priority_unrelated" {
  resource_id          = "0oasp3g29b1hqkcYE1d7"
  approval_sequence_id = "69251ae704a4d0a7fcdb870f"
  name                 = "test-condition-unrelated-after"
  description          = "after"
  priority             = 5
  access_scope_settings {
    type = "RESOURCE_DEFAULT"
  }
  requester_settings {
    type = "EVERYONE"
  }
}
