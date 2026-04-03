data "okta_default_policy" "default-replace_with_uuid" {
	type = "PASSWORD"
}

resource "okta_policy_rule_password" "testAcc_replace_with_uuid" {
	policy_id = data.okta_default_policy.default-replace_with_uuid.id
	name      = "testAcc_replace_with_uuid"
	status    = "ACTIVE"

	password_change = "ALLOW"
	password_reset  = "ALLOW"
	password_unlock = "ALLOW"

	password_reset_access_control = "LEGACY"

	password_reset_requirement {
		method_constraints {
			method                 = "otp"
			allowed_authenticators = ["google_otp"] # must be passed in case method is otp
		}
		primary_methods = ["otp"]
		step_up_enabled = true
	}
}
