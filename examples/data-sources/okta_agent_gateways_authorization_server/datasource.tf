# Reads the authorization server linked to the org's existing "Production MCP Servers" agent gateway.
data "okta_agent_gateways_authorization_server" "test" {
  id               = "ausv9w2engBehiZ9tao5"
  agent_gateway_id = "wlpv9w2b57iZgZrnMao5"
}
