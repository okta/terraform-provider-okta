resource "okta_policy_device_assurance_chromeos" "example" {
  name     = "My ChromeOS Policy"
  platform = "CHROMEOS"

  third_party_signal_providers {
    device_posture_id_p {
      compliant = true
    }
  }
}

resource "okta_policy_device_assurance_chromeos" "example_with_providers" {
  name                     = "My ChromeOS Policy with Providers"
  platform                 = "CHROMEOS"
  display_remediation_mode = "HIDE"

  third_party_signal_providers {
    device_posture_id_p {
      compliant = true
      managed   = true
    }

    dtc {
      allow_screen_lock                   = true
      disk_encrypted                      = true
      managed_device                      = true
      os_firewall                         = true
      screen_lock_secured                 = true
      site_isolation_enabled              = true
      realtime_url_check_mode             = true
      safe_browsing_protection_level      = "ENHANCED_PROTECTION"
      password_protection_warning_trigger = "PASSWORD_PROTECTION_OFF"
      device_enrollment_domain            = "example.com"
      key_trust_level                     = "CHROME_BROWSER_HW_KEY"

      os_version {
        minimum = "110"
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
