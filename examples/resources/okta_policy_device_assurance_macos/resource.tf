resource "okta_policy_device_assurance_macos" "example" {
  name                    = "My macOS Policy"
  platform                = "MACOS"
  secure_hardware_present = true
}

resource "okta_policy_device_assurance_macos" "example_with_providers" {
  name                     = "My macOS Policy with Providers"
  platform                 = "MACOS"
  display_remediation_mode = "HIDE"

  third_party_signal_providers {
    device_posture_id_p {
      compliant = true
      managed   = true
    }

    dtc {
      built_in_dns_client_enabled         = true
      disk_encrypted                      = true
      os_firewall                         = true
      screen_lock_secured                 = true
      site_isolation_enabled              = true
      realtime_url_check_mode             = true
      safe_browsing_protection_level      = "ENHANCED_PROTECTION"
      password_protection_warning_trigger = "PASSWORD_PROTECTION_OFF"
      device_enrollment_domain            = "example.com"
      key_trust_level                     = "CHROME_BROWSER_HW_KEY"

      os_version {
        minimum = "13"
      }

      browser_version {
        minimum = "110"
      }
    }
  }

  grace_period {
    type   = "DAYS"
    expiry = "30"
  }
}
