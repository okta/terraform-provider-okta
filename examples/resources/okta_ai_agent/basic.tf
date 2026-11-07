resource "okta_ai_agent" "test" {
  profile {
    name        = "test-ai-agent"
    description = "Test AI Agent"
  }

  resource_url = "https://example.com/api"
}
