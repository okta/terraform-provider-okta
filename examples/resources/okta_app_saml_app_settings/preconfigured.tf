resource "okta_app_saml" "test" {
  preconfigured_app = "amazon_aws"
  label             = "testAcc_replace_with_uuid"
  status            = "ACTIVE"
}

resource "okta_app_saml_app_settings" "test" {
  app_id = okta_app_saml.test.id
  settings = jsonencode(
    {
      "appFilter" : "okta",
      "awsEnvironmentType" : "aws.amazon",
      "groupFilter" : "aws_(?{{accountid}}\\\\d+)_(?{{role}}[a-zA-Z0-9+=,.@\\\\-_]+)",
      "joinAllRoles" : false,
      "loginURL" : "https://console.aws.amazon.com/ec2/home",
      "roleValuePattern" : "arn:aws:iam::$${accountid}:saml-provider/OKTA,arn:aws:iam::$${accountid}:role/$${role}",
      "sessionDuration" : 7600,
      "useGroupMapping" : false
    }
  )
}
