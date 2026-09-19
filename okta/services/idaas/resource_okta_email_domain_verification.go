package idaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func resourceEmailDomainVerification() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEmailDomainVerificationCreate,
		ReadContext:   utils.ResourceFuncNoOp,
		DeleteContext: utils.ResourceFuncNoOp,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				emailDomain, _, err := getOktaV3ClientFromMetadata(meta).EmailDomainAPI.GetEmailDomain(ctx, d.Id()).Execute()
				if err != nil {
					return nil, fmt.Errorf("failed to get email domain for import: %v", err)
				}
				if !IsDomainValidated(emailDomain.GetValidationStatus()) {
					return nil, fmt.Errorf("cannot import email domain verification: email domain %q is not verified (current validation status: %s)", d.Id(), emailDomain.GetValidationStatus())
				}
				_ = d.Set("email_domain_id", d.Id())
				return []*schema.ResourceData{d}, nil
			},
		},
		Description: "Verifies the email domain. The resource won't be created if the email domain could not be verified.",
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
	emailDomainId := d.Get("email_domain_id").(string)
	// Verifying an already-verified email domain returns a 400, so check the
	// current state first and skip verification if it has already been validated.
	if emailDomain, _, err := getOktaV3ClientFromMetadata(meta).EmailDomainAPI.GetEmailDomain(ctx, emailDomainId).Execute(); err == nil && emailDomain != nil && IsDomainValidated(emailDomain.GetValidationStatus()) {
		d.SetId(emailDomainId)
		return nil
	}
	_, _, err := getOktaV3ClientFromMetadata(meta).EmailDomainAPI.VerifyEmailDomain(ctx, emailDomainId).Execute()
	if err != nil {
		return diag.Errorf("failed to verify email domain: %v", err)
	}
	d.SetId(emailDomainId)
	return nil
}
