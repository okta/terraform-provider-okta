resource "okta_app_oauth" "test1" {
    preconfigured_app = "strongdm"
    label = "StrongDM"
    type  = "web"
}

resource "okta_app_oauth" "test2" {
    preconfigured_app = "Applauz"
    label = "Applauz"
    type  = "web"
}

resource "okta_app_oauth" "test3" {
    preconfigured_app = "Deel"
    label = "Deel"
    type  = "web"
}

resource "okta_app_oauth" "test4" {
    label = "StrongDM_CUSTOM"
    type  = "web"
    redirect_uris  = ["http://redirect-uri-2.com/"]
}

resource "okta_app_oauth" "test5" {
  label          = "CustomApp001"
  type           = "web"
  grant_types    = ["implicit", "authorization_code"]
  redirect_uris  = ["http://redirect-uri.com/"]
  response_types = ["code", "token", "id_token"]
  issuer_mode    = "ORG_URL"
}
