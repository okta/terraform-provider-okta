resource "okta_policy_device_assurance_macos" "example" {
  name                                      = "testAcc-replace_with_uuid"
  os_version                                = "12.4.6"
  disk_encryption_type                      = toset(["ALL_INTERNAL_VOLUMES"])
  secure_hardware_present                   = true
  screenlock_type                           = toset(["BIOMETRIC", "PASSCODE"])
  third_party_signal_providers              = true
  tpsp_browser_version                      = "15393.27.0"
  tpsp_builtin_dns_client_enabled           = true
  tpsp_chrome_remote_desktop_app_blocked    = true
  tpsp_device_enrollment_domain             = "exampleDomain"
  tpsp_disk_encrypted                       = true
  tpsp_key_trust_level                      = "CHROME_BROWSER_HW_KEY"
  tpsp_os_firewall                          = true
  tpsp_os_version                           = "10.0.19041"
  tpsp_password_proctection_warning_trigger = "PASSWORD_PROTECTION_OFF"
  tpsp_realtime_url_check_mode              = true
  tpsp_safe_browsing_protection_level       = "ENHANCED_PROTECTION"
  tpsp_screen_lock_secured                  = true
  tpsp_site_isolation_enabled               = true
}
