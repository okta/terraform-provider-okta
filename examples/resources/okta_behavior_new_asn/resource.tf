resource "okta_behavior_new_asn" "example" {
  name   = "My New ASN Behavior"
  type   = "NEW_ASN"
  status = "ACTIVE"
  settings {
    max_events_used_for_evaluation   = 20
    min_events_needed_for_evaluation = 0
  }
}
