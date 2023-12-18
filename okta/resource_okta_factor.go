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
		Description: "Allows you to manage the activation of Okta MFA methods. This resource allows you to manage Okta MFA methods.",
		Schema: map[string]*schema.Schema{
			"provider_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The MFA provider name. Allowed values are `duo`, `fido_u2f`, `fido_webauthn`, `google_otp`, `okta_call`, `okta_otp`, `okta_password`, `okta_push`, `okta_question`, `okta_sms`, `okta_email`, `rsa_token`, `symantec_vip`, `yubikey_token`, or `hotp`.",
				ForceNew:    true,
			},
			"active": {
				Type:        schema.TypeBool,
				Description: "Whether to activate the provider, by default, it is set to `true`.",
				Optional:    true,
				Default:     true,
			},
		},
	}
}

func resourceFactorPut(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	factor, _, err := getAPISupplementFromMetadata(m).GetOrgFactor(ctx, d.Get("provider_id").(string))
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
	factor, resp, err := getAPISupplementFromMetadata(m).GetOrgFactor(ctx, d.Id())
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
	_, resp, err := getAPISupplementFromMetadata(m).DeactivateOrgFactor(ctx, d.Id())
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
		_, _, err = getAPISupplementFromMetadata(m).ActivateOrgFactor(ctx, id)
	} else {
		_, _, err = getAPISupplementFromMetadata(m).DeactivateOrgFactor(ctx, id)
	}
	return err
}

func statusMismatch(d *schema.ResourceData, factor *sdk.OrgFactor) bool {
	status := d.Get("active").(bool)

	// I miss ternary operators
	if factor != nil && factor.Status == statusActive {
		return !status
	}

	return status
}
