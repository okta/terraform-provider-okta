resource "okta_policy_device_assurance_windows" "example" {
  name                    = "My Windows Policy"
  platform                = "WINDOWS"
  secure_hardware_present = true
}

resource "okta_policy_device_assurance_windows" "example_with_providers" {
  name                     = "My Windows Policy with Providers"
  platform                 = "WINDOWS"
  display_remediation_mode = "HIDE"

  third_party_signal_providers {
    device_posture_id_p {
      compliant = true
      managed   = true
    }

    dtc {
      antivirus_enabled                   = true
      built_in_dns_client_enabled         = true
      disk_encrypted                      = true
      os_firewall                         = true
      screen_lock_secured                 = true
      secure_boot_enabled                 = true
      site_isolation_enabled              = true
      third_party_blocking_enabled        = true
      realtime_url_check_mode             = true
      safe_browsing_protection_level      = "ENHANCED_PROTECTION"
      password_protection_warning_trigger = "PASSWORD_PROTECTION_OFF"
      device_enrollment_domain            = "example.com"
      key_trust_level                     = "CHROME_BROWSER_HW_KEY"
      windows_machine_domain              = "example.com"
      windows_user_domain                 = "example.com"
      crowd_strike_agent_id               = "agent-id"
      crowd_strike_customer_id            = "customer-id"

      os_version {
        minimum = "10"
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
