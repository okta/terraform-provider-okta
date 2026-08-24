resource "okta_template_sms" "test" {
  name     = "testAcc_replace_with_uuid"
  type     = "SMS_VERIFY_CODE"
  template = "$${org.name} code is: $${code}"
}
