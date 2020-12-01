// This is to server as a example of a full stack including Okta resources.
// This example is of a serverless SPA that allows users to authenticate via Okta credentials and invoke a Lambda function.
// For brevity only the core resources are defined.
locals = {
  // Potentially could be scripted https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_providers_create_oidc_verify-thumbprint.html
  jwk_thumbprint = "myThumbprint"
  uri            = "https://something.com"
}

provider "okta" {
  org_name = "something"
  base_url = "okta.com"
}

data "okta_group" "peeps" {
  name        = "Peeps"
  description = "For my peeps"
}

data "okta_user" "garth" {
  first_name        = "Garth"
  last_name         = "B"
  login             = "garth@something.com"
  email             = "garth@something.com"
  group_memberships = [okta_group.peeps.id]
}

resource "okta_app_oauth" "app" {
  label          = "G Studios"
  type           = "browser"
  client_uri     = local.uri
  redirect_uris  = [local.uri]
  groups         = [okta_group.peeps.id]
  response_types = ["token", "id_token"]
  grant_types    = ["implicit"]
}

resource "okta_trusted_origin" "app" {
  name   = "Something"
  origin = local.uri
  scopes = ["CORS", "REDIRECT"]
}

// AWS resources for auth to bring it full circle
resource "aws_cognito_identity_pool" "slick_idpool" {
  identity_pool_name               = "Cool Stuff Neat Stuff Slick stuff"
  allow_unauthenticated_identities = false

  openid_connect_provider_arns = [aws_iam_openid_connect_provider.openid.arn]
}

resource "aws_iam_openid_connect_provider" "openid" {
  url = "https://articulate.okta.com"

  client_id_list = [
    okta_app_oauth.app.id,
  ]

  thumbprint_list = [
    local.jwk_thumbprint,
  ]
}

resource "aws_iam_role" "authenticated" {
  name = "cognito_authenticated"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Federated": "cognito-identity.amazonaws.com"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringEquals": {
          "cognito-identity.amazonaws.com:aud": "${aws_cognito_identity_pool.slick_idpool.id}"
        },
        "ForAnyValue:StringLike": {
          "cognito-identity.amazonaws.com:amr": "authenticated"
        }
      }
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "authenticated" {
  name = "auth"
  role = aws_iam_role.authenticated.id

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "mobileanalytics:PutEvents",
        "cognito-sync:*",
        "cognito-identity:*"
      ],
      "Resource": [
        "*"
      ]
    },
    {
        "Effect": "Allow",
        "Action": [
            "lambda:InvokeFunction",
            "lambda:InvokeAsync"
        ],
        "Resource": "${aws_lambda_function.doSlickStuff.arn}"
    }
  ]
}
EOF
}

resource "aws_cognito_identity_pool_roles_attachment" "role_attachment" {
  identity_pool_id = aws_cognito_identity_pool.slick_idpool.id

  role_mapping {
    identity_provider         = aws_iam_openid_connect_provider.openid.arn
    ambiguous_role_resolution = "AuthenticatedRole"
    type                      = "Rules"

    mapping_rule {
      claim      = "staff"
      match_type = "Equals"
      role_arn   = aws_iam_role.authenticated.arn
      value      = "true"
    }

    mapping_rule {
      claim      = "aud"
      match_type = "Equals"
      role_arn   = aws_iam_role.authenticated.arn

      // Need to add data source for apps to Okta provider, oauth app only created in stage so hardcoding for now
      value = okta_app_oauth.app.id
    }
  }

  roles {
    authenticated = aws_iam_role.authenticated.arn
  }
}
