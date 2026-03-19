resource "okta_policy_mfa_default" "classic_example" {
  is_oie = false

  okta_password = {
    enroll = "REQUIRED"
  }

  okta_otp = {
    enroll = "REQUIRED"
  }
}

resource "okta_policy_mfa_default" "oie_example" {
  is_oie = true

  okta_password = {
    enroll = "REQUIRED"
  }

  # The following authenticator can only be used when `is_oie` is set to true
  okta_verify = {
    enroll = "REQUIRED"
  }
}

### -> If the `okta_policy_mfa_default` is used in conjunction with `okta_policy_mfa` resources, ensure to use a `depends_on` attribute for the default policy to ensure that all other policies are created/updated first such that the `priority` field can be appropriately computed on the first plan/apply.
