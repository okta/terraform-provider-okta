resource "okta_group" "test" {
  name = "あいうえおかきくけこさしすせそたちつてとなにぬねのはひふへほまみむめもらりるれろ"
}

resource "okta_group_rule" "test" {
  name              = "ABCDEFGHIJKLMNOPQRSTUVWXYZあいうえおかきくけこさしすせそたちつてとなにぬねのはひふへほまみむめもらりるれろ"
  status            = "ACTIVE"
  group_assignments = [okta_group.test.id]
  expression_type   = "urn:okta:expression:1.0"
  expression_value  = "String.startsWith(user.firstName,\"andy\")"
}