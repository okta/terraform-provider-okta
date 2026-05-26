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

# Rule with Keep Me Signed In - ALLOWED with prompt frequency
resource "okta_app_signon_policy_rule" "test" {
  policy_id = data.okta_app_signon_policy.test.id
  name      = "testAcc_replace_with_uuid"
  access    = "ALLOW"

  keep_me_signed_in {
    post_auth                  = "ALLOWED"
    post_auth_prompt_frequency = "PT168H"
  }
}
