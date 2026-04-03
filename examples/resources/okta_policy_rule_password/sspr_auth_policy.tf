data "okta_default_policy" "default-replace_with_uuid" {
	type = "PASSWORD"
}

resource "okta_policy_rule_password" "testAcc_replace_with_uuid" {
	policy_id = data.okta_default_policy.default-replace_with_uuid.id
	name      = "testAcc_replace_with_uuid"
	status    = "ACTIVE"

	password_change = "ALLOW"
	password_reset  = "ALLOW"
	password_unlock = "DENY"

	password_reset_access_control = "AUTH_POLICY"
}
