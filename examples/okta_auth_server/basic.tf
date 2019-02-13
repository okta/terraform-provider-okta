resource "okta_auth_server" "sun_also_rises" {
  audiences = ["api://something"]

  credentials {
    rotation_mode = "MANUAL"
  }

  description = "The best way to find out if you can trust somebody is to trust them."
  name        = "Cheeky Hemingway quip %[1]d"
}
