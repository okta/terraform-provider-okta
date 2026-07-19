resource "okta_behavior_anomalous_location" "example" {
  type = "<type>"
  name = "Example Name"
  granularity = "<granularity>"

  # Optional fields
  # max_events_used_for_evaluation = 0
  # min_events_needed_for_evaluation = 0
  # radius_kilometers = 0
  # status = "ACTIVE"
}
