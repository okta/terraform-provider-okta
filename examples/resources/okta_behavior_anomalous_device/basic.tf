resource "okta_behavior_anomalous_device" "test" {
  name   = "testAcc_replace_with_uuid"
  type   = "ANOMALOUS_DEVICE"
  status = "ACTIVE"
}
