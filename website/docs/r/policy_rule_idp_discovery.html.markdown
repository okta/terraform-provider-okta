---
layout: "okta"
page_title: "Okta: okta_policy_rule_idp_discovery"
sidebar_current: "docs-okta-resource-policy-rule-idp-discovery"
description: |-
  Creates an IdP Discovery Policy Rule.
---

# okta_policy_rule_idp_discovery

Creates an IdP Discovery Policy Rule.

This resource allows you to create and configure an IdP Discovery Policy Rule.

## Example Usage

```hcl
resource "okta_policy_rule_idp_discovery" "example" {
  policyid                  = "<policy id>"
  priority                  = 1
  name                      = "example"
  idp_type                  = "SAML2"
  idp_id                    = "<idp id>"
  user_identifier_type      = "ATTRIBUTE"
  user_identifier_attribute = "company"

  user_identifier_patterns {
    match_type = "EQUALS"
    value      = "Articulate"
  }
}
```

## Argument Reference

The following arguments are supported:

* `policyid` - (Required) Policy ID.

* `name` - (Required) Policy Rule Name.

* `priority` - (Optional) Policy Rule Priority, this attribute can be set to a valid priority. To avoid endless diff situation we error if an invalid priority is provided. API defaults it to the last/lowest if not there.

* `status` - (Optional) Policy Rule Status: `"ACTIVE"` or `"INACTIVE"`.

* `network_connection` - (Optional) Network selection mode: `"ANYWHERE"`, `"ZONE"`, `"ON_NETWORK"`, or `"OFF_NETWORK"`.

* `network_includes` - (Optional) The network zones to include. Conflicts with `network_excludes`.

* `network_excludes` - (Optional) The network zones to exclude. Conflicts with `network_includes`.

## Attributes Reference

* `id` - ID of the Rule.

* `policyid` - Policy ID.

## Import

A Policy Rule can be imported via the Policy and Rule ID.

```
$ terraform import okta_policy_rule_idp_discovery.example <policy id>/<rule id>
```
