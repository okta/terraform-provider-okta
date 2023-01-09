# resource "okta_app_ws_federation" "test" {
#   preconfigured_app = "aws_console"
#   label             = "testAcc_replace_with_uuid"
#   visibility        = false
# }


resource "okta_app_ws_federation" "test" {
		label    = "example"
		site_url = "https://signin.example.com/saml"
		# realm = "EXAMPLE"
		reply_url = "https://example-two.com"
		allow_override = false
		name_id_format = "uid"
		audience_restriction = "https://signin.example.com"
		authn_context_class_ref = "Kerberos"
		group_filter = "app1.*"
		# group_name = "username"
		group_value_format = "dn"
		username_attribute = "username"
		attribute_statements = "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/givenname|${user.firstName}|,http://schemas.xmlsoap.org/ws/2005/05/identity/claims/surname|${user.lastName}|"
		visibility = true  
}
