resource "okta_group" "group" {
  name = "testAcc_replace_with_uuid"
}

resource "okta_group" "group1" {
  name = "testAcc_replace_with_uuid_1"
}

resource "okta_group" "group2" {
  name = "testAcc_replace_with_uuid_2"
}

resource "okta_user" "user" {
  admin_roles = ["APP_ADMIN", "USER_ADMIN"]
  first_name  = "TestAcc"
  last_name   = "blah"
  login       = "testAcc-replace_with_uuid@example.com"
  email       = "testAcc-replace_with_uuid@example.com"
  status      = "ACTIVE"
}

resource "okta_user" "user1" {
  first_name = "TestAcc1"
  last_name  = "blah"
  login      = "testAcc-1-replace_with_uuid@example.com"
  email      = "testAcc-1-replace_with_uuid@example.com"
  status     = "ACTIVE"
}

resource "okta_app_saml" "test" {
  preconfigured_app = "amazon_aws"
  label             = "testAcc_replace_with_uuid"

  users {
    id       = okta_user.user.id
    username = okta_user.user.email
  }

  users {
    id       = okta_user.user1.id
    username = okta_user.user1.email
  }

  groups = [okta_group.group.id, okta_group.group1.id, okta_group.group2.id]

  key_years_valid = 3
  key_name        = "hello"

  app_settings_json = <<EOT
{
  "appFilter":"okta",
  "awsEnvironmentType":"aws.amazon",
  "groupFilter": "aws_(?{{accountid}}\\d+)_(?{{role}}[a-zA-Z0-9+=,.@\\-_]+)",
  "joinAllRoles": false,
  "loginURL": "https://console.aws.amazon.com/ec2/home",
  "roleValuePattern": "arn:aws:iam::$${accountid}:saml-provider/OKTA,arn:aws:iam::$${accountid}:role/$${role}",
  "sessionDuration": 3600,
  "useGroupMapping": false
}
EOT
}
