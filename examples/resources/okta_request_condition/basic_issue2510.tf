resource "okta_group" "requester" {
  name        = "requester_test"
  description = "requester_test"
}

resource "okta_group" "requester_groupA" {
  name        = "requester_test_groupA"
  description = "requester_test_groupA"
}

resource "okta_app_group_assignment" "testA" {
  app_id   = "0oa15zse2by0vGzyP1d8"
  group_id = okta_group.requester_groupA.id
}


resource "okta_request_condition" "test" {
  resource_id          = "0oa15zse2by0vGzyP1d8"
  approval_sequence_id = "69cbcb93a55844b4f3090b11"
  name                 = "issue-2510"

  access_scope_settings {
    type = "GROUPS"

    ids {
      id = okta_group.requester_groupA.id
    }
  }

  requester_settings {
    type = "GROUPS"

    ids {
      id = okta_group.requester.id
    }
  }
  depends_on = ["okta_app_group_assignment.testA"]
}