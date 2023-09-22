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

resource "okta_idp_saml_key" "test" {
  x5c = [okta_app_saml.test.certificate]
}

resource "okta_idp_saml" "test" {
  name                     = "testAcc_replace_with_uuid"
  acs_type                 = "INSTANCE"
  sso_url                  = "https://idp.example.com/test"
  sso_destination          = "https://idp.example.com/test"
  sso_binding              = "HTTP-POST"
  username_template        = "idpuser.email"
  kid                      = okta_idp_saml_key.test.id
  issuer                   = "https://idp.example.com/issuer"
  response_signature_scope = "RESPONSE"
  request_signature_scope  = "REQUEST"
  max_clock_skew           = 60
}
