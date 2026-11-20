data "okta_authenticator" "tac" {
  key = "tac"
}

data "okta_authenticator_tac" "test" {
  id = data.okta_authenticator.tac.id
}
