resource "okta_app_signon_policy_rule" "test_with_reauthenticate_in_chains_only" {
  policy_id          = "rsttqnoz5vo4GIoAD1d7"
  name               = "test_with_reauthenticate_in_chains_only"
  type               = "AUTH_METHOD_CHAIN"
  priority           = 3
  network_connection = "ANYWHERE"
  access             = "ALLOW"
  factor_mode        = "2FA"
  chains = [
    jsonencode(
      {
        "authenticationMethods" : [
          {
            "key" : "webauthn",
            "userVerification" : "REQUIRED",
            "method" : "webauthn"
          }
        ],
        "reauthenticateIn" : "PT43800H"
    }),
    jsonencode(
      {
        "authenticationMethods" : [
          {
            "key" : "okta_password",
            "method" : "password"
          }
        ],
        "reauthenticateIn" : "PT43800H",
        "next" : [
          {
            "authenticationMethods" : [
              {
                "key" : "custom_otp",
                "method" : "otp"
              },
              {
                "key" : "google_otp",
                "method" : "otp"
              },
              {
                "key" : "phone_number",
                "method" : "sms"
              },
              {
                "key" : "phone_number",
                "method" : "voice"
              },
              {
                "key" : "webauthn",
                "userVerification" : "REQUIRED",
                "method" : "webauthn"
              }
            ],
            "reauthenticateIn" : "PT43800H"
          }
        ]
      }
    )
  ]
}


resource "okta_app_signon_policy_rule" "test_with_re_authentication_frequency_only" {
  policy_id                   = "rsttqnoz5vo4GIoAD1d7"
  name                        = "test_with_re_authentication_frequency_only"
  type                        = "AUTH_METHOD_CHAIN"
  priority                    = 4
  network_connection          = "ANYWHERE"
  access                      = "ALLOW"
  factor_mode                 = "2FA"
  re_authentication_frequency = "PT2H10M"
  inactivity_period           = "PT1H"
  chains = [
    jsonencode({
      "authenticationMethods" : [
        {
          "key" : "okta_password",
          "method" : "password"
        }
      ],
      "next" : [
        {
          "authenticationMethods" : [
            {
              "key" : "phone_number",
              "method" : "sms"
            }
          ]
        }
      ]
    }),
    jsonencode({
      "authenticationMethods" : [
        {
          "key" : "okta_verify",
          "method" : "signed_nonce",
          "userVerification" : "REQUIRED"
        }
      ]
    })
  ]
}