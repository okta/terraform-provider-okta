resource "okta_rate_limiting" "example" {
  login                  = "ENFORCE"
  authorize              = "ENFORCE"
  communications_enabled = true
}
