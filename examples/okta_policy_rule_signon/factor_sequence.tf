data "okta_group" "all" {
  name = "Everyone"
}

resource "okta_policy_signon" "test" {
  name        = "testAcc_replace_with_uuid"
  status      = "ACTIVE"
  description = "Terraform Acceptance Test SignOn Policy"
  groups_included = [
    data.okta_group.all.id
  ]
}

data "okta_behavior" "new_city" {
  name = "New City"
}

resource "okta_network_zone" "test" {
  name = "testAcc_replace_with_uuid"
  type = "IP"
  gateways = [
    "1.2.3.4/24",
    "2.3.4.5-2.3.4.15"
  ]
  proxies = [
    "2.2.3.4/24",
    "3.3.4.5-3.3.4.15"
  ]
  depends_on = [okta_policy_rule_signon.test]
  status   = "ACTIVE"
}

resource "okta_policy_rule_signon" "test" {
  policy_id = okta_policy_signon.test.id
  name      = "testAcc_replace_with_uuid"
  status    = "ACTIVE"
  access    = "CHALLENGE"
  behaviors = [
    data.okta_behavior.new_city.id
  ]
  factor_sequence {
    primary_criteria_factor_type = "password"
    primary_criteria_provider    = "OKTA"
    secondary_criteria {
      factor_type = "push"
      provider    = "OKTA"
    }
    secondary_criteria {
      factor_type = "token:hotp"
      provider    = "CUSTOM"
    }
    secondary_criteria {
      factor_type = "token:software:totp"
      provider    = "OKTA"
    }
  }
  factor_sequence {
    primary_criteria_factor_type = "token:hotp"
    primary_criteria_provider    = "CUSTOM"
  }

  depends_on = [
    okta_factor.okta_sms,
    okta_factor.okta_email,
    okta_factor.hotp
  ]
}

resource "okta_factor" "okta_sms" {
  provider_id = "okta_sms"
}

resource "okta_factor" "okta_email" {
  provider_id = "okta_email"
}

resource "okta_factor" "hotp" {
  provider_id = "hotp"
}

