# Create a resource set to manage collections of Okta resources
resource "okta_resource_set" "example" {
  label       = "My Resource Set"
  description = "Resource set for managing API endpoints"
  
  # Use URL-based resources
  resources = [
    "https://your-domain.okta.com/api/v1/users",
    "https://your-domain.okta.com/api/v1/groups"
  ]
  
  # Alternatively, use ORN-based resources
  # resources_orn = [
  #   "orn:okta:directory:123:users",
  #   "orn:okta:directory:123:groups"
  # ]
}
