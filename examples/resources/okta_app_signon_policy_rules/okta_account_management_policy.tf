resource "okta_app_signon_policy_rules" "okta_app_signon_policy_rules" {
  policy_id = "rstj21rq8pboPiddo1d7"

  rule {
    name                        = "Password Expiry Rule"
    priority                    = 0
    status                      = "ACTIVE"
    access                      = "ALLOW"
    factor_mode                 = "1FA"
    type                        = "ASSURANCE"
    network_connection          = "ANYWHERE"
    custom_expression           = "accessRequest.operation=='recover' && accessRequest.metadata.type=='expiry'"
    re_authentication_frequency = "PT0S"
    risk_score                  = "ANY"

    constraints = [
      jsonencode({
        knowledge = {
          required = true
          types = [
            "password",
          ]
        }
      })
    ]
  }

  rule {
    name                        = "Account management policy rule"
    priority                    = 1
    status                      = "ACTIVE"
    access                      = "ALLOW"
    factor_mode                 = "2FA"
    type                        = "AUTH_METHOD_CHAIN"
    network_connection          = "ANYWHERE"
    risk_score                  = "ANY"
    groups_included             = ["00gimgbwy8cKIj5r61d7"]
    re_authentication_frequency = null

    chains = [
      jsonencode({
        authenticationMethods = [
          {
            key    = "phone_number"
            method = "sms"
          },
          {
            key    = "phone_number"
            method = "voice"
          },
        ]
        next = [
          {
            authenticationMethods = [
              {
                key    = "google_otp"
                method = "otp"
              },
              {
                key              = "okta_verify"
                userVerification = "OPTIONAL"
                method           = "push"
              },
              {
                key    = "okta_verify"
                method = "totp"
              },
              {
                key              = "okta_verify"
                userVerification = "OPTIONAL"
                method           = "signed_nonce"
              },
              {
                key              = "webauthn"
                userVerification = "OPTIONAL"
                method           = "webauthn"
              },
              {
                key    = "okta_password"
                method = "password"
              },
              {
                key    = "okta_email"
                method = "email"
              },
            ]
            reauthenticateIn = "PT0S"
          },
        ]
        reauthenticateIn = "PT0S"
      }),
      jsonencode({
        authenticationMethods = [
          {
            key    = "okta_password"
            method = "password"
          },
        ]
        next = [
          {
            authenticationMethods = [
              {
                key    = "google_otp"
                method = "otp"
              },
              {
                key              = "okta_verify"
                userVerification = "OPTIONAL"
                method           = "push"
              },
              {
                key    = "okta_verify"
                method = "totp"
              },
              {
                key              = "okta_verify"
                userVerification = "OPTIONAL"
                method           = "signed_nonce"
              },
              {
                key              = "webauthn"
                userVerification = "OPTIONAL"
                method           = "webauthn"
              },
              {
                key    = "phone_number"
                method = "sms"
              },
              {
                key    = "phone_number"
                method = "voice"
              },
              {
                key    = "okta_email"
                method = "email"
              },
            ]
            reauthenticateIn = "PT0S"
          },
        ]
        reauthenticateIn = "PT0S"
      }),
    ]
  }
}
