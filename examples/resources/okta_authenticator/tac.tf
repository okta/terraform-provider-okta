resource "okta_authenticator" "tac" {
  name   = "Temporary Access Code"
  key    = "tac"
  status = "ACTIVE"
  provider_json = jsonencode({
    type = "TAC"
    configuration = {
      length          = 8
      minTtl          = 10
      maxTtl          = 480
      defaultTtl      = 60
      multiUseAllowed = false
      complexity = {
        numbers           = true
        letters           = true
        specialCharacters = false
      }
    }
  })
}
