resource "okta_app_features" "example" {
  app_id = "0oarblaf7hWdLawNg1d7"
  name   = "USER_PROVISIONING"
  capabilities {
    create {
      lifecycle_create {
        status = "DISABLED"
      }
    }
    update {
      lifecycle_deactivate {
        status = "ENABLED"
      }
      profile {
        status = "DISABLED"
      }
      password {
        status = "ENABLED"
        seed   = "RANDOM"
        change = "KEEP_EXISTING"
      }
    }
  }
  status = "ENABLED"
}