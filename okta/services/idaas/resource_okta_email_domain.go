package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v4/okta"
	"github.com/okta/terraform-provider-okta/okta/utils"
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
		Description: `Creates email domain. This resource allows you to create and configure an email domain. 
		
**IMPORTANT:** Due to the way Okta's API conflict with terraform design principle, updating the relationship between email_domain and brand is not configurable through terraform and has to be done through clickOps`,
		Schema: map[string]*schema.Schema{
			"brand_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Brand id of the email domain.",
			},
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Mail domain to send from.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name of the email domain.",
			},
			"user_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "User name of the email domain.",
			},
			"validation_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the email domain. Values: NOT_STARTED, IN_PROGRESS, VERIFIED, COMPLETED",
			},
			"validation_subdomain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Subdomain for the email sender's custom mail domain.",
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
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS record value",
						},
						"expiration": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS TXT record expiration",
							Deprecated:  "This field has been removed in the newest go sdk version and has become noop",
						},
					},
				},
			},
		},
	}
}

func resourceEmailDomainCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	emailDomain, _, err := getOktaV3ClientFromMetadata(meta).EmailDomainAPI.CreateEmailDomain(ctx).EmailDomain(buildEmailDomain(d)).Execute()
	if err != nil {
		return diag.Errorf("failed to create email domain: %v", err)
	}
	d.SetId(emailDomain.GetId())
	return resourceEmailDomainRead(ctx, d, meta)
}

func resourceEmailDomainRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	emailDomain, resp, err := getOktaV3ClientFromMetadata(meta).EmailDomainAPI.GetEmailDomain(ctx, d.Id()).Execute()
	if err := utils.SuppressErrorOn404_V3(resp, err); err != nil {
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
			"value":       dnsValidation[i].GetVerificationValue(),
		}
	}
	err = utils.SetNonPrimitives(d, map[string]interface{}{"dns_validation_records": arr})
	if err != nil {
		return diag.Errorf("failed to set DNS validation records: %v", err)
	}
	return nil
}

func resourceEmailDomainUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_, _, err := getOktaV3ClientFromMetadata(meta).EmailDomainAPI.ReplaceEmailDomain(ctx, d.Id()).UpdateEmailDomain(buildUpdateEmailDomain(d)).Execute()
	if err != nil {
		return diag.Errorf("failed to update email domain: %v", err)
	}
	return resourceEmailDomainRead(ctx, d, meta)
}

func resourceEmailDomainDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	emailDomain, resp, err := getOktaV3ClientFromMetadata(meta).EmailDomainAPI.GetEmailDomain(ctx, d.Id()).Execute()
	if err := utils.SuppressErrorOn404_V3(resp, err); err != nil {
		return diag.Errorf("failed to get email domain: %v", err)
	}
	if emailDomain == nil || emailDomain.GetValidationStatus() == "DELETED" {
		return nil
	}
	_, err = getOktaV3ClientFromMetadata(meta).EmailDomainAPI.DeleteEmailDomain(ctx, emailDomain.GetId()).Execute()
	if err := utils.SuppressErrorOn404_V3(resp, err); err != nil {
		return diag.Errorf("failed to delete email domain: %v", err)
	}
	return nil
}

func buildEmailDomain(d *schema.ResourceData) okta.EmailDomain {
	return okta.EmailDomain{
		BrandId:             d.Get("brand_id").(string),
		Domain:              d.Get("domain").(string),
		DisplayName:         d.Get("display_name").(string),
		UserName:            d.Get("user_name").(string),
		ValidationSubdomain: d.Get("validation_subdomain").(*string),
	}
}

func buildUpdateEmailDomain(d *schema.ResourceData) okta.UpdateEmailDomain {
	return okta.UpdateEmailDomain{
		DisplayName: d.Get("display_name").(string),
		UserName:    d.Get("user_name").(string),
	}
}
