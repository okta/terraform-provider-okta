resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "Jones"
  login      = "john_replace_with_uuid@ledzeppelin.com"
  email      = "john_replace_with_uuid@ledzeppelin.com"
}

resource "okta_user" "test1" {
  first_name = "TestAcc"
  last_name  = "Entwhistle"
  login      = "john_replace_with_uuid@thewho.com"
  email      = "john_replace_with_uuid@thewho.com"
}

resource "okta_user" "test2" {
  first_name = "TestAcc"
  last_name  = "Doe"
  login      = "john_replace_with_uuid@unknown.com"
  email      = "john_replace_with_uuid@unknown.com"
}

resource "okta_user" "test3" {
  first_name = "TestAcc"
  last_name  = "Astley"
  login      = "rick_astley_replace_with_uuid@rickrollin.com"
  email      = "rick_astley_replace_with_uuid@rickrollin.com"
}
