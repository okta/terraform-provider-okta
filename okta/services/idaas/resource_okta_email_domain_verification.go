package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func resourceEmailDomainVerification() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEmailDomainVerificationCreate,
		ReadContext:   utils.ResourceFuncNoOp,
		DeleteContext: utils.ResourceFuncNoOp,
		Importer:      nil,
		Description:   "Verifies the email domain. The resource won't be created if the email domain could not be verified.",
		Schema: map[string]*schema.Schema{
			"email_domain_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Email domain ID",
			},
		},
	}
}

func resourceEmailDomainVerificationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_, _, err := getOktaV3ClientFromMetadata(meta).EmailDomainAPI.VerifyEmailDomain(ctx, d.Get("email_domain_id").(string)).Execute()
	if err != nil {
		return diag.Errorf("failed to verify email domain: %v", err)
	}
	d.SetId(d.Get("email_domain_id").(string))
	return nil
}
