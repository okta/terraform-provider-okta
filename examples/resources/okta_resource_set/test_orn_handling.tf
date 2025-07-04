variable "id" { type = string }

resource "okta_resource_set" "test" {
  label         = "testAcc_replace_with_uuid"
  description   = "testing, testing"
  resources_orn = [
    "orn:okta:directory:${var.id}:users",
    "orn:okta:directory:${var.id}:groups"
  ]
}
