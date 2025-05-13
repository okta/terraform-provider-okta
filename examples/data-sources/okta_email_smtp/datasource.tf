resource "okta_email_smtp" "smtp_server2" {
  port     = 587
  host     = "smtp.example.com"
  username = "abcd"
  password = "abcd"
  alias    = "server4"
}

data "okta_email_smtp" "test" {
  id = okta_email_smtp.smtp_server2.id
}
