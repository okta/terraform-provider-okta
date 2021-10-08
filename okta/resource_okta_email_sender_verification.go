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
		ReadContext:   resourceEmailSenderVerificationRead,
		DeleteContext: resourceEmailSenderVerificationDelete,
		Importer:      nil,
		Schema: map[string]*schema.Schema{
			"sender_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Email sender ID",
			},
		},
	}
}

func resourceEmailSenderVerificationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sender, _, err := getSupplementFromMetadata(m).GetEmailSender(ctx, d.Get("sender_id").(string))
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
	_, err = getSupplementFromMetadata(m).ValidateEmailSender(ctx, sender.ID, esv)
	if err != nil {
		return diag.Errorf("failed to verify custom email sender: %v", err)
	}
	d.SetId(d.Get("sender_id").(string))
	return nil
}

func resourceEmailSenderVerificationRead(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return nil
}

func resourceEmailSenderVerificationDelete(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return nil
}
