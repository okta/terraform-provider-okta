resource "okta_request_condition" "test" {
  resource_id          = "0oasp3g29b1hqkcYE1d7"
  approval_sequence_id = "69251ae704a4d0a7fcdb870f"
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
  resource_id = "0oasp3g29b1hqkcYE1d7"
}
