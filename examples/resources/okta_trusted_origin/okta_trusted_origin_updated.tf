resource "okta_trusted_origin" "testAcc_replace_with_uuid" {
  name   = "testAcc-replace_with_uuid"
  status = "INACTIVE"
  origin = "https://example2-replace_with_uuid.com"
  scopes {
    type = "CORS"
  }
  scopes {
    type = "REDIRECT"
  }
}
