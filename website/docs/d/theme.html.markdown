---
layout: 'okta'
page_title: 'Okta: okta_theme'
sidebar_current: 'docs-okta-datasource-theme'
description: |-
  Get a single Theme of a Brand of an Okta Organization.
---

# okta_theme

Use this data source to retrieve a 
[Theme](https://developer.okta.com/docs/reference/api/brands/#theme-response-object)
of a brand for an Okta orgnanization.

## Example Usage

```hcl
data "okta_brands" "test" {
}

data "okta_themes" "test" {
  brand_id = tolist(data.okta_brands.test.brands)[0].id
}

data "okta_theme" "test" {
  brand_id = tolist(data.okta_brands.test.brands)[0].id
  theme_id = tolist(data.okta_themes.test.themes)[0].id
}
```

## Arguments Reference

- `brand_id` - (Required) Brand ID
- `theme_id` - (Required) Theme ID

## Attributes Reference

Related Okta API [Theme Response Object](https://developer.okta.com/docs/reference/api/brands/#theme-response-object)

- `id` - Theme URL
- `logo` - Logo URL
- `favicon` - Favicon URL
- `background_image` - Background image URL
- `primary_color_hex` - Primary color hex code
- `primary_color_contrast_hex` - Primary color contrast hex code
- `secondary_color_hex` - Secondary color hex code
- `secondary_color_contrast_hex` Secondary color contrast hex code
- `sign_in_page_touch_point_variant` - (Enum) Variant for the Okta Sign-In Page (`OKTA_DEFAULT`, `BACKGROUND_SECONDARY_COLOR`, `BACKGROUND_IMAGE`)
- `end_user_dashboard_touch_point_variant` - (Enum) Variant for the Okta End-User Dashboard (`OKTA_DEFAULT`, `WHITE_LOGO_BACKGROUND`, `FULL_THEME`, `LOGO_ON_FULL_WHITE_BACKGROUND`)
- `error_page_touch_point_variant` - (Enum) Variant for the error page (`OKTA_DEFAULT`, `BACKGROUND_SECONDARY_COLOR`, `BACKGROUND_IMAGE`)
- `email_template_touch_point_variant` - (Enum) Variant for email templates (`OKTA_DEFAULT`, `FULL_THEME`)
- `links` - Link relations for this object - JSON HAL - Discoverable resources related to the brand
