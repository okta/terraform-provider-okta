resource "okta_ai_agent" "test" {
  resource_url = "https://my-agent.test.com/test"

  profile {
    name        = "my-ai-agent-test"
    description = "Example AI agent registration-update"
  }
}