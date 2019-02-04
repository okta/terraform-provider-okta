resource "okta_trusted_origin" "testAcc_%[1]d" {
  name   = "test-acc-%[1]d"
  active = false
  origin = "https://example2-%[1]d.com"
  scopes = ["CORS", "REDIRECT"]
}
