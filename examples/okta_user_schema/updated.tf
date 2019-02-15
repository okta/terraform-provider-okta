resource "okta_user_schema" "testAcc_%[1]d" {
  index       = "testAcc_%[1]d"
  title       = "terraform acceptance test updated"
  type        = "string"
  description = "terraform acceptance test updated"
  required    = true
  min_length  = 1
  max_length  = 70
  permissions = "READ_WRITE"
  master      = "OKTA"
  enum        = ["S", "M", "L", "XXL"]

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
      const = "XXL"
      title = "Extra Extra Large"
    },
  ]
}
