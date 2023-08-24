---
layout: 'okta'
page_title: 'Okta: okta_brand'
sidebar_current: 'docs-okta-resource-brand'
description: |-
  Gets, updates, an Okta Brand.
---

# okta_brand

This resource allows you to create and configure an Okta
[Brand](https://developer.okta.com/docs/reference/api/brands/#brand-object).

## Example Usage

```hcl
# resource has been imported into current state
# $ terraform import okta_brand.example <brand id>
resource "okta_brand" "example" {
  name = "example
}
```

## Attention
Due to the way brand api works, you will need to first create a brand with only the name, then you can update other attirbutes later. Failure to do so will cause terraform to error out

## Argument Reference

- `name` - (Required) Name of the brand

- `email_domain_id` - (Optional) Email Domain ID tied to this brand

- `locale` - (Optional) The language specified as an IETF BCP 47 language tag

- `agree_to_custom_privacy_policy` - (Optional) Is a required input flag with when changing custom_privacy_url, shouldn't be considered as a readable property

- `custom_privacy_policy_url` - (Optional) Custom privacy policy URL

- `remove_powered_by_okta` - (Optional) Removes "Powered by Okta" from the Okta-hosted sign-in page, and "© 2021 Okta, Inc." from the Okta End-User Dashboard

- `default_app_app_instance_id` - (Optional) Default app app instance id

- `default_app_app_link_name` - (Optional) Default app app link name

- `default_app_classic_application_uri` - (Optional) Default app classic application uri

## Attributes Reference

- `id` - (Read-only) Brand ID

- `is_default` - (Read-only) Is this the default brand

- `links` - (Read-only) Link relations for this object - JSON HAL - Discoverable resources related to the brand

- `brand_id` - (Read-only) Brand ID, used for read (faux-create). Setting `brand_id` to `default` is equivalent to importing the default brand by its ID.
  - `DEPRECATED`: Please stop using this field as it has become noop.

## Import

An Okta Brand can be imported via the ID.

```
$ terraform import okta_brand.example &#60;brand id&#62;
```