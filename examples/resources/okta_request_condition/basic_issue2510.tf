resource "okta_group" "requester" {
  name        = "requester_test"
  description = "requester_test"
}

resource "okta_request_condition" "test" {
  resource_id          = "0oasp3g29b1hqkcYE1d7"
  approval_sequence_id = "69251ae704a4d0a7fcdb870f"
  name                 = "issue-2510"

  access_scope_settings {
    type = "GROUPS"

    ids {
      id = "00gouu5aq9Gq0JGLH1d7"
    }
  }

  requester_settings {
    type = "GROUPS"

    ids {
      id = okta_group.requester.id
    }
  }
}