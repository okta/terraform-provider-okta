resource "okta_behavior" "test" {
  name = "testAcc_replace_with_uuid"
  type = "ANOMALOUS_LOCATION"
  number_of_authentications = 50
  location_granularity_type = "LAT_LONG"
  radius_from_location = 20
}
