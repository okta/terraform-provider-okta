resource "okta_request_conditions" "test" {
  resource_id = "0oa10z0gowhJuJBAo1d8"

  condition {
    name                 = "test-conditions-c"
    approval_sequence_id = "69cbcb93a55844b4f3090b11"
    priority             = 0
    status               = "ACTIVE"
    access_scope_settings {
      type = "RESOURCE_DEFAULT"
    }
    requester_settings {
      type = "EVERYONE"
    }
  }

  condition {
    name                 = "test-conditions-a"
    approval_sequence_id = "69cbcb93a55844b4f3090b11"
    priority             = 1
    status               = "INACTIVE"
    access_scope_settings {
      type = "RESOURCE_DEFAULT"
    }
    requester_settings {
      type = "EVERYONE"
    }
  }

  condition {
    name                 = "test-conditions-d"
    approval_sequence_id = "69cbcb93a55844b4f3090b11"
    priority             = 2
    status               = "ACTIVE"
    access_scope_settings {
      type = "RESOURCE_DEFAULT"
    }
    requester_settings {
      type = "EVERYONE"
    }
  }
}
