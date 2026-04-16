resource "okta_user" "included-replace_with_uuid" {
	first_name = "TestAcc"
	last_name  = "Included"
	login      = "testAcc-included-replace_with_uuid@example.com"
	email      = "testAcc-included-replace_with_uuid@example.com"
}

resource "okta_user" "excluded-replace_with_uuid" {
	first_name = "TestAcc"
	last_name  = "Excluded"
	login      = "testAcc-excluded-replace_with_uuid@example.com"
	email      = "testAcc-excluded-replace_with_uuid@example.com"
}

resource "okta_group" "included-replace_with_uuid" {
	name = "testAcc_included_replace_with_uuid"
}

resource "okta_group" "excluded-replace_with_uuid" {
	name = "testAcc_excluded_replace_with_uuid"
}

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

	users_included  = [okta_user.included-replace_with_uuid.id]
	users_excluded  = [okta_user.excluded-replace_with_uuid.id]
	groups_included = [okta_group.included-replace_with_uuid.id]
	groups_excluded = [okta_group.excluded-replace_with_uuid.id]

	password_reset_access_control = "AUTH_POLICY"
}
