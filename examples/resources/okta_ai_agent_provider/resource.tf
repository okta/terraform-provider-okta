resource "okta_ai_agent_provider" "test" {
  source_orn = "orn:okta:idp:00op58x1t9iNay3Ojao5:apps:office365:0oat1a3yIbqyw7VY0ao5"

  configuration = {
    clientId     = "1b2af17d-98de-408c-a4d3"
    clientSecret = "<REDACTED>"
  }
}
