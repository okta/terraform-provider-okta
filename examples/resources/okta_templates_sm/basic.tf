resource "okta_templates_sm" "test" {
  name     = "testAcc_replace_with_uuid"
  template = "$${org.name} code is: $${code}"
  type     = "SMS_VERIFY_CODE"
}
