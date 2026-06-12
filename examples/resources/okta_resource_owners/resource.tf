# Assign owners to an entitlement bundle
resource "okta_resource_owners" "example" {
  resource_orn = "orn:okta:governance:00o1234567890abcdef:entitlement-bundles:enb1234567890abcdef"

  principal_orns = [
    "orn:okta:directory:00o1234567890abcdef:users:00u1234567890abcdef",
    "orn:okta:directory:00o1234567890abcdef:groups:00g1234567890abcdef",
  ]
}
