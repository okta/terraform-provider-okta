resource "okta_behavior" "my_location" {
  name                      = "My Location"
  type                      = "ANOMALOUS_LOCATION"
  number_of_authentications = 50
  location_granularity_type = "LAT_LONG"
  radius_from_location      = 20
}

resource "okta_behavior" "my_city" {
  name                      = "My City"
  type                      = "ANOMALOUS_LOCATION"
  number_of_authentications = 50
  location_granularity_type = "CITY"
}

resource "okta_behavior" "my_device" {
  name                      = "My Device"
  type                      = "ANOMALOUS_DEVICE"
  number_of_authentications = 50
}

resource "okta_behavior" "my_ip" {
  name                      = "My IP"
  type                      = "ANOMALOUS_IP"
  number_of_authentications = 50
}

resource "okta_behavior" "my_velocity" {
  name     = "My Velocity"
  type     = "VELOCITY"
  velocity = 25
}
