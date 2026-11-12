# The Okta Account Management Policy's Catch-all Rule is the only rule that
# Okta sets to 2FA_If_Possible, and Okta does not allow it to be updated, so
# this fixture is only used for an import and plan.
data "okta_policy" "account_management" {
  name = "Okta Account Management Policy"
  type = "ACCESS_POLICY"
}

resource "okta_app_signon_policy_rules" "test" {
  policy_id = data.okta_policy.account_management.id

  rule {
    name                        = "Password Expiry Rule"
    priority                    = 0
    factor_mode                 = "1FA"
    re_authentication_frequency = "PT0S"
    custom_expression           = "accessRequest.operation=='recover' && accessRequest.metadata.type=='expiry'"
    constraints                 = [jsonencode({ knowledge = { required = true, types = ["password"] } })]
  }

  rule {
    name                        = "Catch-all Rule"
    priority                    = 99
    factor_mode                 = "2FA_If_Possible"
    re_authentication_frequency = "PT0S"
  }
}
