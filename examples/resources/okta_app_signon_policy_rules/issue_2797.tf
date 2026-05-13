resource "okta_app_saml" "test" {
  label                    = "testAcc_replace_with_uuid"
  sso_url                  = "http://google.com"
  recipient                = "http://here.com"
  destination              = "http://its-about-the-journey.com"
  audience                 = "http://audience.com"
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

resource "okta_app_signon_policy_rules" "test" {
  policy_id = data.okta_app_signon_policy.test.id

  rule {
    name        = "testAcc_replace_with_uuid"
    priority    = 1
    factor_mode = "2FA"
    status      = "ACTIVE"
  }

  rule {
    name        = "testAcc_replace_with_uuid_2"
    priority    = 2
    factor_mode = "1FA"
    access      = "ALLOW"
    status      = "INACTIVE"
  }
}