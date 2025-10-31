resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "Smith"
  login      = "testAcc-replace_with_uuid@example.com"
  email      = "testAcc-replace_with_uuid@example.com"
}

resource "okta_campaign" "test" {
  name          = "Monthly access review of sales team"
  description   = "Multi app campaign"
  campaign_type = "RESOURCE"

  schedule_settings {
    type             = "ONE_OFF"
    start_date       = "2025-10-04T13:43:40.000Z"
    duration_in_days = 21
    time_zone        = "America/Vancouver"
  }

  resource_settings {
    type                                    = "APPLICATION"
    include_entitlements                    = true
    individually_assigned_apps_only         = false
    individually_assigned_groups_only       = false
    only_include_out_of_policy_entitlements = false
    target_resources {
      resource_id                          = "0oao01ardu8r8qUP91d7"
      resource_type                        = "APPLICATION"
      include_all_entitlements_and_bundles = true
    }
    target_resources {
      resource_id                          = "0oanlpd3xkLkePi3W1d7"
      resource_type                        = "APPLICATION"
      include_all_entitlements_and_bundles = false
    }
  }

  principal_scope_settings {
    type                      = "USERS"
    include_only_active_users = false
  }

  reviewer_settings {
    type                   = "USER"
    reviewer_id            = okta_user.test.id
    self_review_disabled   = true
    justification_required = true
    bulk_decision_disabled = true
  }

  notification_settings {
    notify_reviewer_when_review_assigned      = false
    notify_reviewer_at_campaign_end           = false
    notify_reviewer_when_overdue              = false
    notify_reviewer_during_midpoint_of_review = false
    notify_review_period_end                  = false
  }

  remediation_settings {
    access_approved = "NO_ACTION"
    access_revoked  = "NO_ACTION"
    no_response     = "NO_ACTION"
  }
}