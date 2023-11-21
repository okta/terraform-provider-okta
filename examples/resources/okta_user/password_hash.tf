resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "Smith"
  login      = "testAcc-replace_with_uuid@example.com"
  email      = "testAcc-replace_with_uuid@example.com"
  status     = "STAGED"
  password_hash {
    algorithm  = "SHA-512"
    salt       = "TXlTYWx0"
    salt_order = "PREFIX"
    value      = "QrozP8a+KfoHu6mPFysxLoO5LMQsd2Fw6IclZUf8xQjetJOCGS93vm68h+VaFX0LHSiF/GxQkykq1vofmx6NGA=="
  }
}
