resource "okta_app_saml" "test" {
  label                    = "testAcc_replace_with_uuid"
  sso_url                  = "http://google.com"
  recipient                = "http://here.com"
  destination              = "http://its-about-the-journey.com"
  audience                 = "http://audience.com"
  subject_name_id_template = "$${source.login}"
  subject_name_id_format   = "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified"
  response_signed          = true
  assertion_signed         = true
  signature_algorithm      = "RSA_SHA1"
  digest_algorithm         = "SHA1"
  honor_force_authn        = true
  authn_context_class_ref  = "urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport"

  attribute_statements {
    name      = "Attr One"
    namespace = "urn:oasis:names:tc:SAML:2.0:attrname-format:unspecified"
    values    = ["val"]
  }

  attribute_statements {
    type         = "GROUP"
    name         = "Attr Two"
    filter_type  = "STARTS_WITH"
    filter_value = "test"
  }

  acs_endpoints = ["https://example.com", "https://okta.com"]
}
