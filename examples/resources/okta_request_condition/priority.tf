resource "okta_request_condition" "priority_first" {
  resource_id          = "0oasp3g29b1hqkcYE1d7"
  approval_sequence_id = "69251ae704a4d0a7fcdb870f"
  name                 = "priority-first"
  priority             = 1
  access_scope_settings {
    type = "RESOURCE_DEFAULT"
  }
  requester_settings {
    type = "EVERYONE"
  }
}

resource "okta_request_condition" "priority_second" {
  resource_id          = "0oasp3g29b1hqkcYE1d7"
  approval_sequence_id = "69251ae704a4d0a7fcdb870f"
  name                 = "priority-second"
  priority             = 2
  access_scope_settings {
    type = "RESOURCE_DEFAULT"
  }
  requester_settings {
    type = "EVERYONE"
  }

  depends_on = [okta_request_condition.priority_first]
}
