data "okta_group" "test" {
  name = "Everyone"
}

resource "okta_group" "test" {
  name = "testAcc"
}

resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "Smith"
  login      = "testAcc_%[1]d@testing.com"
  email      = "testAcc_%[1]d@testing.com"
}

resource "okta_group_rule" "test" {
  name = "testAcc_%[1]d"
  status = "ACTIVE"
  group_assignments = ["${data.okta_group.id}"]
  user_blacklist = ["${okta_user.test.id}"]
  group_blacklist = ["${okta_group.test.id}"]
  expression_type = "urn:okta:expression:1.0"
  expression_value = "user.role==\\\"Engineer\\\""
}
