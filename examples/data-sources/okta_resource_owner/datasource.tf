resource "okta_resource_owner" "test" {
  resource_orns = ["orn:oktapreview:idp:00onkw1sa89bQ0bncd9A:apps:oidc_client:0oabcde0afs9cbx90bup"]
}

data "okta_resource_owner" "test" {
  parent_resource_orn = "orn:oktapreview:idp:00onabc90Ah3T09I1xe0:apps:oidc_client:0oanbc98zc6kh2bx32Pv"
  depends_on          = [okta_resource_owner.test]
}
