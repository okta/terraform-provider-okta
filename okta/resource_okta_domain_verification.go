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
		ReadContext:   resourceFuncNoOp,
		DeleteContext: resourceFuncNoOp,
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
	boc := newExponentialBackOffWithContext(ctx, 30*time.Second)
	err := backoff.Retry(func() error {
		domain, _, err := getOktaClientFromMetadata(m).Domain.VerifyDomain(ctx, d.Get("domain_id").(string))
		if doNotRetry(m, err) {
			return backoff.Permanent(err)
		}
		if err != nil {
			return backoff.Permanent(fmt.Errorf("failed to verify domain: %v", err))
		}
		if !isDomainValidated(domain.ValidationStatus) {
			return fmt.Errorf("failed to verify domain after several attempts, current validation status: %s", domain.ValidationStatus)
		}
		return nil
	}, boc)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(d.Get("domain_id").(string))
	return nil
}

// Status of the domain. Accepted values: NOT_STARTED, IN_PROGRESS, VERIFIED, COMPLETED
func isDomainValidated(validationStatus string) bool {
	switch validationStatus {
	case "VERIFIED":
		return true
	case "COMPLETED":
		return true
	}
	return false
}
