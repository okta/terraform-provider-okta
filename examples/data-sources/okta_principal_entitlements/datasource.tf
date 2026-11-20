data "okta_principal_entitlements" "test" {
  parent {
    external_id = "0oa15zse2by0vGzyP1d8"
    type        = "APPLICATION"
  }
  target_principal {
    external_id = "00utivnuu1aaKzqqO1d7"
    type        = "OKTA_USER"
  }
}
