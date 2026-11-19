resource "okta_behavior_anomalous_device" "test" {
  name   = "testAcc_replace_with_uuid"
  type   = "ANOMALOUS_DEVICE"
  status = "ACTIVE"
  settings {
    max_events_used_for_evaluation   = 20
    min_events_needed_for_evaluation = 0
  }
}