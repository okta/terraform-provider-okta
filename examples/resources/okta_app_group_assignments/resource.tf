resource "okta_app_group_assignments" "example" {
  app_id = "<app id>"
  group {
    id       = "<group id>"
    priority = 1
  }
  group {
    id       = "<another group id>"
    priority = 2
    profile  = jsonencode({ "application profile field" : "application profile value" })
  }
}
