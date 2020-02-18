---
layout: "okta"
page_title: "Okta: okta_user_profile_mapping_source"
sidebar_current: "docs-okta-datasource-user-profile-mapping-source"
description: |-
  Get the base user Profile Mapping source or target from Okta.
---

# okta_user_profile_mapping_source

Use this data source to retrieve the base user Profile Mapping source or target from Okta.

## Example Usage

```hcl
data "okta_user_profile_mapping_source" "example" {}
```

## Attributes Reference

* `id` - id of the source.

* `name` - name of source.

* `type` - type of source.
