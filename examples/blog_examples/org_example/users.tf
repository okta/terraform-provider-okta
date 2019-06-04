resource okta_user_schema dept {
  index  = "dept"
  title  = "Department"
  type   = "string"
  master = "OKTA"
  enum   = ["Over Lords", "Plebians"]
}

resource okta_user john_johnson {
  first_name = "John"
  last_name  = "Johnson"
  login      = "jj@example.com"
  email      = "jj@example.com"

  custom_profile_attributes = {
    dept = "Plebians"
  }
}

resource okta_user bill_billson {
  first_name  = "Bill"
  last_name   = "Billson"
  login       = "bb@example.com"
  email       = "bb@example.com"
  admin_roles = ["SUPER_ADMIN"]

  custom_profile_attributes = {
    dept = "Over Lords"
  }
}
