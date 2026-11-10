data "okta_authenticator" "tac" {
  key = "tac"
}

resource "okta_authenticator_tac" "test" {
  key  = "tac"
  name = "Test TAC-3"

  provider_tac {
    type = "TAC"

    configuration {
      length            = 16
      min_ttl           = 20
      default_ttl       = 30
      max_ttl           = 60
      multi_use_allowed = true

      complexity {
        letters            = true
        numbers            = true
        special_characters = true
      }
    }
  }
}
