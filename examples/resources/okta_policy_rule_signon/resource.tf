resource "okta_policy_signon" "test" {
  name        = "Example Policy"
  status      = "ACTIVE"
  description = "Example Policy"
}

data "okta_behavior" "new_city" {
  name = "New City"
}

resource "okta_policy_rule_signon" "example" {
  access             = "CHALLENGE"
  authtype           = "RADIUS"
  name               = "Example Policy Rule"
  network_connection = "ANYWHERE"
  policy_id          = okta_policy_signon.example.id
  status             = "ACTIVE"
  risc_level         = "HIGH"
  behaviors          = [data.okta_behavior.new_city.id]
  factor_sequence {
    primary_criteria_factor_type = "token:hotp" // TOTP
    primary_criteria_provider    = "CUSTOM"
    secondary_criteria {
      factor_type = "token:software:totp" // Okta Verify
      provider    = "OKTA"
    }
    secondary_criteria { // Okta Verify Push
      factor_type = "push"
      provider    = "OKTA"
    }
    secondary_criteria { // Password
      factor_type = "password"
      provider    = "OKTA"
    }
    secondary_criteria { // Security Question
      factor_type = "question"
      provider    = "OKTA"
    }
    secondary_criteria { // SMS
      factor_type = "sms"
      provider    = "OKTA"
    }
    secondary_criteria { // Google Auth
      factor_type = "token:software:totp"
      provider    = "GOOGLE"
    }
    secondary_criteria { // Email
      factor_type = "email"
      provider    = "OKTA"
    }
    secondary_criteria { // Voice Call
      factor_type = "call"
      provider    = "OKTA"
    }
    secondary_criteria { // FIDO2 (WebAuthn)
      factor_type = "webauthn"
      provider    = "FIDO"
    }
    secondary_criteria { // RSA
      factor_type = "token"
      provider    = "RSA"
    }
    secondary_criteria { // Symantec VIP
      factor_type = "token"
      provider    = "SYMANTEC"
    }
  }
  factor_sequence {
    primary_criteria_factor_type = "token:software:totp" // Okta Verify
    primary_criteria_provider    = "OKTA"
  }
}
