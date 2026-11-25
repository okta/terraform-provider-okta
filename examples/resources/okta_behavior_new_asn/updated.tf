resource "okta_behavior_new_asn" "test" {
  name   = "testAcc_replace_with_uuid_updated"
  type   = "NEW_ASN"
  status = "INACTIVE"
  settings {
    max_events_used_for_evaluation   = 50
    min_events_needed_for_evaluation = 5
  }
}
