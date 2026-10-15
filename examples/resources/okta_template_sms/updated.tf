resource "okta_template_sms" "test" {
  name     = "testAcc_replace_with_uuid"
  type     = "SMS_VERIFY_CODE"
  template = "$${org.name} updated code is: $${code}"
  translations = {
    "en" = "$${org.name} updated code is: $${code}"
    "es" = "$${org.name} es: $${code}."
    "fr" = "$${org.name} est: $${code}."
  }
}
