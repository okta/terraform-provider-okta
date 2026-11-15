resource "okta_request_condition" "test" {
  resource_id          = "0oa15zse2by0vGzyP1d8"
  approval_sequence_id = "69cbcb93a55844b4f3090b11"
  name                 = "test-condition"
  access_scope_settings {
    type = "RESOURCE_DEFAULT"
  }
  requester_settings {
    type = "GROUPS"
    ids {
      id = "00g117xeeqyPXz13r1d8"
    }
  }
}
