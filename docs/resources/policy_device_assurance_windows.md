---
page_title: "Resource: okta_policy_device_assurance_windows"
description: |-
  Manages device assurance on policy
---

# Resource: okta_policy_device_assurance_windows

Manages device assurance on policy



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Policy device assurance name

### Optional

- `disk_encryption_type` (Set of String) List of disk encryption type, can be ALL_INTERNAL_VOLUMES
- `os_version` (String) The device os minimum version
- `screenlock_type` (Set of String) List of screenlock type, can be BIOMETRIC or BIOMETRIC, PASSCODE
- `secure_hardware_present` (Boolean) Indicates if the device contains a secure hardware functionality
- `third_party_signal_providers` (Boolean) Check to include third party signal provider
- `tpsp_browser_version` (String) Third party signal provider minimum browser version
- `tpsp_builtin_dns_client_enabled` (Boolean) Third party signal provider builtin dns client enable
- `tpsp_chrome_remote_desktop_app_blocked` (Boolean) Third party signal provider chrome remote desktop app blocked
- `tpsp_crowd_strike_agent_id` (String) Third party signal provider crowdstrike agent id
- `tpsp_crowd_strike_customer_id` (String) Third party signal provider crowdstrike user id
- `tpsp_device_enrollment_domain` (String) Third party signal provider device enrollment domain
- `tpsp_disk_encrypted` (Boolean) Third party signal provider disk encrypted
- `tpsp_key_trust_level` (String) Third party signal provider key trust level
- `tpsp_os_firewall` (Boolean) Third party signal provider os firewall
- `tpsp_os_version` (String) Third party signal provider minimum os version
- `tpsp_password_proctection_warning_trigger` (String) Third party signal provider password protection warning trigger
- `tpsp_realtime_url_check_mode` (Boolean) Third party signal provider realtime url check mode
- `tpsp_safe_browsing_protection_level` (String) Third party signal provider safe browsing protection level
- `tpsp_screen_lock_secured` (Boolean) Third party signal provider screen lock secure
- `tpsp_secure_boot_enabled` (Boolean) Third party signal provider secure boot enabled
- `tpsp_site_isolation_enabled` (Boolean) Third party signal provider site isolation enabled
- `tpsp_third_party_blocking_enabled` (Boolean) Third party signal provider third party blocking enabled
- `tpsp_windows_machine_domain` (String) Third party signal provider windows machine domain
- `tpsp_windows_user_domain` (String) Third party signal provider windows user domain

### Read-Only

- `created_by` (String) Created by
- `created_date` (String) Created date
- `id` (String) Policy assurance id
- `last_update` (String) Last update
- `last_updated_by` (String) Last updated by
- `platform` (String) Policy device assurance platform


