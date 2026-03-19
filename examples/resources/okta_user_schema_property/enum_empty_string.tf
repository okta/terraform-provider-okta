resource "okta_user_schema_property" "test" {
  index       = "testAcc_replace_with_empty_enum"
  title       = "terraform acceptance test"
  type        = "string"
  description = "testing"
  master      = "PROFILE_MASTER"
  enum        = ["", "one", "two", "three"]
  one_of {
    const = ""
    title = "(none)"
  }
  one_of {
    const = "one"
    title = "string One"
  }
  one_of {
    const = "two"
    title = "string Two"
  }
  one_of {
    const = "three"
    title = "string Three"
  }
}
