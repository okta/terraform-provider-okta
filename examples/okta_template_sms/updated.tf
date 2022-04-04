resource "okta_template_sms" "test" {
  type     = "SMS_VERIFY_CODE"
  template = "Your $${org.name} updated code is: $${code}"

  translations {
    language = "en"
    template = "Your $${org.name} updated code is: $${code}"
  }

  translations {
    language = "es"
    template = "Tu código actualizado de $${org.name} es: $${code}."
  }

  translations {
    language = "fr"
    template = "Votre code mis à jour $${org.name} est: $${code}."
  }
}
