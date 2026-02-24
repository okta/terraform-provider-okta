resource "okta_app_swa" "test" {
  app_links_json = jsonencode(
    {
      deelhr_link = true
    }
  )
  label                   = "Deel HR"
  preconfigured_app       = "deelhr"
  status                  = "ACTIVE"
  user_name_template_type = "BUILT_IN"
}

resource "okta_group" "test" {
  name = "Some Test Push Group"
}

resource "okta_push_group" "sample" {
  app_id                         = okta_app_swa.test.id
  source_group_id                = okta_group.test.id
  status                         = "ACTIVE"
  delete_target_group_on_destroy = true
}

resource "okta_push_group" "ad_sample" {
  app_id          = okta_app_swa.test.id
  source_group_id = okta_group.test.id
  status          = "ACTIVE"
  app_config = {
    distinguished_name = "CN=Test,OU=Groups,DC=example,DC=com"
    group_scope        = "DOMAIN_LOCAL"
    group_type         = "SECURITY"
    sam_account_name   = "something"
  }
}
