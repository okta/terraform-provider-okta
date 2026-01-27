resource "okta_agent_pool_update" "example" {
  name          = "schedule_test"
  agent_type    = "AD"
  notify_admins = true
  pool_id       = "0oaspf3cfatE1nDO31d7"
  agents {
    id      = "a53slzqkptH2xEJ1r1d7"
    pool_id = "0oaspf3cfatE1nDO31d7" # this is required in schema
  }

  schedule {
    cron     = "0 3 * * WED"
    timezone = "Asia/Calcutta"
    delay    = 0
    duration = 1020
  }
}