resource "okta_policy_device_assurance_ios" "example" {
  name      = "My iOS Policy"
  platform  = "IOS"
  jailbreak = false
}

data "okta_policy_device_assurance_ios" "example" {
  id = okta_policy_device_assurance_ios.example.id
}
