package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v3/okta"
)

func resourceEmailDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEmailDomainCreate,
		ReadContext:   resourceEmailDomainRead,
		UpdateContext: resourceEmailDomainUpdate,
		DeleteContext: resourceEmailDomainDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"brand_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Brand id",
			},
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Domain name",
			},
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name",
			},
			"user_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "User name",
			},
			"validation_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the email domain. Values: NOT_STARTED, IN_PROGRESS, VERIFIED, COMPLETED",
			},
			"dns_validation_records": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "TXT and cname records to be registered for the email Domain",
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
							Description: "Record type can be TXT or cname",
						},
						"values": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "DNS record values",
						},
						"expiration": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS TXT record expiration",
						},
					},
				},
			},
		},
	}
}

func resourceEmailDomainCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	emailDomain, _, err := getOktaV3ClientFromMetadata(m).EmailDomainApi.CreateEmailDomain(ctx).EmailDomain(buildEmailDomain(d)).Execute()
	if err != nil {
		return diag.Errorf("failed to create email domain: %v", err)
	}
	d.SetId(emailDomain.GetId())
	return resourceEmailDomainRead(ctx, d, m)
}

func resourceEmailDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	emailDomain, resp, err := getOktaV3ClientFromMetadata(m).EmailDomainApi.GetEmailDomain(ctx, d.Id()).Execute()
	if err := v3suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get email domain: %v", err)
	}
	if emailDomain == nil || emailDomain.GetValidationStatus() == "DELETED" {
		d.SetId("")
		return nil
	}
	_ = d.Set("validation_status", emailDomain.GetValidationStatus())
	_ = d.Set("domain", emailDomain.GetDomain())
	_ = d.Set("display_name", emailDomain.GetDisplayName())
	_ = d.Set("user_name", emailDomain.GetUserName())
	dnsValidation := emailDomain.GetDnsValidationRecords()
	arr := make([]map[string]interface{}, len(dnsValidation))
	for i := range dnsValidation {
		arr[i] = map[string]interface{}{
			"fqdn":        dnsValidation[i].GetFqdn(),
			"record_type": dnsValidation[i].GetRecordType(),
			"expiration":  dnsValidation[i].GetExpiration(),
		}
		if len(dnsValidation[i].GetValues()) > 0 {
			arr[i]["value"] = dnsValidation[i].GetValues()
		}
	}
	err = setNonPrimitives(d, map[string]interface{}{"dns_validation_records": arr})
	if err != nil {
		return diag.Errorf("failed to set DNS validation records: %v", err)
	}
	return nil
}

func resourceEmailDomainUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	_, _, err := getOktaV3ClientFromMetadata(m).EmailDomainApi.ReplaceEmailDomain(ctx, d.Id()).UpdateEmailDomain(buildUpdateEmailDomain(d)).Execute()
	if err != nil {
		return diag.Errorf("failed to update email domain: %v", err)
	}
	return resourceEmailDomainRead(ctx, d, m)
}

func resourceEmailDomainDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	emailDomain, resp, err := getOktaV3ClientFromMetadata(m).EmailDomainApi.GetEmailDomain(ctx, d.Id()).Execute()
	if err := v3suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get email domain: %v", err)
	}
	if emailDomain == nil || emailDomain.GetValidationStatus() == "DELETED" {
		return nil
	}
	_, err = getOktaV3ClientFromMetadata(m).EmailDomainApi.DeleteEmailDomain(ctx, emailDomain.GetId()).Execute()
	if err := v3suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to delete email domain: %v", err)
	}
	return nil
}

func buildEmailDomain(d *schema.ResourceData) okta.EmailDomain {
	return okta.EmailDomain{
		BrandId:     d.Get("brand_id").(string),
		Domain:      d.Get("domain").(string),
		DisplayName: d.Get("display_name").(string),
		UserName:    d.Get("user_name").(string),
	}
}

func buildUpdateEmailDomain(d *schema.ResourceData) okta.UpdateEmailDomain {
	return okta.UpdateEmailDomain{
		DisplayName: d.Get("display_name").(string),
		UserName:    d.Get("user_name").(string),
	}
}
