resource "okta_user" "included-replace_with_uuid" {
  first_name = "TestAcc"
  last_name  = "Included"
  login      = "testAcc-included-replace_with_uuid@example.com"
  email      = "testAcc-included-replace_with_uuid@example.com"
}

resource "okta_user" "excluded-replace_with_uuid" {
  first_name = "TestAcc"
  last_name  = "Excluded"
  login      = "testAcc-excluded-replace_with_uuid@example.com"
  email      = "testAcc-excluded-replace_with_uuid@example.com"
}

resource "okta_group" "included-replace_with_uuid" {
  name = "testAcc_included_replace_with_uuid"
}

resource "okta_group" "excluded-replace_with_uuid" {
  name = "testAcc_excluded_replace_with_uuid"
}

resource "okta_network_zone" "test" {
  name     = "testAcc_replace_with_uuid"
  type     = "IP"
  gateways = ["1.2.3.4/24", "2.3.4.5-2.3.4.15"]
  proxies  = ["2.2.3.4/24", "3.3.4.5-3.3.4.15"]
  status   = "ACTIVE"
}

data "okta_default_policy" "default-replace_with_uuid" {
  type = "PASSWORD"
}

resource "okta_policy_rule_password" "testAcc_replace_with_uuid" {
  policy_id = data.okta_default_policy.default-replace_with_uuid.id
  name      = "testAcc_replace_with_uuid"
  status    = "ACTIVE"

  password_change = "ALLOW"
  password_reset  = "ALLOW"
  password_unlock = "ALLOW"

  users_included     = [okta_user.included-replace_with_uuid.id]
  users_excluded     = [okta_user.excluded-replace_with_uuid.id]
  groups_included    = [okta_group.included-replace_with_uuid.id]
  groups_excluded    = [okta_group.excluded-replace_with_uuid.id]
  network_connection = "ZONE"
  network_includes   = [okta_network_zone.test.id]

  password_reset_access_control = "LEGACY"

  password_reset_requirement {
    method_constraints {
      method                 = "otp"
      allowed_authenticators = ["google_otp"]
    }
    primary_methods = ["otp", "email"]
    step_up_enabled = true
    step_up_methods = ["security_question"]
  }
}
