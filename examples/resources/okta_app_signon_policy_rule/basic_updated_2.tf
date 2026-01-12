resource "okta_app_saml" "test" {
  label                     = "testAcc_replace_with_uuid"
  sso_url                   = "http://google.com"
  recipient                 = "http://here.com"
  destination               = "http://its-about-the-journey.com"
  audience                  = "http://audience.com"
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
  status   = "ACTIVE"
}

data "okta_user_type" "default" {
  name = "user"
}

resource "okta_policy_device_assurance_android" "test" {
  name       = "testAcc-replace_with_uuid"
  os_version = "12"
  jailbreak  = false
}

resource "okta_app_signon_policy_rule" "test" {
  name                 = "testAcc_replace_with_uuid_updated"
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
  # FIXME Okta API for /api/v1/policies/{policyId}/rules/{ruleId}
  # is not returning os_expression even when it has been set throwing off the TF state.
  #  platform_include {
  #    os_expression = ".*"
  #    os_type = "OTHER"
  #    type    = "DESKTOP"
  #  }
  #  platform_include {
  #    os_expression = ".*"
  #    os_type = "OTHER"
  #    type    = "MOBILE"
  #  }
  risk_score = "MEDIUM"
  platform_include {
    os_type = "WINDOWS"
    type    = "DESKTOP"
  }
  platform_include {
    os_type = "CHROMEOS"
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
        "types" : ["password"],
        "required" : false
      },
      "possession" : {
        "deviceBound" : "REQUIRED",
        "required" : false
      }
    }),
    jsonencode({
      "possession" : {
        "deviceBound" : "REQUIRED",
        "hardwareProtection" : "REQUIRED",
        "userPresence" : "OPTIONAL",
        "userVerification" : "OPTIONAL",
        "required" : false
      }
    })
  ]
}
