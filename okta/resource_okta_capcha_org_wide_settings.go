package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourceCaptchaOrgWideSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCaptchaOrgWideSettingsCreate,
		ReadContext:   resourceCaptchaOrgWideSettingsRead,
		UpdateContext: resourceCaptchaOrgWideSettingsUpdate,
		DeleteContext: resourceCaptchaOrgWideSettingsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"captcha_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the CAPTCHA",
			},
			"enabled_for": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: elemInSlice([]string{"SSR", "SSPR", "SIGN_IN"}),
				},
				Description:  "Set of pages that have CAPTCHA enabled",
				RequiredWith: []string{"captcha_id"},
			},
		},
	}
}

func resourceCaptchaOrgWideSettingsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if isClassicOrg(m) {
		return resourceOIEOnlyFeatureError(captchaOrgWideSettings)
	}

	captcha, _, err := getSupplementFromMetadata(m).UpdateOrgWideCaptchaSettings(ctx, buildCaptchaOrgWideSettings(d))
	if err != nil {
		return diag.Errorf("failed to set org-wide CAPTCHA settings: %v", err)
	}
	_ = d.Set("captcha_id", captcha.CaptchaId)
	_ = d.Set("enabled_for", convertStringSliceToSetNullable(captcha.EnabledPages))
	d.SetId("org_wide_captcha")
	return nil
}

func resourceCaptchaOrgWideSettingsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if isClassicOrg(m) {
		return resourceOIEOnlyFeatureError(captchaOrgWideSettings)
	}

	captcha, _, err := getSupplementFromMetadata(m).GetOrgWideCaptchaSettings(ctx)
	if err != nil {
		return diag.Errorf("failed to get org-wide CAPTCHA settings: %v", err)
	}
	if captcha == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("captcha_id", captcha.CaptchaId)
	_ = d.Set("enabled_for", convertStringSliceToSetNullable(captcha.EnabledPages))
	d.SetId("org_wide_captcha")
	return nil
}

func resourceCaptchaOrgWideSettingsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if isClassicOrg(m) {
		return resourceOIEOnlyFeatureError(captchaOrgWideSettings)
	}

	captcha, _, err := getSupplementFromMetadata(m).UpdateOrgWideCaptchaSettings(ctx, buildCaptchaOrgWideSettings(d))
	if err != nil {
		return diag.Errorf("failed to update org-wide CAPTCHA settings: %v", err)
	}
	_ = d.Set("captcha_id", captcha.CaptchaId)
	_ = d.Set("enabled_for", convertStringSliceToSetNullable(captcha.EnabledPages))
	return nil
}

func resourceCaptchaOrgWideSettingsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if isClassicOrg(m) {
		return resourceOIEOnlyFeatureError(captchaOrgWideSettings)
	}

	_, err := getSupplementFromMetadata(m).DeleteOrgWideCaptchaSettings(ctx)
	if err != nil {
		return diag.Errorf("failed to delete org-wide CAPTCHA settings: %v", err)
	}
	return nil
}

func buildCaptchaOrgWideSettings(d *schema.ResourceData) sdk.OrgWideCaptchaSettings {
	s := sdk.OrgWideCaptchaSettings{
		EnabledPages: convertInterfaceToStringSet(d.Get("enabled_for")),
	}
	captchID, ok := d.GetOk("captcha_id")
	if ok {
		id := captchID.(string)
		s.CaptchaId = &id
	}
	return s
}
