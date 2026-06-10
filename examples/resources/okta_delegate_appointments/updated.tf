resource "okta_user" "delegator" {
  first_name = "TestAcc"
  last_name  = "Delegator"
  login      = "testAcc-delegator-replace_with_uuid@example.com"
  email      = "testAcc-delegator-replace_with_uuid@example.com"
}

resource "okta_user" "delegate2" {
  first_name = "TestAcc"
  last_name  = "Delegate2"
  login      = "testAcc-delegate2-replace_with_uuid@example.com"
  email      = "testAcc-delegate2-replace_with_uuid@example.com"
}

resource "okta_delegate_appointments" "test" {
  principal_id = okta_user.delegator.id

  appointments {
    delegate_id = okta_user.delegate2.id
    note        = "Switched delegate"
  }
}
