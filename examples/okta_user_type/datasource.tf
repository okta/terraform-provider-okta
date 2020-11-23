resource okta_user_type test {
  name         = "testAcc_replace_with_uuid"
  display_name = "testing"
  description  = "testing"
}

data okta_user_type test {
  name = okta_user_type.test.name
}
