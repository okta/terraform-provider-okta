resource "okta_request_condition" "test_priority_0" {
  status               = "ACTIVE"
  resource_id          = "0oasp3g29b1hqkcYE1d7"
  approval_sequence_id = "69251ae704a4d0a7fcdb870f"
  name                 = "test-condition-priority-0"
  priority             = 0
  access_scope_settings {
    type = "RESOURCE_DEFAULT"
  }
  requester_settings {
    type = "EVERYONE"
  }
}

resource "okta_request_condition" "test_priority_1" {
  status               = "ACTIVE"
  resource_id          = "0oasp3g29b1hqkcYE1d7"
  approval_sequence_id = "69251ae704a4d0a7fcdb870f"
  name                 = "test-condition-priority-1"
  priority             = 1
  access_scope_settings {
    type = "RESOURCE_DEFAULT"
  }
  requester_settings {
    type = "EVERYONE"
  }
}
