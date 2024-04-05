resource "okta_group" "example" {
  name        = "Example"
  description = "My Example Group"
}

# Custom profile attributes
resource "okta_group" "example" {
  name        = "Example"
  description = "My Example Group"
  custom_profile_attributes = jsonencode({
    "example1" = "testing1234",
    "example2" = true,
    "example3" = 54321
  })
}
