resource "okta_group" "test" {
  name        = "testAcc_replace_with_uuid"
  description = "testing"
}

resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "Smith"
  login      = "testAcc-replace_with_uuid@example.com"
  email      = "testAcc-replace_with_uuid@example.com"
}

resource "okta_campaign" "test" {
  name        = "testAcc_replace_with_uuid"
  description = "Test campaign for datasource"

  schedule_settings {
    type             = "ONE_OFF"
    start_date       = "2027-10-04T13:43:40.000Z"
    duration_in_days = 30
    time_zone        = "America/Los_Angeles"
  }

  resource_settings {
    type = "GROUP"

    target_resources {
      resource_id   = okta_group.test.id
      resource_type = "GROUP"
    }
  }

  principal_scope_settings {
    type = "USERS"
  }

  reviewer_settings {
    type        = "USER"
    reviewer_id = okta_user.test.id
  }

  notification_settings {
    notify_reviewer_at_campaign_end           = false
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

data "okta_campaigns" "test" {
  depends_on = [okta_campaign.test]
}
