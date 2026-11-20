resource "okta_collection" "test" {
  name        = "testAcc-collection-resource-replace_with_uuid"
  description = "Collection for a collection-resource test"
}

resource "okta_collections_resource" "test" {
  collection_id = okta_collection.test.id
  resource_orn  = "orn:oktapreview:idp:00otivnuq1X5554yc1d7:apps:oidc_client:0oa15zse2by0vGzyP1d8"

  entitlements {
    id = "esp1mrruxopcczTEr1d7"
    values {
      id = "ent1mrruxps9YxXFb1d7"
    }
  }
}
