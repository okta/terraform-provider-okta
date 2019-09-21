---
layout: "okta"
page_title: "Okta: okta_auth_server"
sidebar_current: "docs-okta-datasource-auth-server"
description: |-
  Get an auth server from Okta.
---

# okta_auth_server

Use this data source to retrieve the collaborators for a given repository.

## Example Usage

```hcl
data "okta_auth_server" "example" {
  label = "Example App"
}
```

## Arguments Reference

 * `label` - (Optional) The label of the app to retrieve, conflicts with `id`.

 * `id` - (Optional) `id` of application to retrieve, conflicts with `label`.

 * `active_only` - (Optional) tells the provider to query for only `ACTIVE` applications.

## Attributes Reference

 * `id` - `id` of application.

 * `label` - `label` of application.

 * `description` - `description` of application.

 * `name` - `name` of application.

 * `status` - `status` of application.
