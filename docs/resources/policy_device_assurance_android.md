---
page_title: "Resource: okta_policy_device_assurance_android"
description: |-
  Manages device assurance on policy
---

# Resource: okta_policy_device_assurance_android

Manages device assurance on policy



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Policy device assurance name

### Optional

- `disk_encryption_type` (Set of String) List of disk encryption type, can be FULL, USER
- `jailbreak` (Boolean) The device jailbreak. Only for android and iOS platform
- `os_version` (String) The device os minimum version
- `screenlock_type` (Set of String) List of screenlock type, can be BIOMETRIC or BIOMETRIC, PASSCODE
- `secure_hardware_present` (Boolean) Indicates if the device contains a secure hardware functionality

### Read-Only

- `created_by` (String) Created by
- `created_date` (String) Created date
- `id` (String) Policy assurance id
- `last_update` (String) Last update
- `last_updated_by` (String) Last updated by
- `platform` (String) Policy device assurance platform


