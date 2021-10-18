package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourceEmailSender() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEmailSenderCreate,
		ReadContext:   resourceEmailSenderRead,
		UpdateContext: resourceEmailSenderUpdate,
		DeleteContext: resourceEmailSenderDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"from_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of sender",
			},
			"from_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Email address to send from ",
			},
			"subdomain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Mail domain to send from",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Verification status",
			},
			"dns_records": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "TXT and CNAME records to be registered for the Domain",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"fqdn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS record name",
						},
						"record_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Record type can be TXT or CNAME",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS verification value",
						},
					},
				},
			},
		},
	}
}

func resourceEmailSenderCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sender, _, err := getSupplementFromMetadata(m).CreateEmailSender(ctx, buildEmailSender(d))
	if err != nil {
		return diag.Errorf("failed to create custom email sender: %v", err)
	}
	d.SetId(sender.ID)
	return resourceEmailSenderRead(ctx, d, m)
}

func resourceEmailSenderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sender, resp, err := getSupplementFromMetadata(m).GetEmailSender(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get custom email sender: %v", err)
	}
	if sender == nil || sender.Status == "DELETED" {
		d.SetId("")
		return nil
	}
	_ = d.Set("from_name", sender.FromName)
	_ = d.Set("from_address", sender.FromAddress)
	_ = d.Set("subdomain", sender.ValidationSubdomain)
	_ = d.Set("status", sender.Status)
	arr := make([]map[string]interface{}, len(sender.DNSValidation))
	for i := range sender.DNSValidation {
		arr[i] = map[string]interface{}{
			"fqdn":        sender.DNSValidation[i].Fqdn,
			"record_type": sender.DNSValidation[i].RecordType,
			"value":       sender.DNSValidation[i].VerificationValue,
		}
	}
	err = setNonPrimitives(d, map[string]interface{}{"dns_records": arr})
	if err != nil {
		return diag.Errorf("failed to set DNS records: %v", err)
	}
	return nil
}

func resourceEmailSenderUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	_, _, err := getSupplementFromMetadata(m).UpdateEmailSender(ctx, buildEmailSender(d))
	if err != nil {
		return diag.Errorf("failed to update custom email sender: %v", err)
	}
	return resourceEmailSenderRead(ctx, d, m)
}

func resourceEmailSenderDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sender, resp, err := getSupplementFromMetadata(m).GetEmailSender(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get custom email sender: %v", err)
	}
	if sender == nil || sender.Status == "DELETED" {
		return nil
	}
	if sender.Status == "VERIFIED" {
		resp, err = getSupplementFromMetadata(m).DisableVerifiedEmailSender(ctx, sdk.DisableActiveEmailSender{ActiveID: sender.ID})
	} else {
		resp, err = getSupplementFromMetadata(m).DisableUnverifiedEmailSender(ctx, sdk.DisableInactiveEmailSender{
			PendingID: sender.ID,
		})
	}
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to delete custom email sender: %v", err)
	}
	return nil
}

func buildEmailSender(d *schema.ResourceData) sdk.EmailSender {
	return sdk.EmailSender{
		FromName:            d.Get("from_name").(string),
		FromAddress:         d.Get("from_address").(string),
		ValidationSubdomain: d.Get("subdomain").(string),
	}
}
