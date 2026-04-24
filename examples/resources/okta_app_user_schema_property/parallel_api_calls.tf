resource "okta_app_oauth" "test" {
  label          = "testAcc_replace_with_uuid"
  type           = "native"
  grant_types    = ["authorization_code"]
  redirect_uris  = ["http://d.com/"]
  response_types = ["code"]
}
resource "okta_app_user_schema_property" "one" {
  app_id      = okta_app_oauth.test.id
  index       = "testAcc_replace_with_uuid_one"
  title       = "one"
  type        = "string"
  permissions = "%s"
}
resource "okta_app_user_schema_property" "two" {
  app_id      = okta_app_oauth.test.id
  index       = "testAcc_replace_with_uuid_two"
  title       = "two"
  type        = "string"
  permissions = "%s"
}
resource "okta_app_user_schema_property" "three" {
  app_id      = okta_app_oauth.test.id
  index       = "testAcc_replace_with_uuid_three"
  title       = "three"
  type        = "string"
  permissions = "%s"
}
resource "okta_app_user_schema_property" "four" {
  app_id      = okta_app_oauth.test.id
  index       = "testAcc_replace_with_uuid_four"
  title       = "four"
  type        = "string"
  permissions = "%s"
}
resource "okta_app_user_schema_property" "five" {
  app_id      = okta_app_oauth.test.id
  index       = "testAcc_replace_with_uuid_five"
  title       = "five"
  type        = "string"
  permissions = "%s"
}