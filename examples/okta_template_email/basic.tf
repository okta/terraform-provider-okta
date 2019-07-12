resource okta_template_email test {
  type = "email.forgotPassword"

  translations {
    language = "en"
    subject  = "Stuff"
    template = "Hi $${user.firstName}"
  }

  translations {
    language = "es"
    subject  = "Cosas"
    template = "Hola $${user.firstName}"
  }
}
