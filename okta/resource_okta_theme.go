package okta

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v3/okta"
)

func resourceTheme() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceThemeCreate,
		ReadContext:   resourceThemeRead,
		UpdateContext: resourceThemeUpdate,
		DeleteContext: resourceThemeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceThemeImportStateContext,
		},
		Schema: themeResourceSchema,
	}
}

func resourceThemeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	bid, ok := d.GetOk("brand_id")
	if !ok {
		return diag.Errorf("brand_id required to create theme")
	}
	brandID := bid.(string)

	var themeID string
	if d.Id() != "" {
		themeID = d.Id()
	}
	if themeID == "" {
		if tid, ok := d.GetOk("theme_id"); ok {
			themeID = tid.(string)
		}
	}
	if themeID == "" {
		return diag.Errorf("theme_id required to create theme")
	}

	theme, _, err := getOktaV3ClientFromMetadata(m).CustomizationApi.GetBrandTheme(ctx, brandID, themeID).Execute()
	if err != nil {
		return diag.Errorf("failed to get theme: %v", err)
	}

	d.SetId(theme.GetId())
	rawMap := flattenTheme(brandID, theme)
	err = setNonPrimitives(d, rawMap)
	if err != nil {
		return diag.Errorf("failed to set theme properties: %v", err)
	}

	return nil
}

func resourceThemeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	logger(m).Info("reading theme", "id", d.Id())

	bid, ok := d.GetOk("brand_id")
	if !ok {
		return diag.Errorf("brand_id required to import theme")
	}
	brandID := bid.(string)

	theme, _, err := getOktaV3ClientFromMetadata(m).CustomizationApi.GetBrandTheme(ctx, brandID, d.Id()).Execute()
	if err != nil {
		return diag.Errorf("failed to get theme: %v", err)
	}

	rawMap := flattenTheme(brandID, theme)
	err = setNonPrimitives(d, rawMap)
	if err != nil {
		return diag.Errorf("failed to set theme properties: %v", err)
	}

	return nil
}

func resourceThemeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	logger(m).Info("updating theme", "id", d.Id())

	bid, ok := d.GetOk("brand_id")
	if !ok {
		return diag.Errorf("brand_id required to update theme")
	}
	brandID := bid.(string)

	// peform delete/upload on the logo/favicon/background_image first so any
	// errors there will interrupt apply on the theme itself
	if d.HasChange("logo") {
		err := handleThemeLogo(ctx, d, m, brandID, d.Id())
		if err != nil {
			return diag.Errorf("failed to handle logo for theme: %v", err)
		}
	}
	if d.HasChange("favicon") {
		err := handleThemeFavicon(ctx, d, m, brandID, d.Id())
		if err != nil {
			return diag.Errorf("failed to handle favicon for theme: %v", err)
		}
	}
	if d.HasChange("background_image") {
		err := handleThemeBackgroundImage(ctx, d, m, brandID, d.Id())
		if err != nil {
			return diag.Errorf("failed to handle background_image for theme: %v", err)
		}
	}

	theme := okta.Theme{}

	if val, ok := d.GetOk("primary_color_hex"); ok {
		theme.PrimaryColorHex = stringPtr(val.(string))
	}

	if val, ok := d.GetOk("primary_color_contrast_hex"); ok {
		theme.PrimaryColorContrastHex = stringPtr(val.(string))
	}

	if val, ok := d.GetOk("secondary_color_hex"); ok {
		theme.SecondaryColorHex = stringPtr(val.(string))
	}

	if val, ok := d.GetOk("secondary_color_contrast_hex"); ok {
		theme.SecondaryColorContrastHex = stringPtr(val.(string))
	}

	if val, ok := d.GetOk("sign_in_page_touch_point_variant"); ok {
		if siptpv, err := okta.NewSignInPageTouchPointVariantFromValue(val.(string)); err == nil {
			theme.SignInPageTouchPointVariant = siptpv
		}
	}

	if val, ok := d.GetOk("end_user_dashboard_touch_point_variant"); ok {
		if eudtpv, err := okta.NewEndUserDashboardTouchPointVariantFromValue(val.(string)); err == nil {
			theme.EndUserDashboardTouchPointVariant = eudtpv
		}
	}

	if val, ok := d.GetOk("error_page_touch_point_variant"); ok {
		if eptpv, err := okta.NewErrorPageTouchPointVariantFromValue(val.(string)); err == nil {
			theme.ErrorPageTouchPointVariant = eptpv
		}
	}

	if val, ok := d.GetOk("email_template_touch_point_variant"); ok {
		if ettpv, err := okta.NewEmailTemplateTouchPointVariantFromValue(val.(string)); err == nil {
			theme.EmailTemplateTouchPointVariant = ettpv
		}
	}

	themeResp, _, err := getOktaV3ClientFromMetadata(m).CustomizationApi.ReplaceBrandTheme(ctx, brandID, d.Id()).Theme(theme).Execute()
	if err != nil {
		return diag.Errorf("failed to update theme: %v", err)
	}

	rawMap := flattenTheme(brandID, themeResp)
	err = setNonPrimitives(d, rawMap)
	if err != nil {
		return diag.Errorf("failed to set theme properties: %v", err)
	}

	return nil
}

func resourceThemeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// fake delete
	d.SetId("")
	return nil
}

func resourceThemeImportStateContext(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid resource import specifier, expecting the following format: <brand_id>/<theme_id>")
	}
	brandID := parts[0]
	themeID := parts[1]

	theme, _, err := getOktaV3ClientFromMetadata(m).CustomizationApi.GetBrandTheme(ctx, brandID, themeID).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get theme: %v", err)
	}

	d.SetId(theme.GetId())
	rawMap := flattenTheme(brandID, theme)
	err = setNonPrimitives(d, rawMap)
	if err != nil {
		return nil, fmt.Errorf("failed to set theme properties: %v", err)
	}

	return []*schema.ResourceData{d}, nil
}

func handleThemeLogo(ctx context.Context, d *schema.ResourceData, m interface{}, brandID, themeID string) error {
	_, newPath := d.GetChange("logo")
	if newPath == "" {
		_, err := getOktaV3ClientFromMetadata(m).CustomizationApi.DeleteBrandThemeLogo(ctx, brandID, themeID).Execute()
		return err
	}
	fo, err := os.Open(newPath.(string))
	if err != nil {
		return err
	}
	defer fo.Close()
	_, _, err = getOktaV3ClientFromMetadata(m).CustomizationApi.UploadBrandThemeLogo(ctx, brandID, themeID).File(fo).Execute()
	return err
}

func handleThemeFavicon(ctx context.Context, d *schema.ResourceData, m interface{}, brandID, themeID string) error {
	_, newPath := d.GetChange("favicon")
	if newPath == "" {
		_, err := getOktaV3ClientFromMetadata(m).CustomizationApi.DeleteBrandThemeFavicon(ctx, brandID, themeID).Execute()
		return err
	}
	fo, err := os.Open(newPath.(string))
	if err != nil {
		return err
	}
	defer fo.Close()
	_, _, err = getOktaV3ClientFromMetadata(m).CustomizationApi.UploadBrandThemeFavicon(ctx, brandID, themeID).File(fo).Execute()
	return err
}

func handleThemeBackgroundImage(ctx context.Context, d *schema.ResourceData, m interface{}, brandID, themeID string) error {
	_, newPath := d.GetChange("background_image")
	if newPath == "" {
		_, err := getOktaV3ClientFromMetadata(m).CustomizationApi.DeleteBrandThemeBackgroundImage(ctx, brandID, themeID).Execute()
		return err
	}
	fo, err := os.Open(newPath.(string))
	if err != nil {
		return err
	}
	defer fo.Close()
	_, _, err = getOktaV3ClientFromMetadata(m).CustomizationApi.UploadBrandThemeBackgroundImage(ctx, brandID, themeID).File(fo).Execute()
	return err
}
