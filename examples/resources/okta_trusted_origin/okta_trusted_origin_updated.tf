resource "okta_trusted_origin" "testAcc_replace_with_uuid" {
  name   = "testAcc-replace_with_uuid"
  active = false
  origin = "https://example2-replace_with_uuid.com"
  scopes = ["CORS", "REDIRECT"]
}
