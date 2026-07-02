resource "okta_request_condition" "test" {
  resource_id          = "0oaty79sxfpt7AVvI1d7"
  approval_sequence_id = "695cd5cd4bfe7d01dbcacb51"
  name                 = "issue-2780"

  requester_settings {
    type = "EVERYONE"
  }

  access_scope_settings {
    type = "RESOURCE_DEFAULT"
  }
}
