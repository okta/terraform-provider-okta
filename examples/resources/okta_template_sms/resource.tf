resource "okta_template_sms" "example" {
  name     = "My SMS Template"
  type     = "SMS_VERIFY_CODE"
  template = "Your $${org.name} code is: $${code}"
  translations = {
    "en" = "Your $${org.name} code is: $${code}"
    "es" = "Tu código de $${org.name} es: $${code}."
  }
}
