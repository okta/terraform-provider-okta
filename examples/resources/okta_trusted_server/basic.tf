resource "okta_auth_server" "test1" {
  audiences                 = ["whatever.rise.zone"]
  credentials_rotation_mode = "AUTO"
  description               = "The best way to find out if you can trust somebody is to trust them."
  name                      = "testAcc-replace_with_uuid"
}

resource "okta_auth_server" "test2" {
  audiences                 = ["whatever.rise.zone"]
  credentials_rotation_mode = "AUTO"
  description               = "The best way to find out if you can trust somebody is to trust them."
  name                      = "testAcc-replace_with_uuid"
}

resource "okta_auth_server" "test3" {
  audiences                 = ["whatever.rise.zone"]
  credentials_rotation_mode = "AUTO"
  description               = "The best way to find out if you can trust somebody is to trust them."
  name                      = "testAcc-replace_with_uuid"
}

resource "okta_auth_server" "test4" {
  audiences                 = ["whatever.rise.zone"]
  credentials_rotation_mode = "AUTO"
  description               = "The best way to find out if you can trust somebody is to trust them."
  name                      = "testAcc-replace_with_uuid"
}

resource "okta_trusted_server" "example" {
  auth_server_id = resource.okta_auth_server.test1.id
  trusted        = [resource.okta_auth_server.test2.id, resource.okta_auth_server.test3.id]
}
