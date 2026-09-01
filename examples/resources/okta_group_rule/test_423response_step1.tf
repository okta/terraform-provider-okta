locals {
  valid_expression_examples = [
    # "(user.firstName == \"TestAcc\" AND user.lastName == \"Smith\")",
    # "(user.firstName == \"TestAcc\" OR user.firstName == \"TestAcc2\")",
    # "user.firstName == \"(John)\"",
    "user.firstName == \"John\" AND (user.lastName == \"Doe\" OR user.lastName == \"Smith\")",
    "(user.department == \"Sales\" AND (user.city == \"SF\" OR user.city == \"NYC\")) OR (user.department == \"Engineering\" AND user.title == \"Senior\")",
    "String.startsWith(user.firstName,\"andy\")",
    "user.email == \"test+alias@example.com\" AND user.department == \"R&D\"",
    # "  user.firstName    ==    \"John\"   AND   user.lastName   ==   \"Doe\"  ",
    # "isMemberOfGroupNameStartsWith(\"IT-\")",
    # "user.login matches \"(a+)+$\"",
    # "hasWorkdayUser()", # hasWorkdayUser has no required args
  ]
}

resource "okta_group_rule" "rule" {
  count             = length(local.valid_expression_examples)
  status            = "ACTIVE"
  name              = "testAcc${count.index}_replace_with_uuid"
  expression_value  = local.valid_expression_examples[count.index]
  group_assignments = [okta_group.test.id]
}

resource "okta_group" "test" {
  name        = "testAcc_replace_with_uuid"
  description = "testing"
}
