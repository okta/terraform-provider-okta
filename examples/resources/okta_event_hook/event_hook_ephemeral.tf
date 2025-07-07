# Example configuration showing ephemeral support for okta_event_hook

# Using ephemeral values from AWS Secrets Manager
data "aws_secretsmanager_secret_version" "webhook_token" {
  secret_id = "webhook-api-token"
}

resource "okta_event_hook" "example" {
  name   = "Example Event Hook with Ephemeral Auth"
  status = "ACTIVE"
  events = [
    "user.lifecycle.create",
    "user.lifecycle.activate",
    "user.lifecycle.deactivate"
  ]

  channel = {
    type    = "HTTP"
    version = "1.0.0"
    uri     = "https://example.com/webhook"
  }

  # The auth.value field now supports ephemeral values
  # This value will not be stored in Terraform state
  auth = {
    type  = "HEADER"
    key   = "x-api-token"
    value = ephemeral.aws_secretsmanager_secret_version.webhook_token.secret_string
  }
}

# Alternative: Using a regular secret (still write-only)
resource "okta_event_hook" "simple" {
  name   = "Simple Event Hook"
  status = "ACTIVE"
  events = ["user.lifecycle.create"]

  channel = {
    type    = "HTTP"
    version = "1.0.0"
    uri     = "https://example.com/simple-webhook"
  }

  # Even with regular values, the auth.value is write-only
  # and won't be stored in state
  auth = {
    type  = "HEADER"
    key   = "Authorization"
    value = "Bearer your-secret-token-here"
  }
}
