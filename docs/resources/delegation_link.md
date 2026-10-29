---
page_title: "okta_delegation_link Resource - terraform-provider-okta"
description: |-
  Manages an Okta Delegation Link resource.
---

# okta_delegation_link (Resource)

Manages an Okta Delegation Link. The Delegation Links API provides operations for managing delegation links.

## Example Usage

```terraform
resource "okta_delegation_link" "example" {
  from {
    type       = "OKTA_AUTHORIZATION_SERVER"
    client_orn = "orn:okta:directory:00otivnuq1X5554yc1d7:workload-principals:ai-agents:wlp114bmffd1byZg51d8"
    token_type = "ACCESS_TOKEN"
  }

  to {
    resource_orn = "orn:okta:directory:00otivnuq1X5554yc1d7:workload-principals:ai-agents:wlp11f2exa8aMDIS41d8"
  }
}
```

## Schema

### Required

- `from` (Block, Optional) From (see [below for nested schema](#nestedblock--from))
- `to` (Block, Optional) The target resource for the delegation link (see [below for nested schema](#nestedblock--to))

### Read-Only

- `id` (String) The unique identifier for the resource.

<a id="nestedblock--from"></a>
### Nested Schema for `from`

Required:

- `type` (String) The type of token source.

Optional:

- `app_instance_orn` (String) The [ORN](https://developer.
- `client_orn` (String) The [Okta Resource Name (ORN)](/openapi/okta-management/guides/roles/#okta-resource-name-orn) of the OAuth 2.
- `token_type` (String) The type of token accepted by the delegation link.


<a id="nestedblock--to"></a>
### Nested Schema for `to`

Required:

- `resource_orn` (String) The [Okta Resource Name (ORN)](/openapi/okta-management/guides/roles/#okta-resource-name-orn) of the target resource.

## Import

Import is supported using the following syntax:

```shell
terraform import okta_delegation_link.example <delegation_link_id>
```

Notes:

- `update` is not supported in place. Changing any attribute forces the resource to be replaced.