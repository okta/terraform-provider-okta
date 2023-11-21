resource "okta_template_sms" "test" {
  type     = "SMS_VERIFY_CODE"
  template = "$${org.name} code is: $${code}"

  translations {
    language = "en"
    template = "$${org.name} code is: $${code}"
  }

  translations {
    language = "es"
    template = "$${org.name} es: $${code}."
  }
}
