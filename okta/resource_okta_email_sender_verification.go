package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourceEmailSenderVerification() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEmailSenderVerificationCreate,
		ReadContext:   resourceFuncNoOp,
		DeleteContext: resourceFuncNoOp,
		Importer:      nil,
		Description:   "Verifies the email sender. The resource won't be created if the email sender could not be verified.",
		Schema: map[string]*schema.Schema{
			"sender_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Email sender ID",
			},
		},
		DeprecationMessage: "The api for this resource has been deprecated. Please use okta_email_domain_verification instead",
	}
}

func resourceEmailSenderVerificationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sender, _, err := getAPISupplementFromMetadata(meta).GetEmailSender(ctx, d.Get("sender_id").(string))
	if err != nil {
		return diag.Errorf("failed to get custom email sender: %v", err)
	}
	esv := sdk.EmailSenderValidation{
		PendingFromAddress:      sender.FromAddress,
		PendingFromName:         sender.FromName,
		PendingValidationDomain: sender.ValidationSubdomain,
		PendingID:               sender.ID,
		PendingDNSValidation:    sender.DNSValidation,
	}
	_, err = getAPISupplementFromMetadata(meta).ValidateEmailSender(ctx, sender.ID, esv)
	if err != nil {
		return diag.Errorf("failed to verify custom email sender: %v", err)
	}
	d.SetId(d.Get("sender_id").(string))
	return nil
}
