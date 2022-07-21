data "okta_group" "all" {
  name = "Everyone"
}

resource "okta_policy_mfa" "test" {
  name        = "testAcc_replace_with_uuid"
  status      = "ACTIVE"
  description = "Terraform Acceptance Test MFA Policy"

  okta_email = {
    enroll = "REQUIRED"
  }

  google_otp = {
    enroll = "REQUIRED"
  }

  groups_included = [data.okta_group.all.id]
  depends_on      = [okta_factor.okta_email, okta_factor.google_otp]
}

resource "okta_factor" "okta_email" {
  provider_id = "okta_email"
}

resource "okta_factor" "google_otp" {
  provider_id = "google_otp"
}
