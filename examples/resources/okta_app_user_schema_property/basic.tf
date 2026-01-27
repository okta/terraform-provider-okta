resource "okta_app_oauth" "test" {
  label          = "testLabel"
  type           = "native"
  grant_types    = ["authorization_code"]
  redirect_uris  = ["http://d.com/"]
  response_types = ["code"]
}

resource "okta_app_user_schema_property" "test" {
  app_id      = okta_app_oauth.test.id
  index       = "testIndex"
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
}
