resource "okta_mcp_server" "test" {
  registration_type = "MANUAL"
  resource_url      = "https://example.com"

  metadata {
    display_name = "Test MCP Server"
    description  = "Test MCP server registered by the terraform provider acceptance tests"
  }
}
