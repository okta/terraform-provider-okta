resource "okta_request_condition" "test" {
  resource_id          = "0oaoum6j3cElINe1z1d7"
  approval_sequence_id = "68920b41386747a673869356"
  name                 = "test-condition"
  access_scope_settings {
    type = "RESOURCE_DEFAULT"
  }
  requester_settings {
    type = "EVERYONE"
  }
}

resource "okta_request_condition" "test_active" {
  status               = "ACTIVE"
  resource_id          = "0oaoum6j3cElINe1z1d7"
  approval_sequence_id = "68920b41386747a673869356"
  name                 = "test-condition-active"
  access_scope_settings {
    type = "RESOURCE_DEFAULT"
  }
  requester_settings {
    type = "EVERYONE"
  }
}
