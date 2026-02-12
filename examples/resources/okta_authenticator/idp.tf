resource "okta_authenticator" "extenral_idp" {
  name = "external idp"
  key  = "external_idp"
  provider_json = jsonencode({
    "type" : "CLAIMS",
    "configuration" : {
      "idpId" : "0oauo0g9snGE2oZcx1d7"
    }
  })
  status = "ACTIVE"
}