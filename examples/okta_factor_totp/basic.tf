resource "okta_factor_totp" "test" {
  name = "testAcc_replace_with_uuid"
  otp_length = 10
  hmac_algorithm = "HMacSHA256"
  time_step = 30
  clock_drift_interval = 10
  shared_secret_encoding = "hexadecimal"
}
