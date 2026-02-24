data "okta_push_groups" "sample" {
  app_id = "<okta_app_id>"
}

output "sample" {
  value = data.okta_push_groups.sample.mappings
}





