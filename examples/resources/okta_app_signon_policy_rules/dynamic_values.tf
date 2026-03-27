resource "okta_app_saml" "test_dynamic" {
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

data "okta_app_signon_policy" "test_dynamic" {
  app_id = okta_app_saml.test_dynamic.id
}

# Network zone whose ID is unknown at plan time — this exercises the
# "unknown value in ListNestedBlock" code path that caused the
# "Value Conversion Error" when Rules was a native Go slice.
resource "okta_network_zone" "test_dynamic" {
  name     = "testAcc_replace_with_uuid-zone"
  type     = "IP"
  status   = "ACTIVE"
  gateways = ["203.0.113.0/24"]
}

# Policy rules defined as a local map, driven via dynamic "rule" blocks.
# The network_includes attribute references the network zone ID, which is
# (known after apply) on the first plan — this is the exact scenario that
# previously triggered the Value Conversion Error.
locals {
  policy_rules_testAcc_replace_with_uuid = {
    "allow_2fa" = {
      name               = "Allow-2FA-testAcc_replace_with_uuid"
      priority           = 1
      status             = "ACTIVE"
      access             = "ALLOW"
      factor_mode        = "2FA"
      network_connection = "ZONE"
      network_includes   = [okta_network_zone.test_dynamic.id]
    }
    "allow_1fa" = {
      name               = "Allow-1FA-testAcc_replace_with_uuid"
      priority           = 2
      status             = "ACTIVE"
      access             = "ALLOW"
      factor_mode        = "1FA"
      network_connection = "ANYWHERE"
      network_includes   = []
    }
    "deny_all" = {
      name               = "Deny-All-testAcc_replace_with_uuid"
      priority           = 3
      status             = "ACTIVE"
      access             = "DENY"
      factor_mode        = "2FA"
      network_connection = "ANYWHERE"
      network_includes   = []
    }
  }
}

resource "okta_app_signon_policy_rules" "policy_rules" {
  policy_id = data.okta_app_signon_policy.test_dynamic.id

  dynamic "rule" {
    for_each = local.policy_rules_testAcc_replace_with_uuid
    content {
      name               = rule.value.name
      priority           = rule.value.priority
      status             = rule.value.status
      access             = rule.value.access
      factor_mode        = rule.value.factor_mode
      network_connection = rule.value.network_connection
      network_includes   = rule.value.network_includes
    }
  }
}
