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
  attribute_statements {
    type         = "GROUP"
    name         = "groups"
    filter_type  = "REGEX"
    filter_value = ".*"
  }
}

data "okta_app_signon_policy" "test" {
  app_id = okta_app_saml.test.id
}

resource "okta_app_signon_policy_rule" "test" {
  policy_id = data.okta_app_signon_policy.test.id
  name      = "testAcc_replace_with_uuid"
  status    = "ACTIVE"
  priority  = 1
  platform_include {
    os_expression = ""
    os_type       = "OTHER"
    type          = "DESKTOP"
  }
}

data "okta_policies_rule_access_policy" "test" {
  id        = okta_app_signon_policy_rule.test.id
  policy_id = data.okta_app_signon_policy.test.id
}

