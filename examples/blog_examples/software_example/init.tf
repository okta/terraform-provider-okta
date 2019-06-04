provider okta {
  org_name  = "example"
  base_url  = "okta.com"
  api_token = "${data.vault_generic_secret.okta_token.data["value"]}"
}

data vault_generic_secret okta_token {
  path = "secret/my_okta_api_token"
}
