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
    permissions = "READ_ONLY"
    master      = "PROFILE_MASTER"
  }

  # This simulates an attribute auto-created by Okta when enabling provisioning.
  custom_property {
    index       = "autoAddedProp"
    title       = "Auto Added Property"
    type        = "string"
    description = "Auto added by Okta"
    required    = false
    scope       = "NONE"
    permissions = "READ_ONLY"
    master      = "PROFILE_MASTER"
  }
}

