data "okta_group" "all" {
  name = "Everyone"
}

resource "okta_policy_mfa" "test" {
  name        = "testAcc_replace_with_uuid"
  status      = "ACTIVE"
  description = "Terraform Acceptance Test MFA Policy Okta Verify Breakdown"
  is_oie      = true

  groups_included = [data.okta_group.all.id]

  okta_password = {
    enroll = "REQUIRED"
  }

  okta_verify_fastpass = {
    enroll = "REQUIRED"
  }

  okta_verify_push = {
    enroll = "OPTIONAL"
  }

  okta_verify_totp = {
    enroll = "NOT_ALLOWED"
  }
}
