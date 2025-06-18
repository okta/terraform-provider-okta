resource "okta_app_signon_policy" "test" {
  name        = "testAcc_Test_App_replace_with_uuid"
  description = "The app signon policy used by our test app"
}

resource "okta_app_signon_policy_rule" "test" {
  access                      = "ALLOW"
  custom_expression           = null
  device_assurances_included  = null
  device_is_managed           = null
  device_is_registered        = null
  factor_mode                 = "2FA"
  groups_excluded             = null
  groups_included             = null
  inactivity_period           = "PT1H"
  name                        = "test2"
  network_connection          = "ANYWHERE"
  network_excludes            = null
  network_includes            = null
  policy_id                   = okta_app_signon_policy.test.id
  priority                    = 0
  re_authentication_frequency = "PT0S"
  status                      = "ACTIVE"
  type                        = "AUTH_METHOD_CHAIN"
  user_types_excluded         = []
  user_types_included         = []
  users_excluded              = []
  users_included              = []
  platform_include {
    os_expression = ""
    os_type       = "OTHER"
    type          = "DESKTOP"
  }
  chains = [jsonencode(
    {
      "authenticationMethods" : [
        {
          "key" : "okta_password",
          "method" : "password"
        }
      ],
      "next" : [{
        "authenticationMethods" : [{
          "key" : "okta_email",
          "method" : "email"
        }]
      }]
    }
  )]
}