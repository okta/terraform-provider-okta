resource "okta_app_oauth" "my_app" {
  label                     = "My App"
  type                      = "web"
  grant_types               = ["authorization_code"]
  redirect_uris             = ["http://localhost:3000"]
  post_logout_redirect_uris = ["http://localhost:3000"]
  response_types            = ["code"]
  // this is needed to associate the application with the policy
  authentication_policy = okta_app_signon_policy.my_app_policy.id
}

resource "okta_app_signon_policy" "my_app_policy" {
  name        = "My App Sign-On Policy"
  description = "Authentication Policy to be used on my app."
}

### The created policy can be extended using `app_signon_policy_rules`.

resource "okta_app_signon_policy" "my_app_policy" {
  name        = "My App Sign-On Policy"
  description = "Authentication Policy to be used on my app."
}

resource "okta_app_signon_policy_rule" "some_rule" {
  policy_id                   = resource.okta_app_signon_policy.my_app_policy.id
  name                        = "Some Rule"
  factor_mode                 = "1FA"
  re_authentication_frequency = "PT43800H"
  constraints = [
    jsonencode({
      "knowledge" : {
        "types" : ["password"]
      }
    })
  ]
}
