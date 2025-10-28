resource "okta_app_oauth" "example" {
  label = "example"
  type  = "web"
}

data "okta_apps" "apps" {
  q = "example"
}