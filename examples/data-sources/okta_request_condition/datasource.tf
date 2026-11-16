resource "okta_request_condition" "test" {
  resource_id          = "0oa10z0gowhJuJBAo1d8"
  approval_sequence_id = "69cbcb93a55844b4f3090b11"
  name                 = "request-condition-test"
  access_scope_settings {
    type = "RESOURCE_DEFAULT"
  }
  requester_settings {
    type = "EVERYONE"
  }
}
data "okta_request_condition" "test" {
  id          = okta_request_condition.test.id
  resource_id = "0oa10z0gowhJuJBAo1d8"
}
