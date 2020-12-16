resource "okta_template_email" "test" {
  type = "email.forgotPassword"

  translations {
    language = "en"
    subject  = "Stuff"
    template = "Hi $${user.firstName},<br/><br/>Blah blah $${resetPasswordLink}"
  }

  translations {
    language = "es"
    subject  = "Cosas"
    template = "Hola $${user.firstName},<br/><br/>Puedo ir al bano $${resetPasswordLink}"
  }
}
