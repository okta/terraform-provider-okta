variable "hostname" {
  type = string
}

data "okta_oauth_authorization_server" "test" {
  base_url = "https://${var.hostname}"
}
