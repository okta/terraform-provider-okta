resource "okta_trusted_origin" "testAcc_%[1]d" {
  name   = "test-acc-%[1]d"
  origin = "https://example-%[1]d.com"
  scopes = ["CORS"]
}
