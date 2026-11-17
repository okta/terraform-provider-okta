resource "okta_behavior_velocity" "test" {
  name   = "testAcc_replace_with_uuid"
  type   = "VELOCITY"
  status = "ACTIVE"
  settings {
    velocity_kph = 805
  }
}

data "okta_behavior_velocity" "test" {
  id = okta_behavior_velocity.test.id
}
