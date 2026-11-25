resource "okta_behavior_new_asn" "test" {
  name   = "testAcc_replace_with_uuid"
  type   = "NEW_ASN"
  status = "ACTIVE"
  settings {
    max_events_used_for_evaluation   = 20
    min_events_needed_for_evaluation = 0
  }
}
