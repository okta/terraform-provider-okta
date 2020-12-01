resource okta_user_type test {
  name         = "testAcc_replace_with_uuid"
  display_name = "Terraform Acceptance Test User Type"
  description  = "Terraform Acceptance Test User Type"
}

data okta_user_type test {
  name = okta_user_type.test.name
  display_name = okta_user_type.test.display_name
}

resource okta_user_type test_1 {
  name         = format("testAcc_%s_1", data.okta_user_type.test.name)
  display_name = format("source_%s_1", data.okta_user_type.test.display_name)
  description  = "Terraform Acceptance Test User Type"
}
