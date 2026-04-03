# Example: Skip Authentication Policy Operations for OAuth App
# This example demonstrates how to use the skip_authentication_policy flag
# to prevent the provider from managing authentication policies for an OAuth application.

resource "okta_app_oauth" "example_skip_policy" {
  label          = "Example OAuth App - Skip Policy"
  type           = "web"
  grant_types    = ["authorization_code"]
  response_types = ["code"]
  redirect_uris  = ["https://example.com/callback"]

  # Skip authentication policy operations
  # When set to true, the provider will not attempt to create, update, or delete
  # authentication policies for this application. This is useful when you want
  # to manage authentication policies manually or when the application should
  # use the default policy without explicit configuration.
  skip_authentication_policy = true
}

# Example: Regular OAuth app with authentication policy management
resource "okta_app_oauth" "example_with_policy" {
  label          = "Example OAuth App - With Policy"
  type           = "web"
  grant_types    = ["authorization_code"]
  response_types = ["code"]
  redirect_uris  = ["https://example.com/callback"]

  # This app will have authentication policy operations performed normally
  # The provider will assign it to the default policy if none is specified
  skip_authentication_policy = false
}
