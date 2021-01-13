resource "okta_template_sms" "test" {
  type     = "SMS_VERIFY_CODE"
  template = "Your $${org.name} code is: $${code}"

  translations {
    language = "en"
    template = "Your $${org.name} code is: $${code}"
  }

  translations {
    language = "es"
    template = "Tu código de $${org.name} es: $${code}."
  }
}
