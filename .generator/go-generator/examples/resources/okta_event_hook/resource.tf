resource "okta_event_hook" "example" {
  uri = "https://example.com"
  type = "<type>"
  version = "<version>"
  items = "<items>"
  type = "<type>"
  name = "Example Name"

  # Optional fields
  # key = "<key>"
  # type = "<type>"
  # value = "<value>"
  # key = "<key>"
  # value = "<value>"
}
