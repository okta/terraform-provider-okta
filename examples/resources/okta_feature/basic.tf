data "okta_features" "example" {
  label = "Android Device Trust"
}

resource "okta_feature" "test" {
  feature_id = tolist(data.okta_features.example.features)[0].id
}
