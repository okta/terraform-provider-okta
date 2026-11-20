resource "okta_ai_agent" "test" {
  profile {
    name        = "test-ai-agent"
    description = "Test AI Agent"
  }
  resource_url = "https://testacc-ai-agent-replace_with_uuid.example.com/api"
}
