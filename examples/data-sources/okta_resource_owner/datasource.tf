resource "okta_resource_owner" "test" {
  resource_orns = ["orn:oktapreview:idp:00otivnuq1X5554yc1d7:apps:oidc_client:0oaws4am895IZbn6Q1d7"]
}

data "okta_resource_owner" "test" {
  parent_resource_orn = "orn:oktapreview:idp:00onabc90Ah3T09I1xe0:apps:oidc_client:0oaws4am895IZbn6Q1d7"
  depends_on          = [okta_resource_owner.test]
}
