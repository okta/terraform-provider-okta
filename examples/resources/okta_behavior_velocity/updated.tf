resource "okta_behavior_velocity" "test" {
  name   = "testAcc_replace_with_uuid"
  type   = "VELOCITY"
  status = "INACTIVE"
  settings {
    velocity_kph = 805
  }
}
