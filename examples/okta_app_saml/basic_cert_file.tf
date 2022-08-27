resource "okta_app_saml" "test" {
  label                     = "testAcc_replace_with_uuid"
  sso_url                   = "http://google.com"
  recipient                 = "http://here.com"
  destination               = "http://its-about-the-journey.com"
  audience                  = "http://audience.com"
  subject_name_id_template  = "$${user.userName}"
  subject_name_id_format    = "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"
  response_signed           = true
  signature_algorithm       = "RSA_SHA256"
  digest_algorithm          = "SHA256"
  honor_force_authn         = false
  authn_context_class_ref   = "urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport"
  single_logout_issuer      = "https://dunshire.okta.com"
  single_logout_url         = "https://dunshire.okta.com/logout"
  single_logout_certificate = file("./slo_cert.crt")

  attribute_statements {
    type         = "GROUP"
    name         = "groups"
    filter_type  = "REGEX"
    filter_value = ".*"
  }
}
