resource "okta_app_oauth" "test1" {
  preconfigured_app = "strongdm"
  label             = "StrongDM_Updated"
  type              = "web"
  redirect_uris     = ["https://strongdm.example.com/callback"]
}

resource "okta_app_oauth" "test2" {
  preconfigured_app = "Applauz"
  label             = "Applauz_Updated"
  type              = "web"
  redirect_uris     = ["https://applauz.example.com/callback"]
}

resource "okta_app_oauth" "test3" {
  preconfigured_app = "Deel"
  label             = "Deel_Updated"
  type              = "web"
  redirect_uris     = ["https://deel.example.com/callback"]
}

resource "okta_app_oauth" "test4" {
  label         = "StrongDM_CUSTOM_Updated"
  type          = "web"
  redirect_uris = ["http://redirect-uri-2-updated.com/"]
}

resource "okta_app_oauth" "test5" {
  label          = "CustomApp001_Updated"
  type           = "web"
  grant_types    = ["implicit", "authorization_code"]
  redirect_uris  = ["http://redirect-uri-updated.com/"]
  response_types = ["code", "token", "id_token"]
  issuer_mode    = "ORG_URL"
}
