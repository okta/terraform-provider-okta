resource "okta_delegation_link" "test" {
  from {
    type       = "OKTA_AUTHORIZATION_SERVER"
    client_orn = "orn:okta:idp:00op58x1t9iNay3Ojao5:apps:oidc_client:0oatkyugHHHcLBCFTao5"
    token_type = "ACCESS_TOKEN"
  }

  to {
    resource_orn = "orn:okta:directory:00op58x1t9iNay3Ojao5:workload-principals:ai-agents:wlpt1ulvUBdpeNxDBao5"
  }
}
data "okta_delegation_link" "test" {
  id = okta_delegation_link.test.id
}
