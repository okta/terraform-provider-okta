resource "okta_email_smtp" "test" {
  alias= "server4"
  host="Issue1950-4"
  port= 8086
  username= "test_user"
  enabled= false
  password = "adot4a"
}