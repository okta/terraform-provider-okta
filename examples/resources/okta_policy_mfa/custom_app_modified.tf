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
    { "enroll" : "NOT_ALLOWED", "id" : "aut123456789abcdef" },
    { "enroll" : "OPTIONAL", "id" : "aut123456789ghijkl" },
    { "enroll" : "OPTIONAL", "id" : "aut123456789mnopqr" }
  ]

  groups_included = [data.okta_group.all.id]
}
