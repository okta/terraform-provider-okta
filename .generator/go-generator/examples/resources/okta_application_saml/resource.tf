resource "okta_application_saml" "example" {
  sign_on_mode = "<sign_on_mode>"
  label = "Example Label"
  index = 0
  url = "https://example.com"
  allow_multiple_acs_endpoints = "https://example.com"
  assertion_signed = true
  audience = "<audience>"
  authn_context_class_ref = "<authn_context_class_ref>"
  destination = "<destination>"
  digest_algorithm = "<digest_algorithm>"
  honor_force_authn = true
  idp_issuer = "<idp_issuer>"
  recipient = "<recipient>"
  request_compressed = true
  response_signed = true
  signature_algorithm = "<signature_algorithm>"
  sso_acs_url = "https://example.com"
  subject_name_id_format = "Example Subject Name Id Format"
  subject_name_id_template = "Example Subject Name Id Template"

  # Optional fields
  # error_redirect_url = "https://example.com"
  # login_redirect_url = "https://example.com"
  # self_service = true
  # kid = "<kid>"
  # rotation_mode = "<rotation_mode>"
}
