resource "okta_email_smtp" "test" {
  alias= "CustomisedServer"
  host="192.168.2.0"
  port= 8086
  username= "test_user"
  enabled= false
  password = "testPwd"
}