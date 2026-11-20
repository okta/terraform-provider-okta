resource "okta_delegation_link" "test" {
  from {
    type       = "OKTA_AUTHORIZATION_SERVER"
    client_orn = "orn:okta:idp:00otxk5i5T0Rk59Z3ao5:apps:office365:0oautxdhNWOdr5J9qao5"
    token_type = "ACCESS_TOKEN"
  }

  to {
    resource_orn = "orn:okta:directory:00otxk5i5T0Rk59Z3ao5:workload-principals:mcp:wlputxd97zgaV8w5Rao5"
  }
}
data "okta_delegation_link" "test" {
  id = okta_delegation_link.test.id
}
