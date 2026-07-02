resource "okta_app_swa" "sample" {
  app_links_json = jsonencode(
    {
      deelhr_link = true
    }
  )
  label                   = "Test Deel HR"
  preconfigured_app       = "deelhr"
  status                  = "ACTIVE"
  user_name_template_type = "BUILT_IN"
}

resource "okta_group" "sample" {
  name = "Some Example Test Push Group"
}

resource "okta_push_group" "sample" {
  app_id                         = okta_app_swa.sample.id
  source_group_id                = okta_group.sample.id
  status                         = "INACTIVE"
  delete_target_group_on_destroy = false
}
