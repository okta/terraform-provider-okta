resource "okta_ai_agent" "test" {
  profile {
    name        = "test-ai-agent-updated"
    description = "Test AI Agent, updated"
  }
  resource_url = "https://testacc-ai-agent-replace_with_uuid.example.com/api"
}
