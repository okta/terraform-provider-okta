resource "okta_principal_rate_limits" "test" {
  principal_id                   = "0oatyx4ukmqlnSQ0P1d7"
  principal_type                 = "OAUTH_CLIENT"
  default_percentage             = 55
  default_concurrency_percentage = 85
}