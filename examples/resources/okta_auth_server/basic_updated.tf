resource "okta_auth_server" "sun_also_rises" {
  audiences   = ["whatever-else.rise.zone"]
  description = "The past is not dead. In fact, it's not even past."
  name        = "testAcc_replace_with_uuid"
  issuer_mode = "DYNAMIC"
}
