---
page_title: "okta_group_rule Resource - terraform-provider-okta"
description: |-
  Manages an Okta Group Rule resource.
---

# okta_group_rule (Resource)

Manages an Okta Group Rule.

## Example Usage

```hcl
resource "okta_group_rule" "example" {

  # Optional fields
  # group_ids = "<group_ids>"
  # type = "<type>"
  # value = "<value>"
  # exclude = "<exclude>"
  # exclude = "<exclude>"
}
```

## Schema

### Optional

- `group_ids` (List) Array of `groupIds` to which users are added
- `type` (String) Expression type. Only valid value is '`urn:okta:expression:1.0`'.
- `value` (String) Okta expression that would result in a Boolean value
- `exclude` (List) Currently not supported
- `exclude` (List) Excluded `userIds` when processing rules
- `name` (String) Name of the group rule
- `status` (String) Status of group rule. You can't update the status of a rule from `INACTIVE` to `ACTIVE`. You must use the activate and deactivate lifecycle operations.
- `type` (String) Type to indicate a group rule operation. Only `group_rule` is allowed.

### Read-Only

- `id` (String) The unique identifier for the resource.
- `embedded` (String) This object appears with embedded resources related to the group rule if you use the `expand` query parameter
- `created` (String) Creation date for group rule
- `last_updated` (String) Date group rule was last updated
