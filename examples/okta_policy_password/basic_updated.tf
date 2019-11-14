data okta_group all {
  name = "Everyone"
}

resource okta_policy_password test {
  name                           = "testAcc_replace_with_uuid"
  status                         = "INACTIVE"
  description                    = "Terraform Acceptance Test Password Policy Updated"
  password_min_length            = 12
  password_min_lowercase         = 0
  password_min_uppercase         = 0
  password_min_number            = 0
  password_min_symbol            = 1
  password_exclude_username      = false
  password_exclude_first_name    = true
  password_exclude_last_name     = true
  password_max_age_days          = 60
  password_expire_warn_days      = 15
  password_min_age_minutes       = 60
  password_history_count         = 5
  password_max_lockout_attempts  = 0
  password_auto_unlock_minutes   = 2
  password_show_lockout_failures = true
  question_min_length            = 10
  recovery_email_token           = 20160
  sms_recovery                   = "ACTIVE"

  groups_included = ["${data.okta_group.all.id}"]
}
