---
page_title: "Resource: okta_app_signon_policy_rule"
description: |-
  Manages a sign-on policy rules for the application.
  ~> WARNING: This feature is only available as a part of the Identity Engine. Contact support mailto:dev-inquiries@okta.com for further information.
  This resource allows you to create and configure a sign-on policy rule for the application.
  A default or 'Catch-all Rule' sign-on policy rule can be imported and managed as a custom rule.
  The only difference is that these fields are immutable and can not be managed: 'networkconnection', 'networkexcludes',
  'networkincludes', 'platforminclude', 'customexpression', 'deviceisregistered', 'deviceismanaged', 'usersexcluded',
  'usersincluded', 'groupsexcluded', 'groupsincluded', 'usertypesexcluded' and 'usertypes_included'.
---

# Resource: okta_app_signon_policy_rule

Manages a sign-on policy rules for the application.

~> **WARNING:** This feature is only available as a part of the Identity Engine. [Contact support](mailto:dev-inquiries@okta.com) for further information.

~> **WARNING:** When managing multiple `okta_app_signon_policy_rule` resources with concurrent operations, the Okta API may encounter concurrency issues. While this provider implements internal locking to prevent conflicts within a single Terraform process, you should use explicit `depends_on` references between rules to ensure proper sequencing, especially when managing rule priorities.

This resource allows you to create and configure a sign-on policy rule for the application.
A default or 'Catch-all Rule' sign-on policy rule can be imported and managed as a custom rule.
The only difference is that these fields are immutable and can not be managed: 'network_connection', 'network_excludes', 
'network_includes', 'platform_include', 'custom_expression', 'device_is_registered', 'device_is_managed', 'users_excluded',
'users_included', 'groups_excluded', 'groups_included', 'user_types_excluded' and 'user_types_included'.

## Example Usage

```terraform
### Simple usage

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

#### This will create an app sign-on policy rule with the following `THEN` block:

#### THEN Access is 'Allowed after successful authentication'
#### AND User must authenticate with 'Any 2 fator types'
#### AND Possession factor constraints are '-'
#### AND Access with Okta FastPass is granted 'If the user approves a prompt in Okta Verify or provides biometrics (meets NIST AAL2 requirements)'

### Rule with Constraints

#### Example 1:

resource "okta_app_signon_policy_rule" "test" {
  policy_id = data.okta_app_signon_policy.test.id
  name      = "testAcc_replace_with_uuid"
  constraints = [
    jsonencode({
      "knowledge" : {
        "types" : ["password"]
      },
    })
  ]
}

#### This will create an app sign-on policy rule with the following `THEN` block:


#### THEN Access is 'Allowed after successful authentication'
#### AND User must authenticate with 'Password + Another factor'
#### AND Possession factor constraints are '-'
#### AND Access with Okta FastPass is granted 'If the user approves a prompt in Okta Verify or provides biometrics (meets NIST AAL2 requirements)'

#### Example 2:

resource "okta_app_signon_policy_rule" "test" {
  policy_id = data.okta_app_signon_policy.test.id
  name      = "testAcc_replace_with_uuid"
  constraints = [
    jsonencode(
      {
        knowledge = {
          reauthenticateIn = "PT2H"
          types = [
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


#### This will create an app sign-on policy rule with the following `THEN` block:


#### THEN Access is 'Allowed after successful authentication'
#### AND User must authenticate with 'Password + Another factor'
#### AND Possession factor constraints are 'Hardware protected' and 'Device Bound (excludes phone and email)'
#### AND Access with Okta FastPass is granted 'Without the user approving a prompt in Okta Verify or providing biometrics'


#### More examples can be
#### found [here](https://developer.okta.com/docs/reference/api/policy/#verification-method-json-examples).

### Complex example

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
  name       = "test"
  os_version = "12"
  jailbreak  = false
}

resource "okta_app_signon_policy_rule" "test" {
  name                 = "testAcc_replace_with_uuid"
  policy_id            = data.okta_app_signon_policy.test.id
  access               = "ALLOW"
  custom_expression    = "user.status == \"ACTIVE\""
  device_is_managed    = false
  device_is_registered = true
  factor_mode          = "2FA"
  groups_excluded = [
    okta_group.this[2].id,
    okta_group.this[3].id,
    okta_group.this[4].id
  ]
  groups_included = [
    okta_group.this[0].id,
    okta_group.this[1].id
  ]
  device_assurances_included = [
    okta_policy_device_assurance_android.test.id
  ]
  network_connection = "ZONE"
  network_includes = [
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
  priority                    = 98
  re_authentication_frequency = "PT43800H"
  type                        = "ASSURANCE"
  user_types_excluded = [
    okta_user_type.test.id
  ]
  user_types_included = [
    data.okta_user_type.default.id
  ]
  users_excluded = [
    okta_user.test[2].id,
    okta_user.test[3].id,
    okta_user.test[4].id
  ]
  users_included = [
    okta_user.test[0].id,
    okta_user.test[1].id
  ]
  constraints = [
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

resource "okta_app_signon_policy_rule" "test" {
  access                      = "ALLOW"
  custom_expression           = null
  device_assurances_included  = null
  device_is_managed           = null
  device_is_registered        = null
  factor_mode                 = "2FA"
  groups_excluded             = null
  groups_included             = null
  inactivity_period           = "PT1H"
  name                        = "test2"
  network_connection          = "ANYWHERE"
  network_excludes            = null
  network_includes            = null
  policy_id                   = okta_app_signon_policy.test.id
  priority                    = 0
  re_authentication_frequency = "PT0S"
  status                      = "ACTIVE"
  type                        = "AUTH_METHOD_CHAIN"
  user_types_excluded         = []
  user_types_included         = []
  users_excluded              = []
  users_included              = []
  platform_include {
    os_expression = ""
    os_type       = "OTHER"
    type          = "DESKTOP"
  }
  chains = [jsonencode(
    {
      "authenticationMethods" : [
        {
          "key" : "okta_password",
          "method" : "password"
        }
      ],
      "next" : [{
        "authenticationMethods" : [{
          "key" : "okta_email",
          "method" : "email"
        }]
      }]
    }
  )]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Policy Rule Name
- `policy_id` (String) ID of the policy

### Optional

- `access` (String) Allow or deny access based on the rule conditions: ALLOW or DENY
- `constraints` (List of String) An array that contains nested Authenticator Constraint objects that are organized by the Authenticator class
- `custom_expression` (String) This is an optional advanced setting. If the expression is formatted incorrectly or conflicts with conditions set above, the rule may not match any users.
- `device_assurances_included` (Set of String) List of device assurance IDs to include
- `device_is_managed` (Boolean) If the device is managed. A device is managed if it's managed by a device management system. When managed is passed, registered must also be included and must be set to true.
- `device_is_registered` (Boolean) If the device is registered. A device is registered if the User enrolls with Okta Verify that is installed on the device.
- `factor_mode` (String) The number of factors required to satisfy this assurance level
- `groups_excluded` (Set of String) List of group IDs to exclude
- `groups_included` (Set of String) List of group IDs to include
- `inactivity_period` (String) The inactivity duration after which the end user must re-authenticate. Use the ISO 8601 Period format for recurring time intervals.
- `network_connection` (String) Network selection mode: ANYWHERE, ZONE, ON_NETWORK, or OFF_NETWORK.
- `network_excludes` (List of String) The zones to exclude
- `network_includes` (List of String) The zones to include
- `platform_include` (Block Set) (see [below for nested schema](#nestedblock--platform_include))
- `priority` (Number) Priority of the rule.
- `re_authentication_frequency` (String) The duration after which the end user must re-authenticate, regardless of user activity. Use the ISO 8601 Period format for recurring time intervals. PT0S - Every sign-in attempt, PT43800H - Once per session. Cannot be set if reauthenticateIn is set in one or more entries of chains.
- `risk_score` (String) The risk score specifies a particular level of risk to match on: ANY, LOW, MEDIUM, HIGH
- `status` (String) Status of the rule
- `type` (String) The Verification Method type
- `user_types_excluded` (Set of String) Set of User Type IDs to exclude
- `user_types_included` (Set of String) Set of User Type IDs to include
- `users_excluded` (Set of String) Set of User IDs to exclude
- `users_included` (Set of String) Set of User IDs to include
- `chains` (Block List) Authentication method chains. Only supports 5 items in the array. Each chain can support maximum 3 steps. To be used only with verification method type `AUTH_METHOD_CHAIN`.(see [below for nested schema](#nestedblock--chains))

### Read-Only

- `id` (String) The ID of this resource.
- `system` (Boolean) Often the `Catch-all Rule` this rule is the system (default) rule for its associated policy

<a id="nestedblock--platform_include"></a>
### Nested Schema for `platform_include`

Optional:

- `os_expression` (String) Only available with OTHER OS type
- `os_type` (String)
- `type` (String)

<a id="nestedblock--chains"></a>
### Nested Schema for `chains`

Optional:

- `authenticationMethods` (List of Authentication Methods) (see [below for nested schema](#nestedblock--authenticationMethods))
- `next` (List) The next steps of the authentication method chain. This is an array of type `chains`. Only supports one item in the array.
- `reauthenticateIn` (String) Specifies how often the user is prompted for authentication using duration format for the time period. This parameter can't be set at the same time as the `re_authentication_frequency` field.

<a id="nestedblock--authenticationMethods"></a>
### Nested Schema for `authenticationMethods`
Required:
- `key` (String) A label that identifies the authenticator.
- `method` (String) Specifies the method used for the authenticator.

Optional:
- `hardwareProtection` (String) Indicates if any secrets or private keys used during authentication must be hardware protected and not exportable. This property is only set for "POSSESSION" constraints. Set to "OPTIONAL" by default. Can only be set to "OPTIONAL" or "REQUIRED".
- `id` (String) An ID that identifies the authenticator
- `phishingResistant` (String) Indicates if phishing-resistant Factors are required. This property is only set for POSSESSION constraints. Set to "OPTIONAL" by default. Can only be set to "OPTIONAL" or "REQUIRED".
- `userVerification` (String) Indicates if a user is required to be verified with a verification method. Set to "OPTIONAL" by default. Can only be set to "OPTIONAL" or "REQUIRED".
- `userVerificationMethods` (List of String) Indicates which methods can be used for user verification. `userVerificationMethods` can only be used when `userVerification` is REQUIRED. BIOMETRICS is currently the only supported method.

## Import

Import is supported using the following syntax:

```shell
terraform import okta_app_signon_policy_rule.example <policy_id>/<rule_id>
```
