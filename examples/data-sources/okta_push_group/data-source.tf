data "okta_push_group" "sample" {
  app_id = "<okta_app_id>"

  id = "<push_group_mapping_id>"
}

data "okta_push_group" "another_sample" {
  app_id = "<okta_app_id>"

  source_group_id = "<okta_source_group_id>"
}

