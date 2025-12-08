---
page_title: "Resource: okta_push_provider"
description: |-
  Manages Okta Push Providers that provide a centralized integration platform to fetch and manage push provider configurations. Okta administrators can use these APIs to provide their push provider credentials, for example from APNs and FCM, so that Okta can send push notifications to their own custom app authenticator applications..
---

# Resource: okta_push_provider

Manages Okta Push Providers that provide a centralized integration platform to fetch and manage push provider configurations. Okta administrators can use these APIs to provide their push provider credentials, for example from APNs and FCM, so that Okta can send push notifications to their own custom app authenticator applications.

## Example Usage

### FCM Push Provider

```terraform
resource "okta_push_provider" "fcm_example" {
  name          = "FCM Push Provider"
  provider_type = "FCM"
  
  configuration {
    fcm_configuration {
      service_account_json {
        type                        = "service_account"
        project_id                  = "my-firebase-project"
        private_key_id              = "abc123def456"
        private_key                 = "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC...\n-----END PRIVATE KEY-----\n"
        client_email                = "firebase-adminsdk-abc123@my-firebase-project.iam.gserviceaccount.com"
        client_id                   = "123456789012345678901"
        auth_uri                    = "https://accounts.google.com/o/oauth2/auth"
        token_uri                   = "https://oauth2.googleapis.com/token"
        auth_provider_x509_cert_url = "https://www.googleapis.com/oauth2/v1/certs"
        client_x509_cert_url        = "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-abc123%40my-firebase-project.iam.gserviceaccount.com"
        file_name                   = "service-account-key.json"
      }
    }
  }
}
```

### APNS (Apple Push Notification Service) Push Provider

```terraform
resource "okta_push_provider" "apns_example" {
  name          = "APNS Push Provider"
  provider_type = "APNS"
  
  configuration {
    apns_configuration {
      key_id             = "ABC123DEFG"  # 10-character Key ID from Apple
      team_id            = "DEF123GHIJ"  # 10-character Team ID from Apple
      token_signing_key  = "-----BEGIN PRIVATE KEY-----\nMIGTAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBHkwdwIBAQQg...\n-----END PRIVATE KEY-----\n"
      file_name          = "AuthKey_ABC123DEFG.p8"
    }
  }
}
```

## Schema

### Required

- `name` (String) The display name of the push provider.
- `provider_type` (String) The type of push provider. Valid values are `APNS` (Apple Push Notification Service) or `FCM` (Firebase Cloud Messaging).
- `configuration` (Block) Configuration block for the push provider. The configuration structure depends on the provider type. (see [below for nested schema](#nestedblock--configuration))

### Optional

- `last_updated_date` (String) Timestamp when the push provider was last modified. (Computed)

### Read-Only

- `id` (String) The unique identifier of the push provider.

<a id="nestedblock--configuration"></a>
### Nested Schema for `configuration`

The configuration block must contain exactly one of the following sub-blocks based on the `provider_type`:

#### For `provider_type = "FCM"`:

**fcm_configuration** (Block) Configuration for Firebase Cloud Messaging:

**service_account_json** (Block) JSON containing the private service account key and service account details obtained from Firebase Console:

- `type` (String) The type of the service account (typically "service_account").
- `project_id` (String) The Firebase project ID.
- `private_key_id` (String) The private key ID from the service account.
- `private_key` (String, Sensitive) The private key from the service account.
- `client_email` (String) The service account email address.
- `client_id` (String) The client ID from the service account.
- `auth_uri` (String) The authentication URI (typically "https://accounts.google.com/o/oauth2/auth").
- `token_uri` (String) The token URI (typically "https://oauth2.googleapis.com/token").
- `auth_provider_x509_cert_url` (String) The auth provider X509 certificate URL.
- `client_x509_cert_url` (String) The client X509 certificate URL.
- `file_name` (String) Optional file name for Admin Console display (e.g., "service-account-key.json").

#### For `provider_type = "APNS"`:

**apns_configuration** (Block) Configuration for Apple Push Notification Service:

- `key_id` (String) 10-character Key ID obtained from the Apple Developer account.
- `team_id` (String) 10-character Team ID used to develop the iOS app.
- `token_signing_key` (String, Sensitive) APNs private authentication token signing key in PEM format.
- `file_name` (String) Optional file name for Admin Console display (e.g., "AuthKey_ABC123DEFG.p8").

## Import

Push providers can be imported using their ID:

```shell
terraform import okta_push_provider.example <push_provider_id>
```
