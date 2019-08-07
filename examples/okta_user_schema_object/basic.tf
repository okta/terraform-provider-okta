resource okta_user_schema_object test {
  custom = [
    {
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
        const = "XL"
        title = "Extra Large"
      }
    },
    {
      index       = "testAccBasic_replace_with_uuid"
      title       = "terraform acceptance test"
      type        = "string"
      description = "terraform acceptance test"
    },
  ]
}
