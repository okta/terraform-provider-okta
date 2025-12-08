data "okta_group" "all" {
  name = "Everyone"
}

resource "okta_policy_mfa" "test" {
  name        = "testAcc_replace_with_uuid"
  status      = "ACTIVE"
  description = "Terraform Acceptance Test MFA Policy with Specific Custom Apps"
  is_oie      = true

  okta_password = {
    enroll = "REQUIRED"
  }

  security_question = {
    enroll = "REQUIRED"
  }

  custom_app = [
    { "enroll" : "OPTIONAL", "id" : "aut1234567890abcdefg" },
    { "enroll" : "OPTIONAL", "id" : "aut1234567890hijklmn" }
  ]

  groups_included = [data.okta_group.all.id]
}
