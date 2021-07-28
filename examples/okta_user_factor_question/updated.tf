resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name = "Jones"
  login = "john_replace_with_uuid@ledzeppelin.com"
  email = "john_replace_with_uuid@ledzeppelin.com"
}

resource "okta_factor" "test_factor" {
  provider_id = "okta_question"
  active = true
}

resource "okta_user_factor_question" "test" {
  user_id = okta_user.test.id
  key = "name_of_first_plush_toy"
  answer = "meatball"
  depends_on = [okta_factor.test_factor]
}
