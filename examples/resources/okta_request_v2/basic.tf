resource "okta_request_v2" "test" {
  requested {
    type     = "CATALOG_ENTRY"
    entry_id = "cen1048kevzlpJ0cf1d7"
  }
  requested_for {
    type        = "OKTA_USER"
    external_id = "00unkw1sfbTw08c0g1d7"
  }
}
