resource "okta_campaign" "test" {
  name        = "Monthly access review of sales team"
  description = "Review access of all sales team members to a specific app"

  schedule_settings {
    type             = "ONE_OFF"
    start_date       = "2025-10-04T13:43:40.000Z"
    duration_in_days = 30
    time_zone        = "America/Los_Angeles"
  }

  resource_settings {
    type = "GROUP"

    target_resources {
      resource_id   = "00gnkw1sdqL30MdGk1d7"
      resource_type = "GROUP"
    }

    target_resources {
      resource_id   = "00go2ywj8vuFO2JzF1d7"
      resource_type = "GROUP"
    }
  }

  principal_scope_settings {
    type = "USERS"
  }

  reviewer_settings {
    type                      = "REVIEWER_EXPRESSION"
    reviewer_scope_expression = "user.profile.managerId"
    fallback_reviewer_id      = "00unkw1sfbTw08c0g1d7"
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
