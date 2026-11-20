# Uses the org's existing Microsoft Teams OIN MCP server and its linked authorization server.
# An org can register an OIN MCP server only once, so the test can't create its own.
resource "okta_ai_agent" "test" {
  profile {
    name        = "testAcc-connection-replace_with_uuid"
    description = "Terraform acceptance test agent for okta_connection"
  }

  resource_url = "https://testacc-connection-replace_with_uuid.example.com/api"
}

resource "okta_connection" "test" {
  actor {
    wlp_id = okta_ai_agent.test.id
  }

  authorization_server {
    orn = "orn:okta:idp:00otxk5i5T0Rk59Z3ao5:authorization_servers:ausu6ayw5WeizA26Vao5"
  }

  protocol {
    protocol_type   = "IDENTITY_ASSERTION"
    scope_condition = "INCLUDE_ONLY"
    scopes          = ["openid", "profile"]
  }

  resource {
    resource_type = "MCP_SERVER"
    orn           = "orn:okta:directory:00otxk5i5T0Rk59Z3ao5:resource-servers:mcp:emsv9wtzWEC6F5C0mao5"
  }
}
