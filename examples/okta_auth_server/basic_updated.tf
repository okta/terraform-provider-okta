resource "okta_auth_server" "sun_also_rises" {
  audiences   = ["api://something_else"]
  description = "The past is not dead. In fact, it's not even past."
  name        = "Cheeky Faulkner quip %[1]d"
}
