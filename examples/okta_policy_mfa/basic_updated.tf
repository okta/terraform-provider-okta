data okta_group all {
  name = "Everyone"
}

resource okta_policy_mfa test {
  name            = "testAcc_replace_with_uuid"
  status          = "INACTIVE"
  description     = "Terraform Acceptance Test MFA Policy Updated"
  groups_included = ["${data.okta_group.all.id}"]

  fido_u2f = {
    enroll = "OPTIONAL"
  }

  google_otp = {
    enroll = "OPTIONAL"
  }

  okta_otp = {
    enroll = "OPTIONAL"
  }

  okta_sms = {
    enroll = "OPTIONAL"
  }
}
