resource "okta_app_oauth" "test" {
  label          = "testAcc_replace_with_uuid"
  type           = "native"
  grant_types    = ["authorization_code"]
  redirect_uris  = ["http://example.com/callback"]
  response_types = ["code"]
}

resource "okta_app_user_schema" "test" {
  app_id = okta_app_oauth.test.id

  custom_property {
    index       = "testCustomProp1"
    title       = "Test Custom Property 1"
    type        = "string"
    description = "Test description 1"
    required    = false
    scope       = "NONE"
    min_length  = 1
    max_length  = 20
    permissions = "READ_ONLY"
    master      = "PROFILE_MASTER"
  }

  custom_property {
    index       = "testCustomProp2"
    title       = "Test Custom Property 2"
    type        = "string"
    description = "Test description 2"
    required    = true
    scope       = "SELF"
    permissions = "READ_ONLY"
    master      = "OKTA"
  }

  custom_property {
    index       = "testArrayProp"
    title       = "Test Array Property"
    type        = "array"
    array_type  = "string"
    description = "Test array description"
    required    = false
    scope       = "NONE"
    permissions = "READ_WRITE"
    master      = "PROFILE_MASTER"
  }
}

