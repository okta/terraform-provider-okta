data "okta_group" "all" {
  name = "Everyone"
}

resource "okta_policy_mfa" "test" {
  name        = "testAcc_replace_with_uuid"
  status      = "ACTIVE"
  description = "Terraform Acceptance Test MFA Policy"
  priority    = 1
  is_oie      = true

  okta_password = {
    enroll = "REQUIRED"
  }

  okta_email = {
    enroll = "NOT_ALLOWED"
  }

  fido_webauthn = {
    enroll = "REQUIRED"
  }

  groups_included = [data.okta_group.all.id]
}
