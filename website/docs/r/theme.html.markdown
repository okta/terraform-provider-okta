---
layout: 'okta'
page_title: 'Okta: okta_theme'
sidebar_current: 'docs-okta-resource-theme'
description: |-
  Gets, updates, a single Theme of a Brand of an Okta Organization.
---

# okta_theme

This resource allows you to get and update an Okta
[Theme](https://developer.okta.com/docs/reference/api/brands/#theme-object).

The Okta Management API does not have a true Create or Delete for a theme. Therefore, the theme resource must be imported
first into the terraform state before updates can be applied to the theme.

## Example Usage

```hcl
data "okta_brands" "test" {
}

# resource has been imported into current state:
# $ terraform import okta_theme.example <theme id>
resource "okta_theme" "example" {
    brand_id = tolist(data.okta_brands.test.brands)[0].id
    logo                                   = "path/to/logo.png"
    favicon                                = "path/to/favicon.png"
    background_image                       = "path/to/background.png"
    primary_color_hex                      = "#1662dd"
    secondary_color_hex                    = "#ebebed"
    sign_in_page_touch_point_variant       = "OKTA_DEFAULT"
    end_user_dashboard_touch_point_variant = "OKTA_DEFAULT"
    error_page_touch_point_variant         = "OKTA_DEFAULT"
    email_template_touch_point_variant     = "OKTA_DEFAULT"
}
```

## Arguments Reference

- `brand_id` - (Required) Brand ID
- `theme_id` - (Optional) Theme ID, used for read (faux-create)

## Attributes Reference

Related Okta API [Theme Response Object](https://developer.okta.com/docs/reference/api/brands/#theme-response-object)

- `id` - (Read-Only) Theme URL
- `logo` - (Optional) Local path to logo file. Setting the value to the blank string `""` will delete the logo on the theme at Okta but will not delete the local file.
- `logo_url` - (Read-Only) Logo URL
- `favicon` - (Optional) Local path to favicon file. Setting the value to the blank string `""` will delete the favicon on the theme at Okta but will not delete the local file.
- `favicon_url` - (Read-Only) Favicon URL
- `background_image` - (Optional) Local path to background image file. Setting the value to the blank string `""` will delete the favicon on the theme at Okta but will not delete the local file.
- `background_image_url` - (Read-Only) Background image URL
- `primary_color_hex` - (Required) Primary color hex code
- `primary_color_contrast_hex` - (Optional) Primary color contrast hex code
- `secondary_color_hex` - (Required) Secondary color hex code
- `secondary_color_contrast_hex` (Optional) Secondary color contrast hex code
- `sign_in_page_touch_point_variant` - (Required) Variant for the Okta Sign-In Page. Valid values: (`OKTA_DEFAULT`, `BACKGROUND_SECONDARY_COLOR`, `BACKGROUND_IMAGE`)
- `end_user_dashboard_touch_point_variant` - (Required) Variant for the Okta End-User Dashboard. Valid values: (`OKTA_DEFAULT`, `WHITE_LOGO_BACKGROUND`, `FULL_THEME`, `LOGO_ON_FULL_WHITE_BACKGROUND`)
- `error_page_touch_point_variant` - (Required) Variant for the error page. Valid values: (`OKTA_DEFAULT`, `BACKGROUND_SECONDARY_COLOR`, `BACKGROUND_IMAGE`)
- `email_template_touch_point_variant` - (Required) Variant for email templates. Valid values: (`OKTA_DEFAULT`, `FULL_THEME`)
- `links` - Link relations for this object - JSON HAL - (Read-Only) Discoverable resources related to the brand

[Variants for the Okta Sign-In Page](https://developer.okta.com/docs/reference/api/brands/#variants-for-the-okta-sign-in-page):

| Enum Value  |  Description  |
| ----------- | ------------- |
| OKTA_DEFAULT | Use the Okta logo, Okta favicon with no background image, and the Okta colors on the Okta Sign-In Page. |
| BACKGROUND_SECONDARY_COLOR | Use the logo and favicon from Theme with the secondaryColorHex as the background color for the Okta Sign-In Page. |
| BACKGROUND_IMAGE | Use the logo, favicon, and background image from Theme. |

[Variants for the Okta End-User Dashboard](https://developer.okta.com/docs/reference/api/brands/#variants-for-the-okta-end-user-dashboard):

| Enum Value  | Description   |
| ----------- | ------------- |
| OKTA_DEFAULT | Use the Okta logo and Okta favicon with a white background color for the logo and the side navigation bar background color. |
| WHITE_LOGO_BACKGROUND | Use the logo from Theme with a white background color for the logo, use favicon from Theme, and use primaryColorHex for the side navigation bar background color. |
| FULL_THEME | Use the logo from Theme, primaryColorHex for the logo and the side navigation bar background color, and use favicon from Theme |
| LOGO_ON_FULL_WHITE_BACKGROUND | Use the logo from Theme, white background color for the logo and the side navigation bar background color, and use favicon from Theme |

## Import

An Okta Brand can be imported via the ID.

```
$ terraform import okta_theme.example &#60;brand id&#62;/&#60;theme id&#62;
```
