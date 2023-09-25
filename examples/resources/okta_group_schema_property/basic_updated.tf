resource "okta_group_schema_property" "test" {
  index       = "testAcc_replace_with_uuid"
  title       = "terraform acceptance test updated 002"
  type        = "string"
  description = "terraform acceptance test updated 002"
  required    = true
  min_length  = 1
  max_length  = 70
  permissions = "READ_WRITE"
  master      = "OKTA"
  enum        = ["S", "M", "L", "XXL"]
  scope       = "NONE"

  one_of {
    const = "S"
    title = "Small"
  }

  one_of {
    const = "M"
    title = "Medium"
  }

  one_of {
    const = "L"
    title = "Large"
  }

  one_of {
    const = "XXL"
    title = "Extra Extra Large"
  }
}
