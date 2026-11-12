resource "okta_behavior_anomalous_ip" "test" {
  name   = "testAcc_replace_with_uuid"
  type   = "ANOMALOUS_IP"
  status = "ACTIVE"
}
