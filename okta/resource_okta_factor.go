package okta

import (
	"github.com/articulate/oktasdk-go/okta"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"provider_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						okta.DuoFactor,
						okta.FidoU2fFactor,
						okta.FidoWebauthnFactor,
						okta.GoogleOtpFactor,
						okta.OktaCallFactor,
						okta.OktaOtpFactor,
						okta.OktaPushFactor,
						okta.OktaQuestionFactor,
						okta.OktaSmsFactor,
						okta.RsaTokenFactor,
						okta.SymantecVipFactor,
						okta.YubikeyTokenFactor,
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
	client := getClientFromMetadata(m)

	if d.Get("active").(bool) {
		_, _, err = client.Org.DeactivateFactor(d.Id())
	}

	return err
}

func resourceFactorRead(d *schema.ResourceData, m interface{}) error {
	factor, err := findFactor(d, m)

	if factor == nil {
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	d.Set("active", factor.Status == "ACTIVE")
	d.Set("provider_id", factor.Id)

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
	client := getClientFromMetadata(m)
	id := d.Get("provider_id").(string)

	if d.Get("active").(bool) {
		_, _, err = client.Org.ActivateFactor(id)
	} else {
		_, _, err = client.Org.DeactivateFactor(id)
	}

	return err
}

// This API is in Beta hence the inability to do a single get. I must list then find.
// Fear is clearly not a factor for me.
func findFactor(d *schema.ResourceData, m interface{}) (*okta.Factor, error) {
	client := getClientFromMetadata(m)
	factorList, _, err := client.Org.ListFactors()

	if err != nil {
		return nil, err
	}

	id := d.Get("provider_id").(string)

	for _, f := range factorList {
		if f.Id == id {
			return f, nil
		}
	}

	return nil, nil
}

func statusMismatch(d *schema.ResourceData, factor *okta.Factor) bool {
	status := d.Get("active").(bool)

	// I miss ternary operators
	if factor.Status == "ACTIVE" {
		return !status
	}

	return status
}
