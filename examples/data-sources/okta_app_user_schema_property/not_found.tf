resource "okta_app_saml" "test" {
  label                    = "testAcc_replace_with_uuid"
  sso_url                  = "http://example.com/sso"
  recipient                = "http://example.com"
  destination              = "http://example.com"
  audience                 = "http://example.com/audience"
  subject_name_id_template = "$${user.userName}"
  subject_name_id_format   = "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"
  response_signed          = true
  signature_algorithm      = "RSA_SHA256"
  digest_algorithm         = "SHA256"
  honor_force_authn        = false
  authn_context_class_ref  = "urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport"
}

# Try to query a non-existent property - should fail
data "okta_app_user_schema_property" "test" {
  app_id = okta_app_saml.test.id
  index  = "nonExistentProperty"
}
