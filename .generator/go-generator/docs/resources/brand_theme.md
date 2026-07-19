---
page_title: "okta_brand_theme Resource - terraform-provider-okta"
description: |-
  Manages an Okta Brand Theme resource.
---

# okta_brand_theme (Resource)

Manages an Okta Brand Theme.

## Example Usage

```hcl
resource "okta_brand_theme" "example" {
  brand_id = "<brand-id>"

  # Optional fields
  # email_template_touch_point_variant = "user@example.com"
  # end_user_dashboard_touch_point_variant = "<end_user_dashboard_touch_point_variant>"
  # error_page_touch_point_variant = "<error_page_touch_point_variant>"
  # loading_page_touch_point_variant = "<loading_page_touch_point_variant>"
  # primary_color_contrast_hex = "<primary_color_contrast_hex>"
}
```

## Schema

### Required

- `brand_id` (String) ID of the brand

### Optional

- `email_template_touch_point_variant` (String) Variant for email templates. You can publish a theme for email templates with different combinations of assets. Variants are preset combinations of those assets.
- `end_user_dashboard_touch_point_variant` (String) Variant for the Okta End-User Dashboard. You can publish a theme for end-user dashboard with different combinations of assets. Variants are preset combinations of those assets.
- `error_page_touch_point_variant` (String) Variant for the error page. You can publish a theme for error page with different combinations of assets. Variants are preset combinations of those assets.
- `loading_page_touch_point_variant` (String) Variant for the Okta loading page. You can publish a theme for Okta loading page with different combinations of assets. Variants are preset combinations of those assets.
- `primary_color_contrast_hex` (String) Primary color contrast hex code
- `primary_color_hex` (String) Primary color hex code
- `secondary_color_contrast_hex` (String) Secondary color contrast hex code
- `secondary_color_hex` (String) Secondary color hex code
- `sign_in_page_touch_point_variant` (String) Variant for the Okta sign-in page. You can publish a theme for sign-in page with different combinations of assets. Variants are preset combinations of those assets. > **Note:**  For a non-`OKTA_DEF...

### Read-Only

- `id` (String) The unique identifier for the resource.
- `background_image` (String) BackgroundImage
- `favicon` (String) Favicon
- `logo` (String) Logo

## Import

Import using `{brand_id}/{id}`:

```shell
terraform import okta_brand_theme.example <brand_id> <id>
```
