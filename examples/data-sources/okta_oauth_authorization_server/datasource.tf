data "okta_oauth_authorization_server" "test" {
  base_url = "https://${var.org_name}.${var.base_url}"
} 
