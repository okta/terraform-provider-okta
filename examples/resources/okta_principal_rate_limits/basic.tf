resource "okta_principal_rate_limits" "test" {
  principal_id                   = "00T1cuuokxSImd8ST1d7"
  principal_type                 = "SSWS_TOKEN"
  default_percentage             = 49
  default_concurrency_percentage = 75
}