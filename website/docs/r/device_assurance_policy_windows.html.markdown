---
layout: 'okta'
page_title: 'Okta: okta_policy_device_assurance_windows'
sidebar_current: 'docs-okta-device-assurance-policy-windows'
description: |-
    Manages a device assurance policy for windows.
---

# okta_policy_device_assurance_windows

This resource allows you to create and configure an device assurance policy for windows.

## Example Usage

```hcl
resource okta_policy_device_assurance_windows example{
    name = "example"
    os_version = "12.4.6"
    disk_encryption_type = toset(["ALL_INTERNAL_VOLUMES"])
    secure_hardware_present = true
    screenlock_type = toset(["BIOMETRIC", "PASSCODE"])
    third_party_signal_providers = true
    tpsp_browser_version = "15393.27.0"
    tpsp_builtin_dns_client_enabled = true
    tpsp_chrome_remote_desktop_app_blocked = true
    tpsp_crowd_strike_agent_id = "exampleAgentId"
    tpsp_crowd_strike_customer_id = "exampleCustomerId"
    tpsp_device_enrollment_domain = "exampleDomain"
    tpsp_disk_encrypted = true
    tpsp_key_trust_level = "CHROME_BROWSER_HW_KEY"
    tpsp_os_firewall = true
    tpsp_os_version = "10.0.19041"
    tpsp_password_proctection_warning_trigger = "PASSWORD_PROTECTION_OFF"
    tpsp_realtime_url_check_mode = true
    tpsp_safe_browsing_protection_level = "ENHANCED_PROTECTION"
    tpsp_screen_lock_secured = true
    tpsp_secure_boot_enabled = true
    tpsp_site_isolation_enabled = true
    tpsp_third_party_blocking_enabled = true
    tpsp_windows_machine_domain = "exampleMachineDomain"
    tpsp_windows_user_domain = "exampleUserDomain"
    }
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) Name of the device assurance policy.

- `disk_encryption_type` - (Optional) List of disk encryption type of the device assurance policy.

- `os_version` - (Optional) Minimum os version of the device in the device assurance policy.

- `secure_hardware_present` - (Optional) Is the device secure with hardware in the device assurance policy.

- `screenlock_type` - (Optional) List of screen lock type of the device assurance policy.

- `third_party_signal_providers` - (Optional) Indicate where the device assurance is using third party signal provider. Must be set if you want to use other tpsp value

- `tpsp_browser_version` - (Optional) Third party signal provider minimum browser version.

- `tpsp_builtin_dns_client_enabled` - (Optional) Third party signal provider builtin dns client enabled.

- `tpsp_chrome_remote_desktop_app_blocked` - (Optional) Third party signal provider chrome remote desktop app blocked.

- `tpsp_crowd_strike_agent_id` - (Optional) Third party signal provider crowdstrike agent id.

- `tpsp_crowd_strike_customer_id` - (Optional) Third party signal provider crowdstrike customer id.

- `tpsp_device_enrollment_domain` - (Optional) Third party signal provider device enrollment domain.

- `tpsp_disk_encrypted` - (Optional) Third party signal provider disk encrypted.

- `tpsp_key_trust_level` - (Optional) Third party signal provider key trust level.

- `tpsp_os_firewall` - (Optional) Third party signal provider os firewall.

- `tpsp_os_version` - (Optional) Third party signal provider minimum os version.

- `tpsp_password_proctection_warning_trigger` - (Optional) Third party signal provider password protection warning trigger.

- `tpsp_realtime_url_check_mode` - (Optional) Third party signal provider realtime url check mode.

- `tpsp_safe_browsing_protection_level` - (Optional) Third party signal provider safe browsing protection level.

- `tpsp_screen_lock_secured` - (Optional) Third party signal provider screen lock secure.

- `tpsp_secure_boot_enabled` - (Optional) Third party signal provider secure boot enabled.

- `tpsp_site_isolation_enabled` - (Optional) Third party signal provider site isolation enabled.

- `tpsp_third_party_blocking_enabled` - (Optional) Third party signal provider third party blocking enabled.

- `tpsp_windows_machine_domain` - (Optional) Third party signal provider windows machine domain.

- `tpsp_windows_user_domain` - (Optional) Third party signal provider windows user domain.

## Attributes Reference

- `id` - ID of the device assurance policy.

- `platform` - Platform of the device assurance policy.

- `created_date` - Created date of the device assurance policy.

- `created_by` - Created by of the device assurance policy.

- `last_update` - Last update of the device assurance policy.

- `last_updated_by` - Last updated by of the device assurance policy.

## Import

Okta Device Assurance Windows can be imported via the Okta ID.

```
$ terraform import okta_policy_device_assurance_windows.example &#60;device assurance id&#62;
```