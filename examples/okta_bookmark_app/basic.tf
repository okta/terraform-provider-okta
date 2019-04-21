resource "okta_group" "group" {
  name = "testAcc_replace_with_uuid"
}

resource "okta_bookmark_app" "test" {
  label  = "testAcc_replace_with_uuid"
  url    = "https://test.com"
  groups = ["${okta_group.group.id}"]
}
