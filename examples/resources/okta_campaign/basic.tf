resource "okta_campaign" "test" {
  name        = "Monthly access review of sales team"
  description = "Review access of all sales team members to a specific app"

  campaign_type = "RESOURCE"

  schedule_settings {
    type            = "ONE_OFF"
    start_date      = "2025-08-07T00:00:00.000Z"
    duration_in_days = 30
    time_zone       = "America/Los_Angeles"
  }

  resource_settings {
    type = "APPLICATION"

    target_resources {
      resource_id = "0oanlpd3xkLkePi3W1d7"
      resource_type = "APPLICATION"
    }
    target_resources{
      resource_id = "0oao01ardu8r8qUP91d7"
      resource_type = "APPLICATION"
    }
  }

  principal_scope_settings {
    type = "USERS"

    predefined_inactive_users_scope {
      inactive_days = 90
    }
  }

  reviewer_settings {
    type        = "USER"
    reviewer_id = "00unkw1sfbTw08c0g1d7"
  }

  remediation_settings {
    access_approved = "NO_ACTION"
    access_revoked  = "NO_ACTION"
    no_response     = "NO_ACTION"
  }

  notification_settings {
    notify_reviewer_at_campaign_end = false
    notify_reviewer_during_midpoint_of_review = false
    notify_reviewer_when_review_assigned = false
    notify_reviewer_when_overdue = false
    notify_review_period_end = false
  }
}
