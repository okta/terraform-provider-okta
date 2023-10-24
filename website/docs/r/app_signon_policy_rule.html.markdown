---
layout: 'okta' 
page_title: 'Okta: okta_app_signon_policy_rule' 
sidebar_current: 'docs-okta-resource-okta-app-signon-policy-rule' 
description: |- 
    Manages a sign-on policy rules for the application.
---

# okta_app_signon_policy_rule

~> **WARNING:** This feature is only available as a part of the Identity Engine. [Contact support](mailto:dev-inquiries@okta.com) for further information.

This resource allows you to create and configure a sign-on policy rule for the application.

A default or `Catch-all Rule` sign-on policy rule can be imported and managed as a custom rule.
The only difference is that these fields are immutable and can not be managed: `network_connection`, `network_excludes`, 
`network_includes`, `platform_include`, `custom_expression`, `device_is_registered`, `device_is_managed`, `users_excluded`,
`users_included`, `groups_excluded`, `groups_included`, `user_types_excluded` and `user_types_included`.

## Example Usage

### Simple usage

```hcl
resource "okta_app_saml" "test" {
  label                    = "My App"
  sso_url                  = "https://google.com"
  recipient                = "https://here.com"
  destination              = "https://its-about-the-journey.com"
  audience                 = "https://audience.com"
  status                   = "ACTIVE"
  subject_name_id_template = "$${user.userName}"
  subject_name_id_format   = "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"
  signature_algorithm      = "RSA_SHA256"
  response_signed          = true
  digest_algorithm         = "SHA256"
  honor_force_authn        = false
  authn_context_class_ref  = "urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport"
}

data "okta_app_signon_policy" "test" {
  app_id = okta_app_saml.test.id
}

resource "okta_app_signon_policy_rule" "test" {
  policy_id = data.okta_app_signon_policy.test.id
  name      = "testAcc_replace_with_uuid"
}
```

This will create an app sign-on policy rule with the following `THEN` block:

```
THEN Access is 'Allowed after successful authentication'
AND User must authenticate with 'Any 2 fator types'
AND Possession factor constraints are '-'
AND Access with Okta FastPass is granted 'If the user approves a prompt in Okta Verify or provides biometrics (meets NIST AAL2 requirements)'
```

### Rule with Constraints

#### Example 1:

```hcl
resource "okta_app_signon_policy_rule" "test" {
  policy_id   = data.okta_app_signon_policy.test.id
  name        = "testAcc_replace_with_uuid"
  constraints = [
    jsonencode({
      "knowledge" : {
        "types" : ["password"]
      },
    })
  ]
}
```

This will create an app sign-on policy rule with the following `THEN` block:

```
THEN Access is 'Allowed after successful authentication'
AND User must authenticate with 'Password + Another factor'
AND Possession factor constraints are '-'
AND Access with Okta FastPass is granted 'If the user approves a prompt in Okta Verify or provides biometrics (meets NIST AAL2 requirements)'
```

#### Example 2:

```hcl
resource "okta_app_signon_policy_rule" "test" {
  policy_id   = data.okta_app_signon_policy.test.id
  name        = "testAcc_replace_with_uuid"
  constraints = [
    jsonencode(
    {
      knowledge  = {
        reauthenticateIn = "PT2H"
        types            = [
          "password",
        ]
      }
      possession = {
        deviceBound        = "REQUIRED"
        hardwareProtection = "REQUIRED"
      }
    }
    )
  ]
}
```

This will create an app sign-on policy rule with the following `THEN` block:

```
THEN Access is 'Allowed after successful authentication'
AND User must authenticate with 'Password + Another factor'
AND Possession factor constraints are 'Hardware protected' and 'Device Bound (excludes phone and email)'
AND Access with Okta FastPass is granted 'Without the user approving a prompt in Okta Verify or providing biometrics'
```

More examples can be
found [here](https://developer.okta.com/docs/reference/api/policy/#verification-method-json-examples).

### Complex example

```hcl
resource "okta_app_saml" "test" {
  label                     = "testAcc_replace_with_uuid"
  sso_url                   = "https://google.com"
  recipient                 = "https://here.com"
  destination               = "https://its-about-the-journey.com"
  audience                  = "https://audience.com"
  subject_name_id_template  = "$${user.userName}"
  subject_name_id_format    = "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"
  response_signed           = true
  signature_algorithm       = "RSA_SHA256"
  digest_algorithm          = "SHA256"
  honor_force_authn         = false
  authn_context_class_ref   = "urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport"
  single_logout_issuer      = "https://dunshire.okta.com"
  single_logout_url         = "https://dunshire.okta.com/logout"
  single_logout_certificate = "MIIFnDCCA4QCCQDBSLbiON2T1zANBgkqhkiG9w0BAQsFADCBjzELMAkGA1UEBhMCVVMxDjAMBgNV\r\nBAgMBU1haW5lMRAwDgYDVQQHDAdDYXJpYm91MRcwFQYDVQQKDA5Tbm93bWFrZXJzIEluYzEUMBIG\r\nA1UECwwLRW5naW5lZXJpbmcxDTALBgNVBAMMBFNub3cxIDAeBgkqhkiG9w0BCQEWEWVtYWlsQGV4\r\nYW1wbGUuY29tMB4XDTIwMTIwMzIyNDY0M1oXDTMwMTIwMTIyNDY0M1owgY8xCzAJBgNVBAYTAlVT\r\nMQ4wDAYDVQQIDAVNYWluZTEQMA4GA1UEBwwHQ2FyaWJvdTEXMBUGA1UECgwOU25vd21ha2VycyBJ\r\nbmMxFDASBgNVBAsMC0VuZ2luZWVyaW5nMQ0wCwYDVQQDDARTbm93MSAwHgYJKoZIhvcNAQkBFhFl\r\nbWFpbEBleGFtcGxlLmNvbTCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIBANMmWDjXPdoa\r\nPyzIENqeY9njLan2FqCbQPSestWUUcb6NhDsJVGSQ7XR+ozQA5TaJzbP7cAJUj8vCcbqMZsgOQAu\r\nO/pzYyQEKptLmrGvPn7xkJ1A1xLkp2NY18cpDTeUPueJUoidZ9EJwEuyUZIktzxNNU1pA1lGijiu\r\n2XNxs9d9JR/hm3tCu9Im8qLVB4JtX80YUa6QtlRjWR/H8a373AYCOASdoB3c57fIPD8ATDNy2w/c\r\nfCVGiyKDMFB+GA/WTsZpOP3iohRp8ltAncSuzypcztb2iE+jijtTsiC9kUA2abAJqqpoCJubNShi\r\nVff4822czpziS44MV2guC9wANi8u3Uyl5MKsU95j01jzadKRP5S+2f0K+n8n4UoV9fnqZFyuGAKd\r\nCJi9K6NlSAP+TgPe/JP9FOSuxQOHWJfmdLHdJD+evoKi9E55sr5lRFK0xU1Fj5Ld7zjC0pXPhtJf\r\nsgjEZzD433AsHnRzvRT1KSNCPkLYomznZo5n9rWYgCQ8HcytlQDTesmKE+s05E/VSWNtH84XdDrt\r\nieXwfwhHfaABSu+WjZYxi9CXdFCSvXhsgufUcK4FbYAHl/ga/cJxZc52yFC7Pcq0u9O2BSCjYPdQ\r\nDAHs9dhT1RhwVLM8RmoAzgxyyzau0gxnAlgSBD9FMW6dXqIHIp8yAAg9cRXhYRTNAgMBAAEwDQYJ\r\nKoZIhvcNAQELBQADggIBADofEC1SvG8qa7pmKCjB/E9Sxhk3mvUO9Gq43xzwVb721Ng3VYf4vGU3\r\nwLUwJeLt0wggnj26NJweN5T3q9T8UMxZhHSWvttEU3+S1nArRB0beti716HSlOCDx4wTmBu/D1MG\r\nt/kZYFJw+zuzvAcbYct2pK69AQhD8xAIbQvqADJI7cCK3yRry+aWtppc58P81KYabUlCfFXfhJ9E\r\nP72ffN4jVHpX3lxxYh7FKAdiKbY2FYzjsc7RdgKI1R3iAAZUCGBTvezNzaetGzTUjjl/g1tcVYij\r\nltH9ZOQBPlUMI88lxUxqgRTerpPmAJH00CACx4JFiZrweLM1trZyy06wNDQgLrqHr3EOagBF/O2h\r\nhfTehNdVr6iq3YhKWBo4/+RL0RCzHMh4u86VbDDnDn4Y6HzLuyIAtBFoikoKM6UHTOa0Pqv2bBr5\r\nwbkRkVUxl9yJJw/HmTCdfnsM9dTOJUKzEglnGF2184Gg+qJDZB6fSf0EAO1F6sTqiSswl+uHQZiy\r\nDaZzyU7Gg5seKOZ20zTRaX3Ihj9Zij/ORnrARE7eM/usKMECp+7syUwAUKxDCZkGiUdskmOhhBGL\r\nJtbyK3F2UvoJoLsm3pIcvMak9KwMjSTGJB47ABUP1+w+zGcNk0D5Co3IJ6QekiLfWJyQ+kKsWLKt\r\nzOYQQatrnBagM7MI2/T4\r\n"

  attribute_statements {
    type         = "GROUP"
    name         = "groups"
    filter_type  = "REGEX"
    filter_value = ".*"
  }
}

data "okta_app_signon_policy" "test" {
  app_id = okta_app_saml.test.id
}

resource "okta_user" "test" {
  count      = 5
  first_name = "TestAcc"
  last_name  = "Smith"
  login      = "testAcc_${count.index}@example.com"
  email      = "testAcc_${count.index}@example.com"
}

resource "okta_group" "this" {
  count       = 5
  name        = "testAcc_${count.index}"
  description = "testAcc_${count.index}"
}

resource "okta_user_type" "test" {
  name         = "testAcc_replace_with_uuid"
  display_name = "Terraform Acceptance Test User Type Updated"
  description  = "Terraform Acceptance Test User Type Updated"
}

resource "okta_network_zone" "test" {
  name     = "testAcc_replace_with_uuid"
  type     = "IP"
  gateways = ["1.2.3.4/24", "2.3.4.5-2.3.4.15"]
  proxies  = ["2.2.3.4/24", "3.3.4.5-3.3.4.15"]
}

data "okta_user_type" "default" {
  name = "user"
}

resource "okta_policy_device_assurance_android" "test" {
  name = "test"
  os_version = "12"
  jailbreak = false
}

resource "okta_app_signon_policy_rule" "test" {
  name                                   = "testAcc_replace_with_uuid"
  policy_id                              = data.okta_app_signon_policy.test.id
  access                                 = "ALLOW"
  custom_expression                      = "user.status == \"ACTIVE\""
  device_is_managed                      = false
  device_is_registered                   = true
  factor_mode                            = "2FA"
  groups_excluded                        = [
    okta_group.this[2].id,
    okta_group.this[3].id,
    okta_group.this[4].id
  ]
  groups_included                        = [
    okta_group.this[0].id,
    okta_group.this[1].id
  ]
  device_assurances_included = [
    okta_policy_device_assurance_android.test.id
  ]
  network_connection                     = "ZONE"
  network_includes                       = [
    okta_network_zone.test.id
  ]
  platform_include {
    os_type = "ANDROID"
    type    = "MOBILE"
  }
  platform_include {
    os_type = "IOS"
    type    = "MOBILE"
  }
  platform_include {
    os_type = "MACOS"
    type    = "DESKTOP"
  }
  platform_include {
    os_type = "OTHER"
    type    = "DESKTOP"
  }
  platform_include {
    os_type = "OTHER"
    type    = "MOBILE"
  }
  platform_include {
    os_type = "WINDOWS"
    type    = "DESKTOP"
  }
  priority                               = 98
  re_authentication_frequency            = "PT43800H"
  type                                   = "ASSURANCE"
  user_types_excluded                    = [
    okta_user_type.test.id
  ]
  user_types_included                    = [
    data.okta_user_type.default.id
  ]
  users_excluded                         = [
    okta_user.test[2].id,
    okta_user.test[3].id,
    okta_user.test[4].id
  ]
  users_included                         = [
    okta_user.test[0].id,
    okta_user.test[1].id
  ]
  constraints                            = [
    jsonencode({
      "knowledge" : {
        "reauthenticateIn" : "PT2H",
        "types" : ["password"]
      },
      "possession" : {
        "deviceBound" : "REQUIRED"
      }
    }),
    jsonencode({
      "possession" : {
        "deviceBound" : "REQUIRED",
        "hardwareProtection" : "REQUIRED",
        "userPresence" : "OPTIONAL"
      }
    })
  ]
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) Name of the policy rule.

- `policy_id` - (Required) ID of the app sign-on policy.

- `priority` - (Optional) Priority of the rule.

- `groups_included` - (Optional) List of groups IDs to be included.

- `groups_excluded` - (Optional) List of groups IDs to be excluded.

- `users_included` - (Optional) List of users IDs to be included.

- `users_excluded` - (Optional) List of users IDs to be excluded.

- `network_connection` - (Optional) Network selection mode: `"ANYWHERE"`, `"ZONE"`, `"ON_NETWORK"`, or `"OFF_NETWORK"`.

- `network_includes` - (Optional) List of network zones IDs to include. Conflicts with `network_excludes`.

- `network_excludes` - (Optional) List of network zones IDs to exclude. Conflicts with `network_includes`.

- `device_is_registered` - (Optional) If the device is registered. A device is registered if the User enrolls with Okta
  Verify that is installed on the device. Can only be set to `true`.

- `device_is_managed` - (Optional) If the device is managed. A device is managed if it's managed by a device management
  system. When managed is passed, `device_is_registered` must also be included and must be set to `true`.

- `device_assurances_included` - (Optional) List of device assurances IDs to be included.

- `platform_include` - (Optional) List of particular platforms or devices to match on.
    - `type` - (Optional) One of: `"ANY"`, `"MOBILE"`, `"DESKTOP"`
    - `os_expression` - (Optional) Only available and required when using `os_type = "OTHER"`
    - `os_type` - (Optional) One of: `"ANY"`, `"IOS"`, `"WINDOWS"`, `"ANDROID"`, `"OTHER"`, `"OSX"`, `"MACOS"`

- `custom_expression` - (Optional) This is an advanced optional setting. If the expression is formatted incorrectly or conflicts with conditions set above, the rule may not match any users.

- `user_types_excluded` - (Optional) List of user types IDs to be excluded.

- `user_types_included` - (Optional) List of user types IDs to be included.

- `access` - (Optional) Allow or deny access based on the rule conditions. It can be set to `"ALLOW"` or `"DENY"`. Default is `"ALLOW"`.

- `factor_mode` - (Optional) The number of factors required to satisfy this assurance level. It can be set to `"1FA"` or `"2FA"`. Default is `"2FA"`.

- `type` - (Optional) The Verification Method type. It can be set to `"ASSURANCE"`. Default is `"ASSURANCE"`.

- `re_authentication_frequency` - (Optional) The duration after which the end user must re-authenticate, regardless of user activity. Use the ISO 8601 Period format for recurring time intervals. `"PT0S"` - every sign-in attempt, `"PT43800H"` - once per session. Default is `"PT2H"`.

- `inactivity_period` - (Optional) The inactivity duration after which the end user must re-authenticate. Use the ISO 8601 Period format for recurring time intervals. Default is `"PT1H"`.

- `constraints` - (Optional) - An array that contains nested Authenticator Constraint objects that are organized by the Authenticator class. Each element should be in JSON format.

- `risk_score` - (Optional) - The risk score specifies a particular level of risk to match on. Valid values are: `"ANY"`, `"LOW"`, `"MEDIUM"`, `"HIGH"`. Default is `"ANY"`.

## Attributes Reference

- `id` - ID of the sign-on policy rule.

## Import

Okta app sign-on policy rule can be imported via the Okta ID.

```
$ terraform import okta_app_signon_policy_rule.example &#60;policy_id&#62;/&#60;rule_id&#62;
```
