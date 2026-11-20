resource "okta_group" "requester_groupA" {
  name        = "requester_test-reorder-ids-1"
  description = "requester_test-reorder-ids-1"
}

resource "okta_group" "requester_groupB" {
  name        = "requester_test-reorder-ids-2"
  description = "requester_test-reorder-ids-2"
}

resource "okta_app_group_assignment" "testA" {
  app_id   = "0oa15zse2by0vGzyP1d8"
  group_id = okta_group.requester_groupA.id
}

resource "okta_app_group_assignment" "testB" {
  app_id   = "0oa15zse2by0vGzyP1d8"
  group_id = okta_group.requester_groupB.id
}

resource "okta_request_condition" "test" {
  resource_id          = "0oa15zse2by0vGzyP1d8"
  approval_sequence_id = "69cbcb93a55844b4f3090b11"
  name                 = "test-reorder-ids"
  access_scope_settings {
    type = "GROUPS"

    ids {
      id = okta_group.requester_groupA.id
    }

    ids {
      id = okta_group.requester_groupB.id
    }

  }
  requester_settings {
    type = "GROUPS"

    ids {
      id = okta_group.requester_groupA.id
    }

    ids {
      id = okta_group.requester_groupB.id
    }
  }
}

