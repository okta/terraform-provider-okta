data "okta_group" "all" {
  name = "Everyone"
}

resource "okta_policy_mfa" "test" {
  name            = "testAcc_replace_with_uuid_new"
  status          = "INACTIVE"
  description     = "Terraform Acceptance Test MFA Policy Updated"
  groups_included = [data.okta_group.all.id]

  okta_email = {
    enroll = "REQUIRED"
  }

  okta_password = {
    enroll = "OPTIONAL"
  }

  google_otp = {
    enroll = "OPTIONAL"
  }

  depends_on = [okta_factor.google_otp, okta_factor.okta_email, okta_factor.okta_password]
}

resource "okta_factor" "okta_email" {
  provider_id = "okta_email"
}

resource "okta_factor" "okta_password" {
  provider_id = "okta_password"
}


resource "okta_factor" "google_otp" {
  provider_id = "google_otp"
}
