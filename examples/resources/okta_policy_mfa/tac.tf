data "okta_group" "all" {
  name = "Everyone"
}

# A standard authentication authenticator is activated alongside TAC: TAC is a
# temporary-access authenticator, so Okta rejects a policy that would leave users
# without another multifactor authenticator for authentication. google_otp is a
# self-contained choice (no telephony/email-mode dependencies).
resource "okta_authenticator" "google_otp" {
  name   = "Google Authenticator"
  key    = "google_otp"
  status = "ACTIVE"
}

resource "okta_authenticator" "tac" {
  name   = "Temporary Access Code"
  key    = "tac"
  status = "ACTIVE"
  provider_json = jsonencode({
    type = "TAC" // provider type must be uppercase "TAC"
    configuration = {
      length          = 8
      minTtl          = 10
      maxTtl          = 480
      defaultTtl      = 60
      multiUseAllowed = false
      complexity = {
        numbers           = true
        letters           = true
        specialCharacters = false
      }
    }
  })
}

resource "okta_policy_mfa" "test" {
  name        = "testAcc_replace_with_uuid"
  description = "Terraform Acceptance Test MFA Policy TAC"
  status      = "ACTIVE"
  is_oie      = true

  groups_included = [data.okta_group.all.id]

  okta_password = {
    enroll = "REQUIRED"
  }

  google_otp = {
    enroll = "REQUIRED"
  }

  tac = {
    enroll = "OPTIONAL"
  }

  depends_on = [okta_authenticator.tac, okta_authenticator.google_otp]
}
