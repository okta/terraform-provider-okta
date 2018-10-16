package okta

import (
	"github.com/articulate/oktasdk-go/okta"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

// Predefined second authentication factors. They must be activated in order to use them in MFA policies.
func resourceFactor() *schema.Resource {
	return &schema.Resource{
		Create: resourceFactorCreate,
		Read:   resourceFactorRead,
		Update: resourceFactorUpdate,
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
				Description: "Factor provider name",
				ForceNew:    true,
			},
			"status": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"ACTIVE", "INACTIVE", "NOT_SETUP"}, false),
				Description:  "Status of MFA provider",
				Optional:     true,
				Default:      "ACTIVE",
			},
		},
	}
}

func resourceFactorExists(d *schema.ResourceData, m interface{}) (bool, error) {
	factor, err := findFactor(d, m)

	return err == nil && factor != nil, err
}

func resourceFactorCreate(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("provider_id").(string))
	err := activateFactor(d, m)

	if err != nil {
		return err
	}

	return resourceFactorRead(d, m)
}

func resourceFactorDelete(d *schema.ResourceData, m interface{}) error {
	var err error
	client := getClientFromMetadata(m)

	if d.Get("status").(string) == "ACTIVE" {
		_, _, err = client.Org.DeactivateFactor(d.Id())
	}

	return err
}

func resourceFactorRead(d *schema.ResourceData, m interface{}) error {
	factor, err := findFactor(d, m)
	if err != nil {
		return err
	}

	d.Set("status", factor.Status)
	d.Set("provider_id", factor.Id)

	return nil
}

func resourceFactorUpdate(d *schema.ResourceData, m interface{}) error {
	if d.HasChange("status") {
		err := activateFactor(d, m)

		if err != nil {
			return err
		}
	}

	return resourceFactorRead(d, m)
}

func activateFactor(d *schema.ResourceData, m interface{}) error {
	var err error
	client := getClientFromMetadata(m)
	id := d.Id()

	switch d.Get("status").(string) {
	case "ACTIVE":
		_, _, err = client.Org.ActivateFactor(id)
	case "INACTIVE":
		_, _, err = client.Org.DeactivateFactor(id)
	}

	return err
}

func findFactor(d *schema.ResourceData, m interface{}) (*okta.Factor, error) {
	client := getClientFromMetadata(m)
	factorList, _, err := client.Org.ListFactors()

	if err != nil {
		return nil, err
	}

	for _, f := range factorList {
		if f.Id == d.Id() {
			return f, nil
		}
	}

	return nil, nil
}
