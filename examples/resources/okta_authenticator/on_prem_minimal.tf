resource "okta_authenticator" "test" {
  name                        = "On-Prem MFA"
  key                         = "onprem_mfa"
  provider_hostname           = "localhost"
  provider_auth_port          = 999
  provider_shared_secret      = "Sh4r3d s3cr3t"
  provider_user_name_template = "global.assign.userName.login"
}
