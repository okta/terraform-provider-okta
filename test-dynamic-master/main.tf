resource "okta_network_zone" "corp_vpn" {
  name     = "Corp VPN"
  type     = "IP"
  status   = "INACTIVE"
  gateways = ["10.0.0.0/8"]
}

locals {
  policy_rules = {
    "allow_1fa" = {
      name               = "Allow 1FA"
      priority           = 1
      status             = "ACTIVE"
      access             = "ALLOW"
      factor_mode        = "1FA"
      network_connection = "ZONE"
      network_includes   = [okta_network_zone.corp_vpn.id]
      chains             = []
      platform_include   = []
    }
    "allow_2fa_vpn" = {
      name               = "Allow 2FA from VPN - UPDATED NAME"
      priority           = 2
      status             = "ACTIVE"
      access             = "ALLOW"
      factor_mode        = "1FA"
      network_connection = "ZONE"
      network_includes   = [okta_network_zone.corp_vpn.id] # unknown at plan time
      platform_include = [{
        os_type = "ANDROID"
        type    = "MOBILE"
      }]
      chains = [jsonencode(
        {
          "authenticationMethods" : [
            {
              "key" : "okta_email",
              "method" : "email"
            }
          ],
          "next" : [{
            "authenticationMethods" : [{
              "key" : "okta_password",
              "method" : "password"
            }]
          }]
        }
      )]
    }
    "new_rule" = {
      name               = "New Rule Added"
      priority           = 3
      status             = "ACTIVE"
      access             = "ALLOW"
      factor_mode        = "1FA"
      network_connection = "ANYWHERE"
      network_includes   = []
      chains             = []
      platform_include   = []
    }
  }
}

resource "okta_app_signon_policy" "dynamic_example" {
  name        = "Dynamic App Sign-On Policy"
  description = "Policy with dynamically defined rules"
}

resource "okta_app_signon_policy_rules" "dynamic_example" {
  policy_id = okta_app_signon_policy.dynamic_example.id

  dynamic "rule" {
    for_each = local.policy_rules
    content {
      name               = rule.value.name
      priority           = rule.value.priority
      status             = rule.value.status
      access             = rule.value.access
      factor_mode        = rule.value.factor_mode
      network_connection = rule.value.network_connection
      network_includes   = rule.value.network_includes
    }
  }
}
