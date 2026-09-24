resource "okta_app_saml" "test" {
  label                    = "testAcc_replace_with_uuid"
  sso_url                  = "https://example.com/sso"
  recipient                = "https://example.com/recipient"
  destination              = "https://example.com/destination"
  audience                 = "https://example.com/audience"
  subject_name_id_template = "$${user.userName}"
  subject_name_id_format   = "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"
  response_signed          = true
  signature_algorithm      = "RSA_SHA256"
  digest_algorithm         = "SHA256"
  honor_force_authn        = false
  authn_context_class_ref  = "urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport"
}

data "okta_app_signon_policy" "test" {
  app_id = okta_app_saml.test.id
}

# Multiple rules exercising the keep_me_signed_in ("Option to stay signed in")
# block in different combinations:
#   - rule[0]: NOT_ALLOWED with no prompt frequency
#   - rule[1]: ALLOWED with a 50h prompt frequency
#   - rule[2]: ALLOWED with a 168h (7 day) prompt frequency
#   - rule[3]: NOT_ALLOWED with no prompt frequency (regression case for the
#              "null -> empty string" inconsistent-result-after-apply bug)
resource "okta_app_signon_policy_rules" "test" {
  policy_id = data.okta_app_signon_policy.test.id

  rule {
    name               = "Rule1-replace_with_uuid"
    priority           = 1
    status             = "ACTIVE"
    factor_mode        = "2FA"
    inactivity_period  = "PT1H"
    network_connection = "ANYWHERE"

    keep_me_signed_in {
      post_auth = "NOT_ALLOWED"
    }
  }

  rule {
    name               = "Rule2-replace_with_uuid"
    priority           = 2
    status             = "ACTIVE"
    factor_mode        = "2FA"
    inactivity_period  = "PT1H"
    network_connection = "ANYWHERE"

    keep_me_signed_in {
      post_auth                  = "ALLOWED"
      post_auth_prompt_frequency = "PT50H"
    }
  }

  rule {
    name               = "Rule3-replace_with_uuid"
    priority           = 3
    status             = "ACTIVE"
    factor_mode        = "2FA"
    inactivity_period  = "PT1H"
    network_connection = "ANYWHERE"

    keep_me_signed_in {
      post_auth                  = "ALLOWED"
      post_auth_prompt_frequency = "PT168H"
    }
  }

  rule {
    name               = "Rule4-replace_with_uuid"
    priority           = 4
    status             = "ACTIVE"
    factor_mode        = "2FA"
    inactivity_period  = "PT1H"
    network_connection = "ANYWHERE"

    keep_me_signed_in {
      post_auth = "NOT_ALLOWED"
    }
  }
}
