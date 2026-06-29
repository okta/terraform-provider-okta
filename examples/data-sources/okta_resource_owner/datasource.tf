resource "okta_resource_owner" "test" {
  resource_orns = ["orn:oktapreview:idp:00onkw1sbuAh3Q06I1d7:apps:oidc_client:0oanlpn0tzc7kh9bx1d7"]
}

data "okta_resource_owner" "test" {
  parent_resource_orn = "orn:oktapreview:idp:00onkw1sbuAh3Q06I1d7:apps:oidc_client:0oanlpn0tzc7kh9bx1d7"
  depends_on          = [okta_resource_owner.test]
}
