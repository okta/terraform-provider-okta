---
page_title: "Data Source: okta_policy_device_assurance_chromeos"
description: |-
  Retrieves a ChromeOS Device Assurance Policy.
---

# Data Source: okta_policy_device_assurance_chromeos

Retrieves a device assurance policy by `deviceAssuranceId`.

!> **BREAKING CHANGE**: This data source has been migrated to use Okta SDK v7, which introduces schema changes. Migration guide below.

## Breaking Changes

### 1. Data Source Replaces the Generic `okta_device_assurance_policy` Data Source
The previous generic `okta_device_assurance_policy` data source (queryable by `id` or `name`) has been replaced by platform-specific data sources, one per OS. This data source, `okta_policy_device_assurance_chromeos`, retrieves ChromeOS device assurance policies.

### 2. Argument Changes
- `id` is now **required**; querying by `name` is no longer supported
- `platform` is no longer exposed (the platform is implied by the data source itself)

### 3. Schema Changes
- Third-party signal provider attributes previously flattened with a `tpsp_` prefix (for example `tpsp_disk_encrypted`, `tpsp_os_firewall`, `tpsp_browser_version`) have been restructured into nested blocks under `third_party_signal_providers.dtc` (the Google Chrome Device Trust Connector provider), using names without the `tpsp_` prefix
- New computed attributes: `display_remediation_mode`, `grace_period` (`type`, `expiry`), and `third_party_signal_providers.device_posture_id_p` (`compliant`, `managed`)

## Example Usage

```terraform
resource "okta_policy_device_assurance_chromeos" "example" {
  name     = "My ChromeOS Policy"
  platform = "CHROMEOS"

  third_party_signal_providers {
    device_posture_id_p {
      compliant = true
    }
  }
}

data "okta_policy_device_assurance_chromeos" "example" {
  id = okta_policy_device_assurance_chromeos.example.id
}
```

## Argument Reference

- `id` - (Required) The unique identifier of the device assurance policy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `name` - Display name of the device assurance policy.
- `created_by` - User ID who created the policy.
- `created_date` - When the policy was created.
- `last_updated_by` - User ID who last updated the policy.
- `last_update` - When the policy was last updated.
- `display_remediation_mode` - Remediation mode for non-compliant devices.
- `third_party_signal_providers` - Third-party signal provider configuration.
  - `device_posture_id_p` - Device Posture IdP provider settings.
    - `compliant` - Whether device must be compliant.
    - `managed` - Whether device must be managed.
  - `dtc` - Google Chrome Device Trust Connector provider settings.
    - `allow_screen_lock` - Whether the AllowScreenLock enterprise policy is enabled.
    - `built_in_dns_client_enabled` - Whether a software DNS client stack is used.
    - `chrome_remote_desktop_app_blocked` - Whether Chrome Remote Desktop access is blocked.
    - `device_enrollment_domain` - Enrollment domain of the managing customer.
    - `disk_encrypted` - Whether the main disk is encrypted.
    - `key_trust_level` - Attestation strength used by the Chrome Verified Access API.
    - `managed_device` - Whether the device is enrolled in ChromeOS device management.
    - `os_firewall` - Whether an OS-level firewall is enabled.
    - `password_protection_warning_trigger` - Password Protection Warning feature state.
    - `realtime_url_check_mode` - Whether enterprise-grade unsafe URL scanning is enabled.
    - `safe_browsing_protection_level` - Current Safe Browsing protection level.
    - `screen_lock_secured` - Whether the device is password-protected.
    - `site_isolation_enabled` - Whether Site Isolation is enabled.
    - `browser_version` - Current version of the Chrome Browser.
      - `minimum` - Minimum browser version.
    - `os_version` - Current version of the operating system.
      - `minimum` - Minimum OS version.
- `grace_period` - Grace period configuration.
  - `type` - Type of grace period (e.g., "DAYS").
  - `expiry` - Grace period duration.
