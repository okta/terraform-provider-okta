variable "hostname" {
  type = string
}

resource "okta_resource_set" "test" {
  label       = "testAcc_replace_with_uuid"
  description = "testing, testing updated"
  resources = [
    "https://${var.hostname}/api/v1/users",
    "https://${var.hostname}/api/v1/apps",
  ]
}
