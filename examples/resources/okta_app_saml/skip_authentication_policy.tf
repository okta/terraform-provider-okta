# Example: Skip Authentication Policy Operations for SAML App
# This example demonstrates how to use the skip_authentication_policy flag
# to prevent the provider from managing authentication policies for a SAML application.

resource "okta_app_saml" "example_skip_policy" {
  label                    = "Example SAML App - Skip Policy"
  sso_url                  = "https://example.com/sso"
  recipient                = "https://example.com/sso"
  destination              = "https://example.com/sso"
  audience                 = "https://example.com"
  subject_name_id_template = "$${user.userName}"
  subject_name_id_format   = "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"
  response_signed          = true
  signature_algorithm      = "RSA_SHA256"
  digest_algorithm         = "SHA256"

  # Skip authentication policy operations
  # When set to true, the provider will not attempt to create, update, or delete
  # authentication policies for this application. This is useful when you want
  # to manage authentication policies manually or when the application should
  # use the default policy without explicit configuration.
  skip_authentication_policy = true
}

# Example: Regular SAML app with authentication policy management
resource "okta_app_saml" "example_with_policy" {
  label                    = "Example SAML App - With Policy"
  sso_url                  = "https://example.com/sso"
  recipient                = "https://example.com/sso"
  destination              = "https://example.com/sso"
  audience                 = "https://example.com"
  subject_name_id_template = "$${user.userName}"
  subject_name_id_format   = "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"
  response_signed          = true
  signature_algorithm      = "RSA_SHA256"
  digest_algorithm         = "SHA256"

  # This app will have authentication policy operations performed normally
  # The provider will assign it to the default policy if none is specified
  skip_authentication_policy = false
}
