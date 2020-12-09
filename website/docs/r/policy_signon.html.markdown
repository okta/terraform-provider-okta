---
layout: 'okta'
page_title: 'Okta: okta_policy_signon'
sidebar_current: 'docs-okta-resource-policy-signon'
description: |-
  Creates a Sign On Policy.
---

# okta_policy_signon

Creates a Sign On Policy.

This resource allows you to create and configure a Sign On Policy.

## Example Usage

```hcl
resource "okta_policy_signon" "example" {
  name            = "example"
  status          = "ACTIVE"
  description     = "Example"
  groups_included = ["${data.okta_group.everyone.id}"]
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) Policy Name.

- `description` - (Optional) Policy Description.

- `priority` - (Optional) Priority of the policy.

- `status` - (Optional) Policy Status: `"ACTIVE"` or `"INACTIVE"`.

- `groups_included` - List of Group IDs to Include.

## Attributes Reference

- `id` - ID of the Policy.

## Import

A Sign On Policy can be imported via the Okta ID.

```
$ terraform import okta_policy_signon.example <policy id>
```
