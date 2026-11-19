resource "okta_behavior_anomalous_location" "test" {
  name   = "testAcc_replace_with_uuid"
  type   = "ANOMALOUS_LOCATION"
  status = "ACTIVE"
  settings {
    granularity                      = "CITY"
    max_events_used_for_evaluation   = 20
    min_events_needed_for_evaluation = 0
  }
}
