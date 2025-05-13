# data "okta_email_smtp" "test" {
# }

# data "okta_email_smtp" "forgot_password" {
#   brand_id      = tolist(data.okta_brands.test.brands)[0].id
#   template_name = "ForgotPassword"
# }

resource "okta_email_smtp" "test" {
  alias= "server4"
  host="Issue1950-4"
  port= 8086
  username= "test_user"
  enabled= false
  password = "adot4a"
}
