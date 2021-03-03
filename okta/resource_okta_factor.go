package okta

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

// Predefined second authentication factors. They must be activated in order to use them in MFA policies.
// This is not your standard resource as each factor provider is predefined and the create function simply puts it in
// terraform state and activates it. Currently the API is in Beta and it only allows lifecycle interactions, and
// no ability to configure them but the resource was built with future expansion in mind. Also keep in mind this
// is an account level resource.
func resourceFactor() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFactorPut,
		ReadContext:   resourceFactorRead,
		UpdateContext: resourceFactorPut,
		DeleteContext: resourceFactorDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"provider_id": {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: stringInSlice([]string{
					sdk.DuoFactor,
					sdk.FidoU2fFactor,
					sdk.FidoWebauthnFactor,
					sdk.GoogleOtpFactor,
					sdk.OktaCallFactor,
					sdk.OktaOtpFactor,
					sdk.OktaPasswordFactor,
					sdk.OktaPushFactor,
					sdk.OktaQuestionFactor,
					sdk.OktaSmsFactor,
					sdk.OktaEmailFactor,
					sdk.RsaTokenFactor,
					sdk.SymantecVipFactor,
					sdk.YubikeyTokenFactor,
					sdk.HotpFactor,
				}),
				Description: "Factor provider ID",
				ForceNew:    true,
			},
			"active": {
				Type:        schema.TypeBool,
				Description: "Is this provider active?",
				Optional:    true,
				Default:     true,
			},
		},
	}
}

func resourceFactorPut(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	factor, _, err := getSupplementFromMetadata(m).GetFactor(ctx, d.Get("provider_id").(string))
	if err != nil {
		return diag.Errorf("failed to find factor: %v", err)
	}
	// To avoid API errors we check downstream status
	if statusMismatch(d, factor) {
		err := activateFactor(ctx, d, m)
		if err != nil {
			return diag.Errorf("failed to activate factor: %v", err)
		}
	}
	d.SetId(d.Get("provider_id").(string))
	return resourceFactorRead(ctx, d, m)
}

func resourceFactorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	factor, resp, err := getSupplementFromMetadata(m).GetFactor(ctx, d.Get("provider_id").(string))
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to find factor: %v", err)
	}
	if factor == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("active", factor.Status == statusActive)
	_ = d.Set("provider_id", factor.Id)
	return nil
}

func resourceFactorDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if !d.Get("active").(bool) {
		return nil
	}
	_, resp, err := getSupplementFromMetadata(m).DeactivateFactor(ctx, d.Id())
	// http.StatusBadRequest means that factor can not be deactivated
	if resp != nil && resp.StatusCode == http.StatusBadRequest {
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to deactivate '%s' factor: %v", d.Id(), err)
	}
	return nil
}

func activateFactor(ctx context.Context, d *schema.ResourceData, m interface{}) error {
	var err error
	id := d.Get("provider_id").(string)
	if d.Get("active").(bool) {
		_, _, err = getSupplementFromMetadata(m).ActivateFactor(ctx, id)
	} else {
		_, _, err = getSupplementFromMetadata(m).DeactivateFactor(ctx, id)
	}
	return err
}

func statusMismatch(d *schema.ResourceData, factor *sdk.Factor) bool {
	status := d.Get("active").(bool)

	// I miss ternary operators
	if factor != nil && factor.Status == statusActive {
		return !status
	}

	return status
}
