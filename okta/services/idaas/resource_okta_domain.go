package idaas

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v4/okta"
	"github.com/okta/terraform-provider-okta/okta/utils"
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
		Description: "Manages custom domain for your organization.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Custom Domain name",
				ForceNew:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// Okta API downcases domain names
					return strings.EqualFold(old, new)
				},
			},
			"certificate_source_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Certificate source type that indicates whether the certificate is provided by the user or Okta. Accepted values: `MANUAL`, `OKTA_MANAGED`. Warning: Use of OKTA_MANAGED requires a feature flag to be enabled. Default value = MANUAL",
				Default:     "MANUAL",
			},
			"validation_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the domain",
			},
			"brand_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Brand id of the domain",
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

func resourceDomainCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	domainRequest, err := buildDomain(d)
	if err != nil {
		return diag.Errorf("failed to build domain: %v", err)
	}
	domain, _, err := getOktaV3ClientFromMetadata(meta).CustomDomainAPI.CreateCustomDomain(ctx).Domain(domainRequest).Execute()
	if err != nil {
		return diag.Errorf("failed to create domain: %v", err)
	}
	d.SetId(domain.GetId())
	if brandId, ok := d.GetOk("brand_id"); ok {
		_, _, err = getOktaV3ClientFromMetadata(meta).CustomDomainAPI.ReplaceCustomDomain(ctx, d.Id()).UpdateDomain(okta.UpdateDomain{BrandId: brandId.(string)}).Execute()
		if err != nil {
			return diag.Errorf("failed to update domain: %v", err)
		}
		_, err = getOktaV3ClientFromMetadata(meta).CustomizationAPI.DeleteBrand(ctx, domain.GetBrandId()).Execute()
		if err != nil {
			return diag.Errorf("failed to delete brand: %v", err)
		}
	}
	return resourceDomainRead(ctx, d, meta)
}

func resourceDomainRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	domain, resp, err := getOktaV3ClientFromMetadata(meta).CustomDomainAPI.GetCustomDomain(ctx, d.Id()).Execute()
	if err := utils.SuppressErrorOn404_V3(resp, err); err != nil {
		return diag.Errorf("failed to get domain: %v", err)
	}
	if domain == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", domain.GetDomain())
	d.Set("certificate_source_type", domain.GetCertificateSourceType())
	d.Set("brand_id", domain.GetBrandId())
	d.Set("validation_status", domain.GetValidationStatus())

	arr := make([]map[string]interface{}, len(domain.DnsRecords))
	for i := range domain.DnsRecords {
		arr[i] = map[string]interface{}{
			"expiration":  domain.DnsRecords[i].GetExpiration(),
			"fqdn":        domain.DnsRecords[i].GetFqdn(),
			"record_type": domain.DnsRecords[i].GetRecordType(),
			"values":      utils.ConvertStringSliceToInterfaceSlice(domain.DnsRecords[i].GetValues()),
		}
	}
	err = utils.SetNonPrimitives(d, map[string]interface{}{"dns_records": arr})
	if err != nil {
		return diag.Errorf("failed to set DNS records: %v", err)
	}
	if domain.GetValidationStatus() == "IN_PROGRESS" || domain.GetValidationStatus() == "VERIFIED" || domain.GetValidationStatus() == "COMPLETED" {
		return nil
	}

	return nil
}

func resourceDomainDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logger(meta).Info("deleting domain", "id", d.Id())
	_, err := getOktaV3ClientFromMetadata(meta).CustomDomainAPI.DeleteCustomDomain(ctx, d.Id()).Execute()
	if err != nil {
		return diag.Errorf("failed to delete domain: %v", err)
	}
	return nil
}

func resourceDomainUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_, err := validateDomain(ctx, d, meta, d.Get("validation_status").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	_, _, err = getOktaV3ClientFromMetadata(meta).CustomDomainAPI.ReplaceCustomDomain(ctx, d.Id()).UpdateDomain(okta.UpdateDomain{BrandId: d.Get("brand_id").(string)}).Execute()
	if err != nil {
		return diag.Errorf("failed to update domain: %v", err)
	}
	return resourceDomainRead(ctx, d, meta)
}

func validateDomain(ctx context.Context, d *schema.ResourceData, meta interface{}, validationStatus string) (*okta.DomainResponse, error) {
	if validationStatus == "IN_PROGRESS" || validationStatus == "VERIFIED" || validationStatus == "COMPLETED" {
		return nil, nil
	}
	domain, _, err := getOktaV3ClientFromMetadata(meta).CustomDomainAPI.VerifyDomain(ctx, d.Id()).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to verify domain: %v", err)
	}
	return domain, nil
}

func buildDomain(d *schema.ResourceData) (okta.DomainRequest, error) {
	return okta.DomainRequest{
		Domain:                d.Get("name").(string),
		CertificateSourceType: (d.Get("certificate_source_type").(string)),
	}, nil
}
