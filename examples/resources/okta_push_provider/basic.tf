resource "okta_push_provider" "example" {
  name          = "example"
  provider_type = "FCM"
  configuration {
    fcm_configuration {
      service_account_json {
        type                        = "service_account"
        project_id                  = "PROJECT_ID"
        private_key_id              = "KEY_ID"
        private_key                 = "-----BEGIN PRIVATE KEY-----REDACTED-----END PRIVATE KEY-----"
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