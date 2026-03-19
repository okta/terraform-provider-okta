resource "okta_email_smtp_server" "smtp_server_example" {
  port     = 587
  host     = "smtp.example.com"
  username = "abcd"
  password = "abcd"
  alias    = "server4"
}

data "okta_email_smtp_server" "test" {
  id = okta_email_smtp_server.smtp_server_example.id
}
