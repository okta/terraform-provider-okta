resource "okta_app_ws_federation" "exampleWsFedApp" {
	label    = "exampleWsFedApp"
	site_url = "https://signin.test.com/saml"
	reply_url = "https://test.com"
	reply_override = false
	name_id_format = "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified"
	audience_restriction = "https://signin.test.com"
	authn_context_class_ref = "urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport"
	group_filter = "app1.*"
	group_name = "username"
	group_value_format = "dn"
	username_attribute = "username"
	attribute_statements = "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/givenname|bob|,http://schemas.xmlsoap.org/ws/2005/05/identity/claims/surname|hope|"
	visibility = false
	status = "active"
}

data "okta_app_ws_federation" "exampleWsFedApp" {
  label = "example"
}
