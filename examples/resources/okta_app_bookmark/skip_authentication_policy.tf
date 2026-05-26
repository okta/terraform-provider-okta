# Example: Skip Authentication Policy Operations
# This example demonstrates how to use the skip_authentication_policy flag
# to prevent the provider from managing authentication policies for an application.

resource "okta_app_bookmark" "example_skip_policy" {
  label = "Example App - Skip Policy"
  url   = "https://example.com"

  # Skip authentication policy operations
  # When set to true, the provider will not attempt to create, update, or delete
  # authentication policies for this application. This is useful when you want
  # to manage authentication policies manually or when the application should
  # use the default policy without explicit configuration.
  skip_authentication_policy = true
}

# Example: Regular app with authentication policy management
resource "okta_app_bookmark" "example_with_policy" {
  label = "Example App - With Policy"
  url   = "https://example.com"

  # This app will have authentication policy operations performed normally
  # The provider will assign it to the default policy if none is specified
  skip_authentication_policy = false
}
