resource "okta_push_provider" "example" {
  name          = "example"
  provider_type = "FCM"
  configuration {
    fcm_configuration {
      service_account_json {
        type                        = "service_account"
        project_id                  = "PROJECT_ID"
        private_key_id              = "KEY_ID"
        private_key                 = "-----BEGIN PRIVATE KEY-----MIGTAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBHkw\ndwIBAQQgTrVNo3aKMlnSuvhhLttykBgUj80HY+/RXqiWjXlk5M+gCgYIKoZIzj0DAQehRANCAASTa2fP8SDHMiaQ7e+80LZxxBVxYmDb8w5e\nIpHR0GZoSUgFASqU3L7VDqLs675+IbxRHXX0/lMeVKyA2oxSBpZk-----END PRIVATE KEY-----"
        client_email                = "SERVICE_ACCOUNT_EMAIL"
        client_id                   = "CLIENT_ID"
        auth_uri                    = "https://accounts.google.com/o/oauth2/auth"
        token_uri                   = "https://oauth2.googleapis.com/token"
        auth_provider_x509_cert_url = "https://www.googleapis.com/oauth2/v1/certs"
        client_x509_cert_url        = "https://www.googleapis.com/robot/v1/metadata/x509/SERVICE_ACCOUNT_EMAIL"
      }
    }
  }
}

data "okta_push_provider" "test" {
  id = okta_push_provider.example.id
}