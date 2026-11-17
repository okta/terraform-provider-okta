resource "okta_behavior_anomalous_ip" "test" {
  name   = "testAcc_replace_with_uuid"
  type   = "ANOMALOUS_IP"
  status = "ACTIVE"
  settings {
    max_events_used_for_evaluation   = 50
    min_events_needed_for_evaluation = 0
  }
}

data "okta_behavior_anomalous_ip" "test" {
  id = okta_behavior_anomalous_ip.test.id
}
