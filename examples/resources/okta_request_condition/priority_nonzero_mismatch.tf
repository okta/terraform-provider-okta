resource "okta_request_condition" "condition-one" {
  resource_id          = "0oa14jmvwvzsRjALs1d8"
  approval_sequence_id = "68d224058c0cff364ca377e8"
  name                 = "okta-1292840-request-condition-one"
  priority             = 0

  access_scope_settings {
    type = "RESOURCE_DEFAULT"
  }
  requester_settings {
    type = "EVERYONE"
  }
}

resource "okta_request_condition" "condition-two" {
  resource_id          = "0oa14jmvwvzsRjALs1d8"
  approval_sequence_id = "68d224058c0cff364ca377e8"
  name                 = "okta-1292840-request-condition-two"
  priority             = 1

  access_scope_settings {
    type = "RESOURCE_DEFAULT"
  }
  requester_settings {
    type = "EVERYONE"
  }

  depends_on = [okta_request_condition.condition-one]
}

resource "okta_request_condition" "condition-three" {
  resource_id          = "0oa14jmvwvzsRjALs1d8"
  approval_sequence_id = "68d224058c0cff364ca377e8"
  name                 = "okta-1292840-request-condition-three"
  priority             = 50

  access_scope_settings {
    type = "RESOURCE_DEFAULT"
  }
  requester_settings {
    type = "EVERYONE"
  }

  depends_on = [okta_request_condition.condition-two]
}

resource "okta_request_condition" "condition-four" {
  resource_id          = "0oa14jmvwvzsRjALs1d8"
  approval_sequence_id = "68d224058c0cff364ca377e8"
  name                 = "okta-1292840-request-condition-four"
  priority             = 100

  access_scope_settings {
    type = "RESOURCE_DEFAULT"
  }
  requester_settings {
    type = "EVERYONE"
  }

  depends_on = [okta_request_condition.condition-three]
}
