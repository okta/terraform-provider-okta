resource "okta_entitlement" "test" {
  name           = "Entitlement Bundle"
  external_value = "Entitlement Bundle"
  description    = "Some license entitlement"
  multi_value    = true
  data_type      = "array"

  parent {
    external_id = "0oao01ardu8r8qUP91d7"
    type        = "APPLICATION"
  }

  values {
    name           = "value1"
    description    = "description for value1"
    external_value = "value_1"
  }

  values {
    name           = "value2"
    description    = "description for value2"
    external_value = "value_2"
  }
}

data "okta_entitlement" "test" {
  id = okta_entitlement.test.id
}