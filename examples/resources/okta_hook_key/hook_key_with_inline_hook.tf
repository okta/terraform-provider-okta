# Create a Hook Key
resource "okta_hook_key" "example" {
  name = "My Hook Key for Inline Hook"
}

# Create an Inline Hook that uses the Hook Key
resource "okta_inline_hook" "example" {
  name    = "My Inline Hook"
  status  = "ACTIVE"
  type    = "com.okta.oauth2.tokens.transform"
  version = "1.0.0"

  channel = {
    type    = "HTTP"
    version = "1.0.0"
    uri     = "https://example.com/hook"
    method  = "POST"
    
    # Reference the Hook Key name for authentication
    auth_scheme = {
      type  = "HEADER"
      key   = "Authorization"
      value = "Bearer ${okta_hook_key.example.name}"
    }
  }
}