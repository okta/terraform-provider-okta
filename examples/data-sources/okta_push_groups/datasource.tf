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
  name = "Some Example Testing Push Group"
}

resource "okta_group" "sample_two" {
  name = "Some Other Example Testing Push Group"
}

resource "okta_push_group" "sample" {
  app_id          = okta_app_swa.sample.id
  source_group_id = okta_group.sample.id
}

resource "okta_push_group" "sample_two" {
  app_id          = okta_app_swa.sample.id
  source_group_id = okta_group.sample_two.id
}

data "okta_push_groups" "sample" {
  app_id = okta_app_swa.sample.id
}
