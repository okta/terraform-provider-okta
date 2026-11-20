resource "okta_agent_gateway" "test" {
  resource_path = "testacc-replace_with_uuid"

  profile {
    display_name = "testAcc-agent-gateway-updated-replace_with_uuid"
    description  = "Terraform acceptance test agent gateway, updated"
  }
}
