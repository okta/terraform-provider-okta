resource "okta_collection" "test" {
  name        = "Test Collection"
  description = "Test collection"
}

resource "okta_collection_resource" "test" {
  collection_id = okta_collection.test.id
  resource_orn  = "orn:okta:idp:00o123:apps:salesforce:0oa456"
  
  entitlements {
    id = "ent123"
    values {
      id = "val456"
    }
  }
  
  entitlements {
    id = "ent789"
    values {
      id = "val012"
    }
  }
}
