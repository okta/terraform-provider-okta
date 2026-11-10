data "okta_authenticator" "tac" {
  key = "tac"
}

resource "okta_authenticator_tac" "test" {
  key  = "tac"
  name = "Temporary Access Code"

  provider_tac {
    type = "TAC"

    configuration {
      length            = 12
      min_ttl           = 15
      default_ttl       = 45
      max_ttl           = 120
      multi_use_allowed = false

      complexity {
        letters            = true
        numbers            = true
        special_characters = false
      }
    }
  }
}
