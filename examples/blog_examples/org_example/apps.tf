data okta_group everyone {
  name = "Everyone"
}

// This is our SAML app that actually gives users access to Jamf.
resource "okta_saml_app" "jamf_sso" {
  label              = "Jamf SSO"
  preconfigured_app  = "jamfsoftwareserver"
  user_name_template = "$${source.login}"
  hide_ios           = true
  hide_web           = true

  app_settings_json = <<EOT
{
  "domain": "example.jamfcloud.com"
}
EOT

  groups = ["${data.okta_group.everyone.id}"]
}

resource "okta_bookmark_app" "jamf_admin" {
  label = "Jamf Admin Access"
  url = "https://example.jamfcloud.com"

  users {
    id = "${okta_user.bill_billson.id}"
    username = "${okta_user.bill_billson.login}"
  }
}

resource "okta_bookmark_app" "jamf_enrollment" {
  label = "Jamf Enrollment"
  url = "https://example.jamfcloud.com/enroll"
  groups = ["${data.okta_group.everyone.id}"]
}
