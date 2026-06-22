
# Create-then-read pattern (resource exists)
resource "okta_label" "test" {
  label_id = "replace_with_uuid"
  name = "test-name"
}

data "okta_label" "test" {
  id = okta_label.test.id
}
