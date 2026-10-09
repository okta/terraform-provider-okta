resource "okta_request_condition" "test" {
  resource_id = "0oa11ez9vooGRnbZv1d8"
  approval_sequence_id = "68d224058c0cff364ca377e8"
  name                 = "test-reorder-ids"
  access_scope_settings {
    type = "GROUPS"

    ids {
      id = "00gwkax39dIJgiptR1d7"
    }

    ids {
      id = "00gwkaw31mtq5LnuZ1d7"
    }
  }
  requester_settings {
    type = "GROUPS"

    ids {
      id = "00gwkax39dIJgiptR1d7"
    }

    ids {
      id = "00gwkaw31mtq5LnuZ1d7"
    }
  }
}

