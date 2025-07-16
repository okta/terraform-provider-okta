data "okta_features" "test" {
  substring = "MFA"
}

output "first_feature_name" {
  value = try(data.okta_features.test.features[0].name, "")
}