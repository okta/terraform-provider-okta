---
page_title: "okta_brand_email_template_customization Resource - terraform-provider-okta"
description: |-
  Manages an Okta Brand Email Template Customization resource.
---

# okta_brand_email_template_customization (Resource)

Manages an Okta Brand Email Template Customization.

## Example Usage

```hcl
resource "okta_brand_email_template_customization" "example" {
  brand_id = "<brand-id>"
  template_name = "Example Template Name"
  body = "<body>"
  language = "<language>"
  subject = "<subject>"

  # Optional fields
  # is_default = true
}
```

## Schema

### Required

- `brand_id` (String) ID of the brand
- `template_name` (String) Name of the email template
- `body` (String) The HTML body of the email. May contain [variable references](https://velocity.apache.org/engine/1.7/user-guide.html#references).  <x-lifecycle class='ea'></x-lifecycle> Not required if Custom lang...
- `language` (String) The language specified as an [IETF BCP 47 language tag](https://datatracker.ietf.org/doc/html/rfc5646)
- `subject` (String) The email subject. May contain [variable references](https://velocity.apache.org/engine/1.7/user-guide.html#references).  <x-lifecycle class='ea'></x-lifecycle> Not required if Custom languages for...

### Optional

- `is_default` (Boolean) Whether this is the default customization for the email template. Each customized email template must have exactly one default customization. Defaults to `true` for the first customization and `fal...

### Read-Only

- `id` (String) The unique identifier for the resource.
- `created` (String) The UTC time at which this email customization was created.
- `last_updated` (String) The UTC time at which this email customization was last updated.

## Import

Import using `{brand_id}/{template_name}/{id}`:

```shell
terraform import okta_brand_email_template_customization.example <brand_id> <template_name> <id>
```
