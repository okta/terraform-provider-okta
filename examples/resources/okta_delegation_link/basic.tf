resource "okta_delegation_link" "test" {
  from {
    type       = "OKTA_AUTHORIZATION_SERVER"
    client_orn = "orn:okta:directory:00otivnuq1X5554yc1d7:workload-principals:ai-agents:wlp114bmffd1byZg51d8"
    token_type = "ACCESS_TOKEN"
  }

  to {
    resource_orn = "orn:okta:directory:00otivnuq1X5554yc1d7:workload-principals:ai-agents:wlp11f2exa8aMDIS41d8"
  }
}
