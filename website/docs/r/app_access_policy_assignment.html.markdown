---
layout: 'okta'
page_title: 'Okta: okta_app_access_policy_assignment'
sidebar_current: 'docs-okta-resource-group'
description: |-
  Assigns an access policy to an application.
---

# okta_app_access_policy_assignment

Assigns an access policy (colloquially known as a sign-on policy and/or an
authentication policy) to an application. This resource does not perform true
delete as it will not delete an application and the app's access policy can't be
removed; it can only be changed to a different access policy. This resource is
only logical within the context of an application therefore `app_id` is
immutable once set. Use this resource to manage assigning an access policy to an
application. It will assign the given `policy_id` to the application at creation
and during update.

-> Inside the product a sign-on policy is referenced as an _authentication
policy_, in the public API the policy is of type
[`ACCESS_POLICY`](https://developer.okta.com/docs/reference/api/policy/#policy-object).

## Example Usage

```hcl
data "okta_policy" "access" {
  name = "Any two factors"
  type = "ACCESS_POLICY",
}

data "okta_app" "example" {
  label = "Example App"
}

resource "okta_app_access_policy_assignment" "assignment" {
  app_id    = data.okta_app.example.id
  policy_id = data.okta_policy.access.id
}
```


## Argument Reference

The following arguments are supported:

- `app_id` - (Required) The application ID; this value is immutable and can not be updated.
- `policy_id` - (Required) The access policy ID.

## Attributes Reference

- `id` - The ID of the resource. This ID is simply the application ID.

## Import

An Okta App's Access Policy Assignment can be imported via its associated Application ID.

```
$ terraform import okta_app_access_policy_assignment.example &#60;app id&#62;
```