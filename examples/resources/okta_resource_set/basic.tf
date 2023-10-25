# For testing TF_VAR_hostname is set in provider_test.go .
# In a live environment the operator would export `TF_VAR_hostname=[the
# hostname]` in order to expose hostname as a variable below.
variable "hostname" {
  type = string
}

resource "okta_resource_set" "test" {
  label       = "testAcc_replace_with_uuid"
  description = "testing, testing"
  resources = [
    "https://${var.hostname}/api/v1/users",
    "https://${var.hostname}/api/v1/apps",
    "https://${var.hostname}/api/v1/groups"
  ]
}
