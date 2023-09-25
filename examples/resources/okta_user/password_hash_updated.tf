resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "Smith"
  login      = "testAcc-replace_with_uuid@example.com"
  email      = "testAcc-replace_with_uuid@example.com"
  status     = "STAGED"
  password_hash {
    algorithm   = "BCRYPT"
    work_factor = 10
    salt        = "rwh3vH166HCH/NT9XV5FYu"
    value       = "qaMqvAPULkbiQzkTCWo5XDcvzpk8Tna"
  }
}
