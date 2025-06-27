variable "hostname" { type = string }

resource "okta_group" "test" {
  count = 101
  name  = "testAcc_${count.index}"
}

resource "okta_resource_set" "test" {
  label       = "testAcc_replace_with_uuid"
  description = "testing pagination with large resource set"
  resources   = [ for g in okta_group.test[*] : "https://${var.hostname}/api/v1/groups/${g.id}" ]
}
