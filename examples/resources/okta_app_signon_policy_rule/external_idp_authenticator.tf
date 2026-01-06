resource "okta_app_saml" "example" {
  label                    = "My App"
  sso_url                  = "https://google.com"
  recipient                = "https://here.com"
  destination              = "https://its-about-the-journey.com"
  audience                 = "https://audience.com"
  status                   = "ACTIVE"
  subject_name_id_template = "$${user.userName}"
  subject_name_id_format   = "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"
  signature_algorithm      = "RSA_SHA256"
  response_signed          = true
  digest_algorithm         = "SHA256"
  honor_force_authn        = false
  authn_context_class_ref  = "urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport"
}


resource "okta_app_signon_policy" "example" {
  name        = "testAcc_Test_App_replace_with_uuid"
  description = "The app signon policy used by our test app"
}

resource "okta_app_signon_policy_rule" "test" {
  policy_id   = okta_app_signon_policy.example.id
  name        = "test-rule"
  factor_mode = "2FA"
  type        = "ASSURANCE"
  constraints = [
    jsonencode(
      {
        "knowledge" : {
          "required" : true,
          "types" : [
            "password"
          ],
          "reauthenticateIn" : "PT2H"
        },
        "possession" : {
          "required" : true,
          "authenticationMethods" : [
            {
              "key" : "external_idp",
              "id" : "auttixlvpvfch3xmO1d7",
              "method" : "idp"
            }
          ],
          "userPresence" : "OPTIONAL"
        }
      }
    )
  ]
}