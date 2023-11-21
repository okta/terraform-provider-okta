resource "okta_app_bookmark" "test" {
  label = "testAcc_replace_with_uuid"
  url   = "https://test.com"
}

resource "okta_policy_profile_enrollment" "test" {
  name = "testAcc_replace_with_uuid"
}

resource "okta_policy_profile_enrollment" "test_2" {
  name = "testAcc_replace_with_uuid_2"
}

resource "okta_policy_profile_enrollment_apps" "test" {
  policy_id = okta_policy_profile_enrollment.test.id
}

resource "okta_policy_profile_enrollment_apps" "test_2" {
  policy_id  = okta_policy_profile_enrollment.test_2.id
  apps       = [okta_app_bookmark.test.id]
  depends_on = [okta_policy_profile_enrollment_apps.test]
}
