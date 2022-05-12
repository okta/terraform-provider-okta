package okta

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

var themeResourceSchema = map[string]*schema.Schema{
	"brand_id": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Brand ID",
	},
	"theme_id": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Theme ID - Note: Okta API for theme only reads and updates therefore the okta_theme resource needs to act as a quasi data source. Do this by setting theme_id.",
	},
	"id": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Brand ID",
	},
	"logo": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "Path to local file",
		DiffSuppressFunc: suppressDuringCreateFunc("theme_id"),
		StateFunc:        localFileStateFunc,
	},
	"logo_url": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Logo URL",
	},
	"favicon": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "Path to local file",
		DiffSuppressFunc: suppressDuringCreateFunc("theme_id"),
		StateFunc:        localFileStateFunc,
	},
	"favicon_url": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Favicon URL",
	},
	"background_image": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "Path to local file",
		DiffSuppressFunc: suppressDuringCreateFunc("theme_id"),
		StateFunc:        localFileStateFunc,
	},
	"background_image_url": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Background image URL",
	},
	"primary_color_hex": {
		Type: schema.TypeString,
		//Required:         true,
		Optional:         true,
		Description:      "Primary color hex code",
		DiffSuppressFunc: suppressDuringCreateFunc("theme_id"),
	},
	"primary_color_contrast_hex": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "Primary color contrast hex code",
		DiffSuppressFunc: suppressDuringCreateFunc("theme_id"),
	},
	"secondary_color_hex": {
		Type: schema.TypeString,
		//Required:         true,
		Optional:         true,
		Description:      "Secondary color hex code",
		DiffSuppressFunc: suppressDuringCreateFunc("theme_id"),
	},
	"secondary_color_contrast_hex": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "Secondary color contrast hex code",
		DiffSuppressFunc: suppressDuringCreateFunc("theme_id"),
	},
	"sign_in_page_touch_point_variant": {
		Type: schema.TypeString,
		//Required:         true,
		Optional:         true,
		Description:      "Variant for the Okta Sign-In Page (`OKTA_DEFAULT`, `BACKGROUND_SECONDARY_COLOR`, `BACKGROUND_IMAGE`)",
		DiffSuppressFunc: suppressDuringCreateFunc("theme_id"),
	},
	"end_user_dashboard_touch_point_variant": {
		Type: schema.TypeString,
		//Required:         true,
		Optional:         true,
		Description:      "Variant for the Okta End-User Dashboard (`OKTA_DEFAULT`, `WHITE_LOGO_BACKGROUND`, `FULL_THEME`, `LOGO_ON_FULL_WHITE_BACKGROUND`)",
		DiffSuppressFunc: suppressDuringCreateFunc("theme_id"),
	},
	"error_page_touch_point_variant": {
		Type: schema.TypeString,
		//Required:         true,
		Optional:         true,
		Description:      "Variant for the error page (`OKTA_DEFAULT`, `BACKGROUND_SECONDARY_COLOR`, `BACKGROUND_IMAGE`)",
		DiffSuppressFunc: suppressDuringCreateFunc("theme_id"),
	},
	"email_template_touch_point_variant": {
		Type: schema.TypeString,
		//Required:         true,
		Optional:         true,
		Description:      "Variant for email templates (`OKTA_DEFAULT`, `FULL_THEME`)",
		DiffSuppressFunc: suppressDuringCreateFunc("theme_id"),
	},
	"links": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Link relations for this object - JSON HAL - Discoverable resources related to the email template",
	},
}

var themesDataSourceSchema = map[string]*schema.Schema{
	"brand_id": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Brand ID",
	},
	"themes": {
		Type:        schema.TypeSet,
		Computed:    true,
		Description: "List of `okta_them` belonging to the brand in the organization",
		Elem: &schema.Resource{
			Schema: themeDataSourceSchema,
		},
	},
}

var themeDataSourceSchema = map[string]*schema.Schema{
	"id": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The ID of the theme",
	},
	"logo_url": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Logo URL",
	},
	"favicon_url": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Favicon URL",
	},
	"background_image_url": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Background image URL",
	},
	"primary_color_hex": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Primary color hex code",
	},
	"primary_color_contrast_hex": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Primary color contrast hex code",
	},
	"secondary_color_hex": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Secondary color hex code",
	},
	"secondary_color_contrast_hex": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Secondary color contrast hex code",
	},
	"sign_in_page_touch_point_variant": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Variant for the Okta Sign-In Page (`OKTA_DEFAULT`, `BACKGROUND_SECONDARY_COLOR`, `BACKGROUND_IMAGE`)",
	},
	"end_user_dashboard_touch_point_variant": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Variant for the Okta End-User Dashboard (`OKTA_DEFAULT`, `WHITE_LOGO_BACKGROUND`, `FULL_THEME`, `LOGO_ON_FULL_WHITE_BACKGROUND`)",
	},
	"error_page_touch_point_variant": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Variant for the error page (`OKTA_DEFAULT`, `BACKGROUND_SECONDARY_COLOR`, `BACKGROUND_IMAGE`)",
	},
	"email_template_touch_point_variant": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Variant for email templates (`OKTA_DEFAULT`, `FULL_THEME`)",
	},
	"links": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Link relations for this object - JSON HAL - Discoverable resources related to the email template",
	},
}

func flattenTheme(brandID string, theme *okta.ThemeResponse) map[string]interface{} {
	attrs := map[string]interface{}{}

	attrs["id"] = theme.Id
	if brandID != "" {
		attrs["brand_id"] = brandID
	}

	attrs["logo_url"] = theme.Logo
	attrs["favicon_url"] = theme.Favicon
	attrs["background_image_url"] = theme.BackgroundImage
	attrs["primary_color_hex"] = theme.PrimaryColorHex
	attrs["primary_color_contrast_hex"] = theme.PrimaryColorContrastHex
	attrs["secondary_color_hex"] = theme.SecondaryColorHex
	attrs["secondary_color_contrast_hex"] = theme.SecondaryColorContrastHex
	attrs["sign_in_page_touch_point_variant"] = theme.SignInPageTouchPointVariant
	attrs["end_user_dashboard_touch_point_variant"] = theme.EndUserDashboardTouchPointVariant
	attrs["error_page_touch_point_variant"] = theme.ErrorPageTouchPointVariant
	attrs["email_template_touch_point_variant"] = theme.EmailTemplateTouchPointVariant

	links, _ := json.Marshal(theme.Links)
	attrs["links"] = string(links)

	return attrs
}
