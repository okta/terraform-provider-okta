resource "okta_api_token" "example" {
  id      = "00T1gtr35t8ZbfBfV1d7"
  name    = "api-token-test-token"
  user_id = "00unkw1sfbTw08c0g1d7"
  network {
    connection = "ANYWHERE"
  }
  client_name = "Okta API"
}