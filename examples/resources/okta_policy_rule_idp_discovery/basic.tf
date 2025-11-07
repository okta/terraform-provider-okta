data "okta_policy" "test" {
  name = "Idp Discovery Policy"
  type = "IDP_DISCOVERY"
}

resource "okta_policy_rule_idp_discovery" "test" {
  policy_id = data.okta_policy.test.id
  priority  = 1
  name      = "testAcc_replace_with_uuid"

  idp_providers {
    type = "SAML2"
    id   = okta_idp_saml.test.id
  }

  user_identifier_type = "ATTRIBUTE"

  // Don't have a company schema in this account, just chosing something always there
  user_identifier_attribute = "firstName"

  user_identifier_patterns {
    match_type = "EQUALS"
    value      = "Articulate"
  }
}

resource "okta_idp_saml" "test" {
  name                     = "testAcc_replace_with_uuid"
  acs_type                 = "INSTANCE"
  sso_url                  = "https://idp.example.com"
  sso_destination          = "https://idp.example.com"
  sso_binding              = "HTTP-POST"
  username_template        = "idpuser.email"
  issuer                   = "https://idp.example.com"
  request_signature_scope  = "REQUEST"
  response_signature_scope = "ANY"
  kid                      = okta_idp_saml_key.test.id
}

resource "okta_idp_saml_key" "test" {
  x5c = [okta_app_saml.test.certificate]
}

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
    name   = "firstName"
    values = ["user.firstName"]
  }

  attribute_statements {
    name   = "lastName"
    values = ["user.lastName"]
  }

  attribute_statements {
    name   = "email"
    values = ["user.email"]
  }

  attribute_statements {
    name   = "company"
    values = ["Articulate"]
  }
}