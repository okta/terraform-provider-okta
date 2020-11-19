package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
)

// Predefined second authentication factors. They must be activated in order to use them in MFA policies.
// This is not your standard resource as each factor provider is predefined and the create function simply puts it in
// terraform state and activates it. Currently the API is in Beta and it only allows lifecycle interactions, and
// no ability to configure them but the resource was built with future expansion in mind. Also keep in mind this
// is an account level resource.
func resourceFactor() *schema.Resource {
	return &schema.Resource{
		Create: resourceFactorPut,
		Read:   resourceFactorRead,
		Update: resourceFactorPut,
		Exists: resourceFactorExists,
		Delete: resourceFactorDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"provider_id": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						sdk.DuoFactor,
						sdk.FidoU2fFactor,
						sdk.FidoWebauthnFactor,
						sdk.GoogleOtpFactor,
						sdk.OktaCallFactor,
						sdk.OktaOtpFactor,
						sdk.OktaPushFactor,
						sdk.OktaQuestionFactor,
						sdk.OktaSmsFactor,
						sdk.RsaTokenFactor,
						sdk.SymantecVipFactor,
						sdk.YubikeyTokenFactor,
					},
					false,
				),
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

func resourceFactorExists(d *schema.ResourceData, m interface{}) (bool, error) {
	factor, err := findFactor(d, m)
	return err == nil && factor != nil, err
}

func resourceFactorDelete(d *schema.ResourceData, m interface{}) error {
	var err error
	if d.Get("active").(bool) {
		_, _, err = getSupplementFromMetadata(m).DeactivateFactor(context.Background(), d.Id())
	}
	return err
}

func resourceFactorRead(d *schema.ResourceData, m interface{}) error {
	factor, err := findFactor(d, m)
	if err != nil {
		return err
	}
	if factor == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("active", factor.Status == statusActive)
	_ = d.Set("provider_id", factor.Id)
	return nil
}

func resourceFactorPut(d *schema.ResourceData, m interface{}) error {
	factor, err := findFactor(d, m)
	if err != nil {
		return err
	}

	// To avoid API errors we check downstream status
	if statusMismatch(d, factor) {
		err := activateFactor(d, m)
		if err != nil {
			return err
		}
	}
	d.SetId(d.Get("provider_id").(string))

	return resourceFactorRead(d, m)
}

func activateFactor(d *schema.ResourceData, m interface{}) error {
	var err error
	id := d.Get("provider_id").(string)
	if d.Get("active").(bool) {
		_, _, err = getSupplementFromMetadata(m).ActivateFactor(context.Background(), id)
	} else {
		_, _, err = getSupplementFromMetadata(m).DeactivateFactor(context.Background(), id)
	}
	return err
}

// This API is in Beta hence the inability to do a single get. I must list then find.
// Fear is clearly not a factor for me.
func findFactor(d *schema.ResourceData, m interface{}) (*sdk.Factor, error) {
	factorList, _, err := getSupplementFromMetadata(m).ListFactors(context.Background())
	if err != nil {
		return nil, err
	}
	id := d.Get("provider_id").(string)

	for _, f := range factorList {
		if f.Id == id {
			return &f, nil
		}
	}
	return nil, nil
}

func statusMismatch(d *schema.ResourceData, factor *sdk.Factor) bool {
	status := d.Get("active").(bool)

	// I miss ternary operators
	if factor != nil && factor.Status == statusActive {
		return !status
	}

	return status
}
