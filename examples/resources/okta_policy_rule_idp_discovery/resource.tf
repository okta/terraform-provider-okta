### All Okta orgs contain only one IdP Discovery Policy
data "okta_policy" "idp_discovery_policy" {
  name = "Idp Discovery Policy"
  type = "IDP_DISCOVERY"
}

# Example 1: Specific IdP routing - route to a named OIDC IdP
resource "okta_policy_rule_idp_discovery" "example" {
  policy_id                 = data.okta_policy.idp_discovery_policy.id
  name                      = "example"
  network_connection        = "ANYWHERE"
  priority                  = 1
  status                    = "ACTIVE"
  user_identifier_type      = "ATTRIBUTE"
  user_identifier_attribute = "company"

  idp_providers {
    id   = "<idp id>"
    type = "OIDC"
  }

  app_exclude {
    id   = "<app id>"
    type = "APP"
  }

  app_exclude {
    name = "yahoo_mail"
    type = "APP_TYPE"
  }

  app_include {
    id   = "<app id>"
    type = "APP"
  }

  app_include {
    name = "<app type name>"
    type = "APP_TYPE"
  }

  platform_include {
    type    = "MOBILE"
    os_type = "OSX"
  }

  user_identifier_patterns {
    match_type = "EQUALS"
    value      = "Articulate"
  }
}

# Example 2: Dynamic IdP routing - select IdP based on an expression
resource "okta_policy_rule_idp_discovery" "dynamic_example" {
  policy_id           = data.okta_policy.idp_discovery_policy.id
  name                = "dynamic-idp-routing"
  network_connection  = "ANYWHERE"
  priority            = 2
  status              = "ACTIVE"
  selection_type      = "DYNAMIC"
  provider_expression = "login.identifier.substringAfter('@')"
}
