---
layout: 'okta'
page_title: 'Okta: okta_org_metadata'
sidebar_current: 'docs-okta-datasource-org-metadata'
description: |-
  Retrieves the well-known org metadata, which includes the id, configured custom domains, authentication pipeline, and various other org settings.
---

# okta_org_metadata

Retrieves the well-known org metadata, which includes the id, configured custom domains, authentication pipeline, and various other org settings.

- [Org Well Known Metadata Reference](https://developer.okta.com/docs/api/openapi/okta-management/management/tag/OrgSetting/#tag/OrgSetting/operation/getWellknownOrgMetadata)

## Example Usage

```hcl
data "okta_org_metadata" "test" {}
```

## Argument Reference

No arguments are supported, this is based on your provider configuration.

## Attributes Reference

- `id` - The unique identifier of the Org.
- `pipeline` - The authentication pipeline of the org. idx means the org is using the Identity Engine, while v1 means the org is using the Classic authentication pipeline.
- `settings` - The wellknown org settings (safe for public consumption).
  - `analytics_collection_enabled`
  - `bug_reporting_enabled`
  - `om_enabled` - Whether the legacy Okta Mobile application is enabled for the org
- `domains` - The URIs for the org's configured domains.
  - `organization` - Standard Org URI
  - `alternate` - Custom Domain Org URI
