resource "okta_link_definition" "test" {
  primary_name           = "testAcc_replace_with_uuid"
  primary_title          = "Manager"
  primary_description    = "Manager link property"
  associated_name        = "testAcc_subordinate"
  associated_title       = "Subordinate"
  associated_description = "Subordinate link property"
}