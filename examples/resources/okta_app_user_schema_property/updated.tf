resource "okta_app_oauth" "test" {
  label          = "testLabel"
  type           = "native"
  grant_types    = ["authorization_code"]
  redirect_uris  = ["http://d.com/"]
  response_types = ["code"]
}

resource "okta_app_user_schema_property" "test" {
  app_id = okta_app_oauth.test.id
  # need to rename "index" since ForceNew:true for it, refer user_schema_property.go 
  # also, API has restrictions refer to https://support.okta.com/help/s/article/unable-to-create-new-attribute-in-the-profile-editor?language=en_US
  index       = "testIndex_renamed"
  title       = "terraform acceptance test"
  type        = "string"
  description = "terraform acceptance test updated 001"
  required    = true
  master      = "OKTA"
  scope       = "SELF"
}
