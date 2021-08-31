package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func resourceDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainCreate,
		ReadContext:   resourceDomainRead,
		UpdateContext: resourceDomainUpdate,
		DeleteContext: resourceDomainDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Custom Domain name",
				ForceNew:    true,
			},
			"verify": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether the domain should be verified during creation",
				Default:     false,
			},
			"validation_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the domain",
			},
			"dns_records": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "TXT and CNAME records to be registered for the Domain",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"expiration": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "TXT record expiration",
						},
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
						"values": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "DNS verification value",
						},
					},
				},
			},
		},
	}
}

func resourceDomainCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	domain, _, err := getOktaClientFromMetadata(m).Domain.CreateDomain(ctx, buildDomain(d))
	if err != nil {
		return diag.Errorf("failed to create domain: %v", err)
	}
	d.SetId(domain.Id)
	if d.Get("verify").(bool) {
		_, _, err := getOktaClientFromMetadata(m).Domain.VerifyDomain(ctx, domain.Id)
		if err != nil {
			return diag.Errorf("failed to verify domain: %v", err)
		}
	}
	return resourceDomainRead(ctx, d, m)
}

func resourceDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	domain, resp, err := getOktaClientFromMetadata(m).Domain.GetDomain(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get domain: %v", err)
	}
	if domain == nil {
		d.SetId("")
		return nil
	}
	vd, err := validateDomain(ctx, d, m, domain.ValidationStatus)
	if err != nil {
		return diag.FromErr(err)
	}
	if vd != nil {
		_ = d.Set("validation_status", vd.ValidationStatus)
	} else {
		_ = d.Set("validation_status", domain.ValidationStatus)
	}
	arr := make([]map[string]interface{}, len(domain.DnsRecords))
	for i := range domain.DnsRecords {
		arr[i] = map[string]interface{}{
			"expiration":  domain.DnsRecords[i].Expiration,
			"fqdn":        domain.DnsRecords[i].Fqdn,
			"record_type": domain.DnsRecords[i].RecordType,
			"values":      convertStringSliceToInterfaceSlice(domain.DnsRecords[i].Values),
		}
	}
	err = setNonPrimitives(d, map[string]interface{}{"dns_records": arr})
	if err != nil {
		return diag.Errorf("failed to set OAuth application properties: %v", err)
	}
	if domain.ValidationStatus == "IN_PROGRESS" || domain.ValidationStatus == "VERIFIED" || domain.ValidationStatus == "COMPLETED" {
		return nil
	}
	if d.Get("verify").(bool) {
		_, _, err := getOktaClientFromMetadata(m).Domain.VerifyDomain(ctx, d.Id())
		if err != nil {
			return diag.Errorf("failed to verify domain: %v", err)
		}
	}

	return nil
}

func resourceDomainDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	logger(m).Info("deleting domain", "id", d.Id())
	_, err := getOktaClientFromMetadata(m).Domain.DeleteDomain(ctx, d.Id())
	if err != nil {
		return diag.Errorf("failed to delete domain: %v", err)
	}
	return nil
}

func resourceDomainUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	_, err := validateDomain(ctx, d, m, d.Get("validation_status").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceDomainRead(ctx, d, m)
}

func validateDomain(ctx context.Context, d *schema.ResourceData, m interface{}, validationStatus string) (*okta.Domain, error) {
	if validationStatus == "IN_PROGRESS" || validationStatus == "VERIFIED" || validationStatus == "COMPLETED" {
		return nil, nil
	}
	if !d.Get("verify").(bool) {
		return nil, nil
	}
	domain, _, err := getOktaClientFromMetadata(m).Domain.VerifyDomain(ctx, d.Id())
	if err != nil {
		return nil, fmt.Errorf("failed to verify domain: %v", err)
	}
	return domain, nil
}

func buildDomain(d *schema.ResourceData) okta.Domain {
	return okta.Domain{
		Domain:                d.Get("name").(string),
		CertificateSourceType: "MANUAL",
	}
}
