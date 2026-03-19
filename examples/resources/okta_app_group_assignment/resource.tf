resource "okta_app_group_assignment" "example" {
  app_id   = "<app id>"
  group_id = "<group id>"
  profile  = <<JSON
{
  "<app_profile_field>": "<value>"
}
JSON
}
