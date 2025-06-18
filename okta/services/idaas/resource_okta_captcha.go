package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourceCaptcha() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCaptchaCreate,
		ReadContext:   resourceCaptchaRead,
		UpdateContext: resourceCaptchaUpdate,
		DeleteContext: resourceCaptchaDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: `Creates different types of captcha.
		~> **WARNING:** This feature is only available as a part of the Identity Engine. [Contact support](mailto:dev-inquiries@okta.com) for further information.
		This resource allows you to create and configure a CAPTCHA.`,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the CAPTCHA",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Type of the captcha. Valid values: `HCAPTCHA`, `RECAPTCHA_V2`",
			},
			"site_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Site key issued from the CAPTCHA vendor to render a CAPTCHA on a page",
			},
			"secret_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Secret key issued from the CAPTCHA vendor to perform server-side validation for a CAPTCHA token",
			},
		},
	}
}

func resourceCaptchaCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if providerIsClassicOrg(ctx, meta) {
		return resourceOIEOnlyFeatureError(resources.OktaIDaaSCaptcha)
	}

	captcha, _, err := getAPISupplementFromMetadata(meta).CreateCaptcha(ctx, buildCaptcha(d))
	if err != nil {
		return diag.Errorf("failed to create CAPTCHA: %v", err)
	}
	d.SetId(captcha.Id)
	return nil
}

func resourceCaptchaRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if providerIsClassicOrg(ctx, meta) {
		return resourceOIEOnlyFeatureError(resources.OktaIDaaSCaptcha)
	}

	captcha, resp, err := getAPISupplementFromMetadata(meta).GetCaptcha(ctx, d.Id())
	if err := utils.SuppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to find CAPTCHA: %v", err)
	}
	if captcha == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("name", captcha.Name)
	_ = d.Set("type", captcha.Type)
	_ = d.Set("site_key", captcha.SiteKey)
	return nil
}

func resourceCaptchaUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if providerIsClassicOrg(ctx, meta) {
		return resourceOIEOnlyFeatureError(resources.OktaIDaaSCaptcha)
	}

	_, _, err := getAPISupplementFromMetadata(meta).UpdateCaptcha(ctx, d.Id(), buildCaptcha(d))
	if err != nil {
		return diag.Errorf("failed to update CAPTCHA: %v", err)
	}
	return nil
}

func resourceCaptchaDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if providerIsClassicOrg(ctx, meta) {
		return resourceOIEOnlyFeatureError(resources.OktaIDaaSCaptcha)
	}

	logger(meta).Info("deleting Captcha", "name", d.Get("name").(string))
	_, err := getAPISupplementFromMetadata(meta).DeleteCaptcha(ctx, d.Id())
	if err != nil {
		return diag.Errorf("failed to delete CAPTCHA: %v", err)
	}
	return nil
}

func buildCaptcha(d *schema.ResourceData) sdk.Captcha {
	return sdk.Captcha{
		Name:      d.Get("name").(string),
		SiteKey:   d.Get("site_key").(string),
		SecretKey: d.Get("secret_key").(string),
		Type:      d.Get("type").(string),
	}
}
