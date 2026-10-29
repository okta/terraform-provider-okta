resource "okta_group" "test" {
  # name = "あいうえおかきくけこさしすせそたちつてとなにぬねのはひふへほまみむめもらりるれろ"
  name = "testAcc_あいうえおかきくけこさしすせそたちつてとなにぬねのはひふへほまみ"
}

resource "okta_group_rule" "test" {
  # name              = "ABCDEFGHIJKLMNOPQRSTUVWXYZあいうえおかきくけこさしすせそたちつてとなにぬねのはひふへほまみむめもらりるれろ"
  name              = "testAcc_ABCDEFGHIJKLMNOPQRあいうえおかきくけこさしすせそたちつてとなにぬねのはひふへほまみむめもらりるれろ"
  status            = "ACTIVE"
  group_assignments = [okta_group.test.id]
  expression_type   = "urn:okta:expression:1.0"
  expression_value  = "String.startsWith(user.firstName,\"andy\")"
}
