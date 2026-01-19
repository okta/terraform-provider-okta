data "okta_principal_entitlements" "test" {
  parent {
    external_id = "0oao01ardu8r8qUP91d7"
    type        = "APPLICATION"
  }
  target_principal {
    external_id = "00unkw1sfbTw08c0g1d7"
    type        = "OKTA_USER"
  }
}
