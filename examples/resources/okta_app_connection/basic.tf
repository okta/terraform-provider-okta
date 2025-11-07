resource "okta_app_connection" "example" {
  id       = "0oarblaf7hWdLawNg1d7"
  base_url = "https://integrator-8357679-admin.okta.com"
  action   = "activate"
  profile {
    auth_scheme = "TOKEN"
    token       = "004Ah8sGzUBj_38FZRNwnikXgrxi33pu7yRrWII4DB"
  }
}