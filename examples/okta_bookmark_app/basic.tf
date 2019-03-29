resource "okta_group" "group" {
  name = "testAcc_%[1]d"
}

resource "okta_bookmark_app" "test" {
  label  = "testAcc_%[1]d"
  url    = "https://test.com"
  groups = ["${okta_group.group.id}"]
}
