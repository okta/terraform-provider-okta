---
page_title: "okta_service_account Resource - terraform-provider-okta"
description: |-
  Manages an Okta Service Account resource.
---

# okta_service_account (Resource)

Manages an Okta Service Account.

## Example Usage

```hcl
resource "okta_service_account" "example" {
  container_orn = "<container_orn>"
  name = "Example Name"
  username = "Example Username"

  # Optional fields
  # description = "Example description"
  # owner_group_ids = "<owner_group_ids>"
  # owner_user_ids = "<owner_user_ids>"
  # password = "<password>"
}
```

## Schema

### Required

- `container_orn` (String) The [ORN](/openapi/okta-management/guides/roles/#okta-resource-name-orn) of the relevant resource.  Use the specific app ORN format (`orn:{partition}:idp:{yourOrgId}:apps:{appType}:{appId}`) to ide...
- `name` (String) The user-defined name for the app service account
- `username` (String) The username that serves as the direct link to your managed app account. Ensure that this value precisely matches the identifier of the target app account.

### Optional

- `description` (String) The description of the app service account
- `owner_group_ids` (List) A list of IDs of the Okta groups who own the app service account
- `owner_user_ids` (List) A list of IDs of the Okta users who own the app service account
- `password` (String) The app service account password. Required for apps that don't have provisioning enabled or don't support password synchronization.

### Read-Only

- `id` (String) The unique identifier for the resource.
- `container_global_name` (String) The key name of the app in the Okta Integration Network (OIN)
- `container_instance_name` (String) The app instance label
- `created` (String) Timestamp when the app service account was created
- `last_updated` (String) Timestamp when the app service account was last updated
- `status` (String) Describes the current status of a service account
- `status_detail` (String) Describes the detailed status of a service account
