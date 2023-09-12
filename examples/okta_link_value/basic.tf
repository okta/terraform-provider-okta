resource "okta_link_definition" "test" {
  primary_name           = "testAcc_replace_with_uuid"
  primary_title          = "Manager"
  primary_description    = "Manager link property"
  associated_name        = "testAcc_subordinate"
  associated_title       = "Subordinate"
  associated_description = "Subordinate link property"
}

resource "okta_user" "test" {
  count      = 5
  first_name = "TestAcc"
  last_name  = "Smith"
  login      = "testAcc_${count.index}@example.com"
  email      = "testAcc_${count.index}@example.com"
}

resource "okta_link_value" "test" {
  primary_name    = okta_link_definition.test.primary_name
  primary_user_id = okta_user.test[0].id
  associated_user_ids = [
    okta_user.test[1].id,
    okta_user.test[2].id,
    okta_user.test[3].id,
    okta_user.test[4].id,
  ]
}