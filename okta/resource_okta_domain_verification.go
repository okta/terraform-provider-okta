package okta

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDomainVerification() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainVerificationCreate,
		ReadContext:   resourceDomainVerificationRead,
		DeleteContext: resourceDomainVerificationDelete,
		Importer:      nil,
		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Domain's ID",
				ForceNew:    true,
			},
		},
	}
}

func resourceDomainVerificationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	bOff := backoff.NewExponentialBackOff()
	bOff.MaxElapsedTime = time.Second * 30
	bOff.InitialInterval = time.Second
	err := backoff.Retry(func() error {
		domain, _, err := getOktaClientFromMetadata(m).Domain.VerifyDomain(ctx, d.Get("domain_id").(string))
		if err != nil {
			return backoff.Permanent(fmt.Errorf("failed to verify domain: %v", err))
		}
		if !isDomainValidated(domain.ValidationStatus) {
			return fmt.Errorf("failed to verify domain after several attempts, current validation status: %s", domain.ValidationStatus)
		}
		return nil
	}, bOff)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(d.Get("domain_id").(string))
	return nil
}

// nothing to do here, since domain should be already verified during creation.
func resourceDomainVerificationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

// nothing to do here, since domain cannot be re-verified
func resourceDomainVerificationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

// Status of the domain. Accepted values: NOT_STARTED, IN_PROGRESS, VERIFIED, COMPLETED
func isDomainValidated(validationStatus string) bool {
	switch validationStatus {
		case "VERIFIED":
		case "COMPLETED":
			return true
	}
	return false
}
