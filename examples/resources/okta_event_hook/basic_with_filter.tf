resource "okta_event_hook" "test" {
  name = "testAcc_replace_with_uuid"
  events = [
    "group.user_membership.add",
  ]

  filter {
    event     = "group.user_membership.add"
    condition = "event.target.?[type eq 'UserGroup'].size()>0 && event.target.?[displayName eq 'Sales'].size()>0"
  }

  channel = {
    type    = "HTTP"
    version = "1.0.0"
    uri     = "https://example.com/test"
  }

  auth = {
    type  = "HEADER"
    key   = "Authorization"
    value = "123"
  }
}
