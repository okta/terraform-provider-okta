resource "okta_user_schema_property" "test" {
  index       = "testAcc_replace_with_uuid"
  title       = "terraform acceptance test"
  type        = "array"
  description = "testing"
  master      = "OKTA"
  scope       = "SELF"
  array_type  = "number"
  array_enum  = [0.01, 0.02, 0.03]
  array_one_of {
    title = "1"
    const = 0.01
  }
  array_one_of {
    title = "2"
    const = 0.02
  }
  array_one_of {
    title = "3"
    const = 0.03
  }
}
