resource "okta_group" "test" {
  name = "[xx]ZZZ_ああああ_yyyyyyあいうえおかきくけ1ww1"
}

resource "okta_group_rule" "test" {
  name              = "[xx]ZZZ_ああああ_yyyyyyあいうw1w1えおかきくけ1"
  status            = "ACTIVE"
  group_assignments = [okta_group.test.id]
  expression_type   = "urn:okta:expression:1.0"
  expression_value  = "String.startsWith(user.firstName,\"andy\")"
}