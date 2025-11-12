---
page_title: "Data Source: okta_push_provider"
description: |-
  Retrieves information about an Okta Push Provider configuration.
---

# Data Source: okta_push_provider

Retrieves information about an Okta Push Provider configuration. This data source allows you to fetch details about existing push providers for Apple Push Notification Service (APNS) and Firebase Cloud Messaging (FCM).

## Example Usage

### Retrieve an APNS Push Provider

```terraform
data "okta_push_provider" "apns_example" {
  id = "ppc1234567890abcdef"
}

# Use the retrieved data
output "apns_provider_name" {
  value = data.okta_push_provider.apns_example.name
}

output "apns_key_id" {
  value = data.okta_push_provider.apns_example.configuration.apns_configuration.key_id
}
```

### Retrieve an FCM Push Provider

```terraform
data "okta_push_provider" "fcm_example" {
  id = "ppc0987654321fedcba"
}

# Use the retrieved data
output "fcm_provider_name" {
  value = data.okta_push_provider.fcm_example.name
}

output "fcm_project_id" {
  value = data.okta_push_provider.fcm_example.configuration.fcm_configuration.service_account_json.project_id
}
```

## Schema

### Required

- `id` (String) The unique identifier of the push provider to retrieve.

### Read-Only

- `name` (String) The display name of the push provider.
- `provider_type` (String) The type of push provider. Values are `APNS` (Apple Push Notification Service) or `FCM` (Firebase Cloud Messaging).
- `last_updated_date` (String) Timestamp when the push provider was last modified.
- `configuration` (Block) Configuration details for the push provider. The structure depends on the provider type. (see [below for nested schema](#nestedblock--configuration))

<a id="nestedblock--configuration"></a>
### Nested Schema for `configuration`

The configuration block will contain one of the following sub-blocks based on the `provider_type`:

#### For `provider_type = "FCM"`:

**fcm_configuration** (Block) Configuration details for Firebase Cloud Messaging:

**service_account_json** (Block) Service account information from Firebase:

- `project_id` (String) The Firebase project ID.
- `file_name` (String) The file name used for Admin Console display.

**Note**: Sensitive fields like private keys are not returned by the API for security reasons.

#### For `provider_type = "APNS"`:

**apns_configuration** (Block) Configuration details for Apple Push Notification Service:

- `key_id` (String) The 10-character Key ID from the Apple Developer account.
- `team_id` (String) The 10-character Team ID used to develop the iOS app.
- `file_name` (String) The file name used for Admin Console display.
