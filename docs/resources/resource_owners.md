---
page_title: "Resource: okta_resource_owners"
description: |-
  Manages owners for a governance resource (entitlement bundle, entitlement value, collection, or app).
---

# Resource: okta_resource_owners

Manages owners for a governance resource (entitlement bundle, entitlement value, collection, or app). Owners are principals (users or groups) authorized to administer the resource. A maximum of 5 principals may be assigned per resource.

## Example Usage

```terraform
# Assign owners to an entitlement bundle
resource "okta_resource_owners" "example" {
  resource_orn = "orn:okta:governance:00o1234567890abcdef:entitlement-bundles:enb1234567890abcdef"

  principal_orns = [
    "orn:okta:directory:00o1234567890abcdef:users:00u1234567890abcdef",
    "orn:okta:directory:00o1234567890abcdef:groups:00g1234567890abcdef",
  ]
}
```

## Argument Reference

- `resource_orn` - (Required) The ORN of the resource to manage owners for (e.g., an entitlement bundle, entitlement value, collection, or app).
- `principal_orns` - (Required) The ORNs of the principals (users or groups) that own the resource. Maximum 5.
- `parent_resource_orn` - (Optional, Computed) The ORN of the parent resource (typically the app). Required when importing. Can be found from the `target_resource_orn` attribute of `okta_entitlement_bundle`.

## Attributes Reference

- `id` - Resource ID (same as `resource_orn`).
- `parent_resource_orn` - Set by the API on create / read.

## Import

Import is supported using the format `parent_resource_orn/resource_orn`:

```shell
terraform import okta_resource_owners.example "orn:okta:idp:00o...:apps:salesforce:0oa.../orn:okta:governance:00o...:entitlement-bundles:enb..."
```

The `parent_resource_orn` can be found from:
- The `target_resource_orn` attribute of an `okta_entitlement_bundle` resource or data source
- The Okta Admin Console under governance settings
