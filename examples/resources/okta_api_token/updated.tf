resource "okta_api_token" example{
  id="00T1gtq5lsLg3q4dh1d7"
  name = "api-token-test-token"
  user_id="00unkw1sfbTw08c0g1d7"
  network{
    connection= "ZONE"
    exclude{
      ip= "nzonkw1sfwwiH9Hxt1d7"
    }
  }
  client_name="Okta API"
}