resource "okta_collection" "test" {
  name = "testAcc-ds-collection-resource-replace_with_uuid"
}

resource "okta_collections_resource" "test" {
  collection_id = okta_collection.test.id
  resource_orn  = "orn:oktapreview:idp:00otivnuq1X5554yc1d7:apps:oidc_client:0oa15zse2by0vGzyP1d8"

  entitlements {
    id = "esp1mpfbbbqD2NAwG1d7"
    values {
      id = "ent1mpfbbczwILObX1d7"
    }
  }
}

data "okta_collections_resource" "test" {
  collection_id = okta_collection.test.id
  id            = okta_collections_resource.test.resource_id
}
