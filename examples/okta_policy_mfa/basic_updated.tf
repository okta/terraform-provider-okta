data "okta_group" "all" {
  name = "Everyone"
}

resource "okta_policy_mfa" "test" {
  name            = "testAcc_replace_with_uuid_new"
  status          = "INACTIVE"
  description     = "Terraform Acceptance Test MFA Policy Updated"
  groups_included = [data.okta_group.all.id]

  google_otp = {
    enroll = "OPTIONAL"
  }

  okta_sms = {
    enroll = "OPTIONAL"
  }

  okta_email = {
    enroll = "OPTIONAL"
  }

  hotp = {
    enroll = "OPTIONAL"
  }

  depends_on = [okta_factor.google_otp, okta_factor.okta_sms, okta_factor.okta_email, okta_factor.hotp]
}

resource "okta_factor" "google_otp" {
  provider_id = "google_otp"
}

resource "okta_factor" "okta_sms" {
  provider_id = "okta_sms"
}

resource "okta_factor" "okta_email" {
  provider_id = "okta_email"
}

resource "okta_factor" "hotp" {
  provider_id = "hotp"
}
