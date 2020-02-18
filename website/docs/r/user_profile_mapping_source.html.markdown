---
layout: "okta"
page_title: "Okta: okta_profile_mapping"
sidebar_current: "docs-okta-resource-profile-mapping"
description: |-
  Manages a profile mapping.
---

# okta_profile_mapping

Manages a profile mapping.

This resource allows you to manage a profile mapping by source id.

## Example Usage

```hcl
data "okta_user_profile_mapping_source" "user" {}

resource "okta_profile_mapping" "example" {
  source_id          = "<source id>"
  target_id          = "${data.okta_user_profile_mapping_source.user.id}"
  delete_when_absent = true

  mappings {
    id         = "firstName"
    expression = "appuser.firstName"
  }

  mappings {
    id         = "lastName"
    expression = "appuser.lastName"
  }

  mappings {
    id         = "email"
    expression = "appuser.email"
  }

  mappings {
    id         = "login"
    expression = "appuser.email"
  }
}
```

## Argument Reference

The following arguments are supported:

* `source_id` - (Required) Source id of the profile mapping.

* `delete_when_absent` - (Optional) Tells the provider whether to attempt to delete missing mappings under profile mapping.

* `mappings` - (Optional) Priority of the policy.
  * `id` - (Required) Key of mapping.
  * `expression` - (Required) Combination or single source properties that will be mapped to the target property.
  * `push_status` - (Optional) Whether to update target properties on user create & update or just on create.

## Attributes Reference

* `id` - ID of the mappings.

* `target_id` - ID of the mapping target.

* `target_name` - Name of the mapping target.

* `target_type` - ID of the mapping target.

* `source_id` - ID of the mapping source.

* `source_name` - Name of the mapping source.

* `source_type` - ID of the mapping source.

## Import

There is no reason to import this resource. You can simply create the resource config and point it to a source ID. Once the source is deleted this resources will no longer exist.
