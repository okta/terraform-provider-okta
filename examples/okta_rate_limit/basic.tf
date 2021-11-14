resource "okta_rate_limit" "example" {
  login                  = "ENFORCE"
  authorize              = "ENFORCE"
  communications_enabled = true
}
