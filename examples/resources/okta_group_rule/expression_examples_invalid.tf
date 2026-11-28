# Invalid Expression Examples for use in future acceptance tests

resource "okta_group_rule" "invalid_expression_examples" {
  for_each = toset([
    "user.firstName == \"John\")",         # Invalid: unbalanced parenthesis
    "user.firstName == \"John\"))",        # Invalid: unbalanced parenthesis
    "(user.firstName == \"John\"",         # Invalid: unbalanced parenthesis
    "user.firstName ==",                   # Invalid: invalid trailing "==" operator
    "user.firstName == \"John\" OR",       # Invalid: invalid trailing "OR" operator
    "user.firstName == \"John\" AND",      # Invalid: invalid trailing "AND" operator
    "user.firstName ==> \"John\"",         # Invalid: invalid "==>" operator
    "profile.firstName == \"John\"",       # Invalid: group rule attr ref "profile." vs "user.
    "NonExistentFunction(user.firstName)", # Invalid: invalid function name
    "String.to lowercase(user.email)",     # Invalid: function space in name
    "String.1toLowerCase(user.email)",     # Invalid: function name starts with number
    "String.to$lower(user.email)",         # Invalid: special characters in function name
    "string.tolowercase(user.email)",      # Invalid: incorrect case for function name
    "toUpperCase(user.firstName)",         # Invalid: does not resolve to type Boolean
    "",                                    # Invalid: Expression cannot be empty

    # Function signature errors
    "isMemberOfGroup(\"00gxxxxxxxxxxxxxxxx\")",     # Invalid group ID
    "String.toLowerCase()",                         # Missing required arguments
    "Arrays.contains(\"not_an_array\", \"value\")", # Type mismatch
  ])
  name = "invalid_expression_${each.index}"
  expression_value = each.value
  group_assignments = [okta_group.test.id]
}

resource "okta_group" "test" {
  name = "testAcc_replace_with_uuid"
  description = "testing"
}
