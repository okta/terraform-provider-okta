data "okta_group" "all" {
  name = "Everyone"
}

resource "okta_policy_mfa" "test" {
  name            = "testAcc_replace_with_uuid"
  status          = "INACTIVE"
  description     = "Terraform Acceptance Test MFA Policy Updated"
  groups_included = [data.okta_group.all.id]

  google_otp = {
    enroll = "OPTIONAL"
  }

  okta_sms = {
    enroll = "OPTIONAL"
  }

  depends_on = [okta_factor.google_otp, okta_factor.okta_sms]
}

resource "okta_factor" "google_otp" {
  provider_id = "google_otp"
}

resource "okta_factor" "okta_sms" {
  provider_id = "okta_sms"
}
