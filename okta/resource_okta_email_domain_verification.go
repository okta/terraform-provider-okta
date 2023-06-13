package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceEmailDomainVerification() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEmailDomainVerificationCreate,
		ReadContext:   resourceFuncNoOp,
		DeleteContext: resourceFuncNoOp,
		Importer:      nil,
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

func resourceEmailDomainVerificationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	_, _, err := getOktaV3ClientFromMetadata(m).EmailDomainApi.VerifyEmailDomain(ctx, d.Get("email_domain_id").(string)).Execute()
	if err != nil {
		return diag.Errorf("failed to verify email domain: %v", err)
	}
	d.SetId(d.Get("email_domain_id").(string))
	return nil
}
