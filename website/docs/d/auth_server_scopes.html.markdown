---
layout: 'okta'
page_title: 'Okta: okta_auth_server_scopes'
sidebar_current: 'docs-okta-datasource-auth-server-scopes'
description: |-
  Get a list of authorization server scopes from Okta.
---

# okta_auth_server_scopes

Use this data source to retrieve a list of authorization server scopes from Okta.

## Example Usage

```hcl
data "okta_auth_server_scopes" "test" {
  auth_server_id = "default"
}
```

## Arguments Reference

- `auth_server_id` - (Required) Auth server ID.

## Attributes Reference

- `scopes` - collection of authorization server scopes retrieved from Okta with the following properties.
  - `id` - ID of the Scope
  - `name` - Name of the Scope
  - `description` - Description of the Scope
  - `consent` - Indicates whether a consent dialog is needed for the Scope
  - `metadata_publish` - Whether the Scope should be included in the metadata
  - `default` - Whether the Scope is a default Scope
  - `system` - Whether Okta created the Scope
