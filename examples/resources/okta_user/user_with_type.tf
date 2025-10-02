# Example of okta_user with user type

# First create a user type
resource "okta_user_type" "example" {
  name         = "contractor"
  display_name = "Contractor"
  description  = "A contractor user type"
}

# Create a user with the user type
resource "okta_user" "example" {
  first_name = "John"
  last_name  = "Doe"
  login      = "john.doe@example.com"
  email      = "john.doe@example.com"

  type {
    id = okta_user_type.example.id
  }
}

