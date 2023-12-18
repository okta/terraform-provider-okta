resource "okta_policy_device_assurance_chromeos" "example" {
  name                                      = "example"
  tpsp_allow_screen_lock                    = true
  tpsp_browser_version                      = "15393.27.0"
  tpsp_builtin_dns_client_enabled           = true
  tpsp_chrome_remote_desktop_app_blocked    = true
  tpsp_device_enrollment_domain             = "exampleDomain"
  tpsp_disk_encrypted                       = true
  tpsp_key_trust_level                      = "CHROME_OS_VERIFIED_MODE"
  tpsp_os_firewall                          = true
  tpsp_os_version                           = "10.0.19041.1110"
  tpsp_password_proctection_warning_trigger = "PASSWORD_PROTECTION_OFF"
  tpsp_realtime_url_check_mode              = true
  tpsp_safe_browsing_protection_level       = "ENHANCED_PROTECTION"
  tpsp_screen_lock_secured                  = true
  tpsp_site_isolation_enabled               = true
}
