resource "okta_entitlement_bundle" "test" {
  name        = "test-entitlement-bundle"
  description = "testing entitlement bundle"

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