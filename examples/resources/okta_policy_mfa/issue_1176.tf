data "okta_group" "all" {
  name = "Everyone"
}

resource "okta_policy_mfa" "test" {
  name        = "testAcc_replace_with_uuid"
  status      = "ACTIVE"
  description = "Terraform Acceptance Test MFA Policy"
  is_oie      = true

  okta_otp = {
    enroll = "OPTIONAL"
  }

  #phone_number = {
  #  enroll = "OPTIONAL"
  #}

  okta_password = {
    enroll = "REQUIRED"
  }

  okta_email = {
    enroll = "OPTIONAL"
  }

  groups_included = [data.okta_group.all.id]
}
