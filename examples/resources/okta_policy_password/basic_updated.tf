data "okta_group" "all" {
  name = "Everyone"
}

resource "okta_policy_password" "test" {
  name        = "testAcc_replace_with_uuid"
  type        = "PASSWORD"
  status      = "ACTIVE"
  description = "Terraform Acceptance Test Password Policy Updated"

  conditions {
    people {
      groups {
        include = [data.okta_group.all.id]
      }
    }
  }

  settings {
    password {
      age {
        max_age_days     = 60
        expire_warn_days = 15
        min_age_minutes  = 60
        history_count    = 0
      }

      complexity {
        min_length       = 12
        min_lower_case   = 0
        min_upper_case   = 0
        min_number       = 0
        min_symbol       = 1
        exclude_username = false
      }

      lockout {
        max_attempts                       = 10
        auto_unlock_minutes                = 2
        show_lockout_failures              = true
        user_lockout_notification_channels = ["EMAIL"]
      }

      breached_protection {
        expire_after_days = 6
        logout_enabled    = true
      }
    }
  }
}
