resource "okta_auth_server" "sun_also_rises" {
  audiences                 = ["whatever.rise.zone"]
  credentials_rotation_mode = "AUTO"
  description               = "The best way to find out if you can trust somebody is to trust them."
  name                      = "testAcc_%[1]d"
}
