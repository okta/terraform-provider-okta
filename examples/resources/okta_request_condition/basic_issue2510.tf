resource "okta_group" "requester" {
  name        = "requester_test"
  description = "requester_test"
}

resource "okta_request_condition" "test" {
  resource_id          = "0oaqgxmg2n2FjHLzw1d7"
  approval_sequence_id = "68d224058c0cff364ca377e8"
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