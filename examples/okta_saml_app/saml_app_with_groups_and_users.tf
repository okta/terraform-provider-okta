resource "okta_group" "group-%[1]d" {
  name = "testAcc_%[1]d"
}
resource "okta_user" "user-%[1]d" {
  admin_roles = ["APP_ADMIN", "USER_ADMIN"]
  first_name  = "TestAcc"
  last_name   = "blah"
  login       = "test-acc-%[1]d@testing.com"
  email       = "test-acc-%[1]d@testing.com"
  status      = "ACTIVE"
}

resource "okta_saml_app" "testAcc_%[1]d" {
  preconfigured_app = "amazon_aws"
  label             = "testAcc_%[1]d"
  users = [
    {
      id       = "${okta_user.user-%[1]d.id}"
      username = "${okta_user.user-%[1]d.email}"
    }
  ]
  groups = ["${okta_group.group-%[1]d.id}"]
  key = {
    years_valid = 3
  }
  app_settings_json = <<EOT
{
  "awsEnvironmentType":"aws.amazon",
  "groupFilter": "aws_(?{{accountid}}\\d+)_(?{{role}}[a-zA-Z0-9+=,.@\\-_]+)",
  "joinAllRoles": false,
  "loginURL": "https://console.aws.amazon.com/ec2/home",
  "roleValuePattern": "arn:aws:iam::$${accountid}:saml-provider/OKTA,arn:aws:iam::$${accountid}:role/$${role}",
  "sessionDuration": 3600
}
EOT
}
