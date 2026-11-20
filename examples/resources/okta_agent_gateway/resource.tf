resource "okta_agent_gateway" "test" {
  resource_path = "staging"

  profile {
    display_name = "Production MCP Servers"
    description  = "description for the production MCP server group"
  }
}
