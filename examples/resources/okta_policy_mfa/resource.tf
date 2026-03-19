
resource "okta_policy_mfa" "classic_example" {
  name        = "MFA Policy Classic"
  status      = "ACTIVE"
  description = "Example MFA policy using Okta Classic engine with factors."
  is_oie      = false

  okta_password = {
    enroll = "REQUIRED"
  }

  okta_otp = {
    enroll = "REQUIRED"
  }

  groups_included = ["${data.okta_group.everyone.id}"]
}

resource "okta_policy_mfa" "oie_example" {
  name        = "MFA Policy OIE"
  status      = "ACTIVE"
  description = "Example MFA policy that uses Okta Identity Engine (OIE) with authenticators"
  is_oie      = true

  okta_password = {
    enroll = "REQUIRED"
  }

  # The following authenticator can only be used when `is_oie` is set to true
  okta_verify = {
    enroll = "REQUIRED"
  }

  groups_included = ["${data.okta_group.everyone.id}"]
}

// policy_mfa that support multiple external_idps
resource "okta_policy_mfa" "oie_example" {
  name        = "MFA Policy OIE"
  status      = "ACTIVE"
  description = "Example MFA policy that uses Okta Identity Engine (OIE) with authenticators"
  is_oie      = true

  okta_email = {
    enroll = "REQUIRED"
  }

  google_otp = {
    enroll = "REQUIRED"
  }

  external_idps = [
    {
      "enroll" : "OPTIONAL",
      "id" : "id1",
    },
    {
      "enroll" : "NOT_ALLOWED",
      "id" : "id2",
    }
  ]


  groups_included = ["${data.okta_group.everyone.id}"]
}
