resource "okta_user_schema" "testAcc_%[1]d" {
  index       = "testAcc%[1]d"
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
