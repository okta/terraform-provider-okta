resource "okta_behavior_anomalous_location" "test" {
  name   = "testAcc_replace_with_uuid"
  type   = "ANOMALOUS_LOCATION"
  status = "ACTIVE"
  settings {
    granularity = "CITY"
  }
}
