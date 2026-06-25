resource "okta_label" "test" {
  name = "test-label-replace_with_uuid"

  values {
    name = "test-value"
  }
}

data "okta_label" "test" {
  id = okta_label.test.id
}
