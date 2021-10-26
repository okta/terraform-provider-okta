---
layout: 'okta'
page_title: 'Okta: okta_org_configuration'
sidebar_current: 'docs-okta-resource-okta-admin-role-targets'
description: |-
  Manages org settings, logo, support and communication.
---

# okta_org_configuration

This resource allows you manage org settings, logo, support and communication options.

~> **IMPORTANT NOTE:** You must specify all Org Setting properties when you update an org's profile. Any property not specified in the script will be deleted.

## Example Usage

```hcl
resource "okta_org_configuration" "example" {
  company_name = "Umbrella Corporation"
  website      = "https://terraform.io"
}
```

## Argument Reference

`company_name` - (Required) Name of the org.

`website` - (Optional) The org's website.

`phone_number` - (Optional) Phone number of org.

`end_user_support_help_url` - (Optional) Support link of org.

`support_phone_number` - (Optional) Support help phone of org.

`address_1` - (Optional) Primary address of org.

`address_2` - (Optional) Secondary address of org.

`city` - (Optional) City of org.

`state` - (Optional) State of org.

`country` - (Optional) County of org.

`postal_code` - (Optional) Postal code of org.

`logo` - (Optional) Logo of org. The file must be in PNG, JPG, or GIF format and less than 1 MB in size. 
For best results use landscape orientation, a transparent background, and a minimum size of 420px by 120px to prevent upscaling.

`billing_contact_user` - (Optional) User ID representing the billing contact

`technical_contact_user` - (Optional) User ID representing the technical contact.

`opt_out_communication_emails` - (Optional) Indicates whether the org's users receive Okta Communication emails.

## Attributes Reference

`id` - ID of org.

`expires_at` - Expiration of org.

`subdomain` - Subdomain of org.

## Import

Okta Org Configuration can be imported even without specifying the Org ID.

```
$ terraform import okta_org_configuration.example _
```
