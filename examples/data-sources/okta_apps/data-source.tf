data "okta_apps" "example" {
  q                   = "Example"
  active_only         = true
  include_non_deleted = true
}