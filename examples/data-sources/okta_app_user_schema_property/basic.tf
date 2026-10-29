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

  attribute_statements {
    type      = "EXPRESSION"
    name      = "email"
    namespace = "urn:oasis:names:tc:SAML:2.0:attrname-format:unspecified"
    values    = ["user.email"]
  }
}

# Create a custom schema property that we can query with the data source
resource "okta_app_user_schema_property" "test" {
  app_id      = okta_app_saml.test.id
  index       = "testCustomProperty"
  title       = "Test Custom Property"
  type        = "string"
  description = "Test description"
  permissions = "READ_ONLY"
  master      = "PROFILE_MASTER"
}

# Query the property using the data source
data "okta_app_user_schema_property" "test" {
  app_id = okta_app_saml.test.id
  index  = okta_app_user_schema_property.test.index

  depends_on = [okta_app_user_schema_property.test]
}
