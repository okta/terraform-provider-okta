data "okta_group" "all" {
  name = "Everyone"
}

resource "okta_policy_mfa" "test" {
  name        = "testAcc_replace_with_uuid"
  description = "Terraform Acceptance Test MFA Policy Yubikey Token"
  status      = "ACTIVE"
  is_oie      = true

  groups_included = [data.okta_group.all.id]

  okta_password = {
    enroll = "REQUIRED"
  }

  yubikey_token = {
    enroll = "REQUIRED"
  }
}
