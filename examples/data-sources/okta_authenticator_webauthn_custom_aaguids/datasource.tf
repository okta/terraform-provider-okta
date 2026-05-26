data "okta_authenticator" "webauthn" {
  key = "webauthn"
}

data "okta_authenticator_webauthn_custom_aaguids" "all" {
  authenticator_id = data.okta_authenticator.webauthn.id
}

output "custom_aaguids" {
  value = data.okta_authenticator_webauthn_custom_aaguids.all.custom_aaguids
}
