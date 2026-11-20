resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "Smith"
  login      = "testAcc-replace_with_uuid@example.com"
  email      = "testAcc-replace_with_uuid@example.com"
}

resource "okta_campaign" "test" {
  name        = "Monthly access review of sales team"
  description = "Review access of all sales team members to a specific app"

  schedule_settings {
    type             = "ONE_OFF"
    start_date       = "2030-10-04T13:43:40.000Z"
    duration_in_days = 30
    time_zone        = "America/Los_Angeles"
  }

  resource_settings {
    type                                    = "APPLICATION"
    include_entitlements                    = false
    individually_assigned_apps_only         = false
    individually_assigned_groups_only       = false
    only_include_out_of_policy_entitlements = false
    target_resources {
      resource_id                          = "0oaws4am895IZbn6Q1d7"
      resource_type                        = "APPLICATION"
      include_all_entitlements_and_bundles = false
    }
  }

  principal_scope_settings {
    type = "USERS"
  }

  reviewer_settings {
    type                      = "REVIEWER_EXPRESSION"
    reviewer_scope_expression = "user.profile.managerId"
    fallback_reviewer_id      = okta_user.test.id
  }

  notification_settings {
    notify_reviewer_at_campaign_end           = true
    notify_reviewer_during_midpoint_of_review = false
    notify_reviewer_when_review_assigned      = false
    notify_reviewer_when_overdue              = false
    notify_review_period_end                  = false
  }

  remediation_settings {
    access_approved = "NO_ACTION"
    access_revoked  = "NO_ACTION"
    no_response     = "NO_ACTION"
  }
}


data "okta_campaign" "test" {
  id = okta_campaign.test.id
}
