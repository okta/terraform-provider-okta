data "okta_authenticator" "webauthn" {
  key = "webauthn"
}

data "okta_authenticator_method_webauthn" "example" {
  authenticator_id = data.okta_authenticator.webauthn.id
}

output "user_verification" {
  value = data.okta_authenticator_method_webauthn.example.user_verification
}

output "aaguid_groups" {
  value = data.okta_authenticator_method_webauthn.example.aaguid_groups
}
