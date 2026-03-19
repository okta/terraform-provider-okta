resource "okta_admin_role_targets" "example" {
  user_id   = "<user_id>"
  role_type = "APP_ADMIN"
  apps      = ["oidc_client.<app_id>", "facebook"]
}
