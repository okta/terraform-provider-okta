resource "okta_entitlement_bundle" "test" {
  name        = "entitlement bundle data source test"
  description = "test resource for entitlement bundle data source"

  target {
    external_id = "0oao01ardu8r8qUP91d7"
    type        = "APPLICATION"
  }

  entitlements {
    id = "espzcbqd7Suwp4Y7A1d6"

    values {
      id = "entzcbqd8lcD3BRWR1d6"
    }
  }
}

data "okta_entitlement_bundle" "test" {
  id = okta_entitlement_bundle.test.id
}