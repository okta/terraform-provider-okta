resource "okta_user_schema" "testAcc_replace_with_uuid" {
  index       = "testAcc_replace_with_uuid"
  title       = "terraform acceptance test"
  type        = "string"
  description = "terraform acceptance test"
  required    = false
  min_length  = 1
  max_length  = 50
  permissions = "READ_ONLY"
  master      = "PROFILE_MASTER"
  enum        = ["S", "M", "L", "XL"]

  one_of = [
    {
      const = "S"
      title = "Small"
    },
    {
      const = "M"
      title = "Medium"
    },
    {
      const = "L"
      title = "Large"
    },
    {
      const = "XL"
      title = "Extra Large"
    },
  ]
}
