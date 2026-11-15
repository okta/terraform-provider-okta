resource "okta_request_condition" "test_status" {
  status               = "ACTIVE"
  resource_id          = "0oa15zse2by0vGzyP1d8"
  approval_sequence_id = "69cbcb93a55844b4f3090b11"
  name                 = "test-condition-status"
  access_scope_settings {
    type = "RESOURCE_DEFAULT"
  }
  requester_settings {
    type = "EVERYONE"
  }
}

