package idaas

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v4/okta"
	"github.com/okta/terraform-provider-okta/okta/utils"
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
		DiffSuppressFunc: utils.SuppressDuringCreateFunc("theme_id"),
		StateFunc:        utils.LocalFileStateFunc,
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
		DiffSuppressFunc: utils.SuppressDuringCreateFunc("theme_id"),
		StateFunc:        utils.LocalFileStateFunc,
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
		DiffSuppressFunc: utils.SuppressDuringCreateFunc("theme_id"),
		StateFunc:        utils.LocalFileStateFunc,
	},
	"background_image_url": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Background image URL",
	},
	"primary_color_hex": {
		Type: schema.TypeString,
		// Required:         true,
		Optional:         true,
		Description:      "Primary color hex code",
		DiffSuppressFunc: utils.SuppressDuringCreateFunc("theme_id"),
	},
	"primary_color_contrast_hex": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "Primary color contrast hex code",
		DiffSuppressFunc: utils.SuppressDuringCreateFunc("theme_id"),
	},
	"secondary_color_hex": {
		Type: schema.TypeString,
		// Required:         true,
		Optional:         true,
		Description:      "Secondary color hex code",
		DiffSuppressFunc: utils.SuppressDuringCreateFunc("theme_id"),
	},
	"secondary_color_contrast_hex": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "Secondary color contrast hex code",
		DiffSuppressFunc: utils.SuppressDuringCreateFunc("theme_id"),
	},
	"sign_in_page_touch_point_variant": {
		Type: schema.TypeString,
		// Required:         true,
		Optional:         true,
		Description:      "Variant for the Okta Sign-In Page (`OKTA_DEFAULT`, `BACKGROUND_SECONDARY_COLOR`, `BACKGROUND_IMAGE`)",
		DiffSuppressFunc: utils.SuppressDuringCreateFunc("theme_id"),
	},
	"end_user_dashboard_touch_point_variant": {
		Type: schema.TypeString,
		// Required:         true,
		Optional:         true,
		Description:      "Variant for the Okta End-User Dashboard (`OKTA_DEFAULT`, `WHITE_LOGO_BACKGROUND`, `FULL_THEME`, `LOGO_ON_FULL_WHITE_BACKGROUND`)",
		DiffSuppressFunc: utils.SuppressDuringCreateFunc("theme_id"),
	},
	"error_page_touch_point_variant": {
		Type: schema.TypeString,
		// Required:         true,
		Optional:         true,
		Description:      "Variant for the error page (`OKTA_DEFAULT`, `BACKGROUND_SECONDARY_COLOR`, `BACKGROUND_IMAGE`)",
		DiffSuppressFunc: utils.SuppressDuringCreateFunc("theme_id"),
	},
	"email_template_touch_point_variant": {
		Type: schema.TypeString,
		// Required:         true,
		Optional:         true,
		Description:      "Variant for email templates (`OKTA_DEFAULT`, `FULL_THEME`)",
		DiffSuppressFunc: utils.SuppressDuringCreateFunc("theme_id"),
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

	attrs["id"] = theme.GetId()
	if brandID != "" {
		attrs["brand_id"] = brandID
	}

	attrs["logo_url"] = theme.GetLogo()
	attrs["favicon_url"] = theme.GetFavicon()
	attrs["background_image_url"] = theme.GetBackgroundImage()
	attrs["primary_color_hex"] = theme.GetPrimaryColorHex()
	attrs["primary_color_contrast_hex"] = theme.GetPrimaryColorContrastHex()
	attrs["secondary_color_hex"] = theme.GetSecondaryColorHex()
	attrs["secondary_color_contrast_hex"] = theme.GetSecondaryColorContrastHex()
	attrs["sign_in_page_touch_point_variant"] = string(theme.GetSignInPageTouchPointVariant())
	attrs["end_user_dashboard_touch_point_variant"] = string(theme.GetEndUserDashboardTouchPointVariant())
	attrs["error_page_touch_point_variant"] = string(theme.GetErrorPageTouchPointVariant())
	attrs["email_template_touch_point_variant"] = string(theme.GetEmailTemplateTouchPointVariant())

	links, _ := json.Marshal(theme.GetLinks())
	attrs["links"] = string(links)

	return attrs
}
