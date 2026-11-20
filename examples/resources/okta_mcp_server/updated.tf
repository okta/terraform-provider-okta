resource "okta_mcp_server" "test" {
  resource_url = "https://workspacemcp.googleapis.com/mcp/v1"

  metadata {
    display_name = "MCP Server updated"
    description  = "Test MCP server registered by the terraform provider acceptance tests, updated"
  }
}
