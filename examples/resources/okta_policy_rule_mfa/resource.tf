data "okta_default_policy" "example" {
  type = "MFA_ENROLL"
}

resource "okta_policy_rule_mfa" "example" {
  policy_id = data.okta_default_policy.example.id
  name      = "My Rule"
  status    = "ACTIVE"
  enroll    = "LOGIN"
  app_include {
    id   = okta_app_oauth.example.id
    type = "APP"
  }
  app_include {
    type = "APP_TYPE"
    name = "yahoo_mail"
  }
}

resource "okta_app_oauth" "example" {
  label          = "My App"
  type           = "web"
  grant_types    = ["authorization_code"]
  redirect_uris  = ["http://localhost:8000"]
  response_types = ["code"]
}
### Unchecked `Okta` and checked `Applications` (with `Any application that supports MFA enrollment` option) checkboxes in the `User is accessing` section corresponds to the following config:

data "okta_default_policy" "example" {
  type = "MFA_ENROLL"
}

resource "okta_policy_rule_mfa" "example" {
  name      = "Some policy rule"
  policy_id = data.okta_default_policy.example.id

  app_exclude {
    name = "okta"
    type = "APP_TYPE"
  }
}

### Unchecked `Okta` and checked `Applications` (with `Specific applications` option) checkboxes in the `User is accessing` section corresponds to the following config:

data "okta_default_policy" "example" {
  type = "MFA_ENROLL"
}

resource "okta_policy_rule_mfa" "example" {
  name      = "Some policy rule"
  policy_id = data.okta_default_policy.example.id

  app_exclude {
    name = "okta"
    type = "APP_TYPE"
  }

  app_include {
    id   = "some_app_id"
    type = "APP"
  }
}

### Checked `Okta` and unchecked `Applications` checkboxes in the `User is accessing` section corresponds to the following config:

data "okta_default_policy" "example" {
  type = "MFA_ENROLL"
}

resource "okta_policy_rule_mfa" "example" {
  name      = "Some policy rule"
  policy_id = data.okta_default_policy.example.id

  app_include {
    name = "okta"
    type = "APP_TYPE"
  }
}

### Checked `Okta` and checked `Applications` (with `Any application that supports MFA enrollment` option) checkboxes in the `User is accessing` section corresponds to the following config:

data "okta_default_policy" "example" {
  type = "MFA_ENROLL"
}

resource "okta_policy_rule_mfa" "example" {
  name      = "Some policy rule"
  policy_id = data.okta_default_policy.example.id
}

### Checked `Okta` and checked `Applications` (with `Specific applications` option) checkboxes in the `User is accessing` section corresponds to the following config:

data "okta_default_policy" "example" {
  type = "MFA_ENROLL"
}

resource "okta_policy_rule_mfa" "example" {
  name      = "Some policy rule"
  policy_id = data.okta_default_policy.example.id

  app_include {
    name = "okta"
    type = "APP_TYPE"
  }

  app_include {
    id   = "some_app_id"
    type = "APP"
  }
}
