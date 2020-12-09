resource "okta_user_schema" "test" {
  index       = "testAcc_replace_with_uuid"
  title       = "terraform acceptance test"
  type        = "array"
  array_type  = "string"
  description = "testing"
  required    = false
  master      = "OKTA"
  array_enum  = ["test", "1", "2"]

  array_one_of {
    const = "test"
    title = "test"
  }

  array_one_of {
    const = "1"
    title = "1"
  }

  array_one_of {
    const = "2"
    title = "2"
  }
}
