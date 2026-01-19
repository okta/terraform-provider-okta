package idaas

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oktav5sdk "github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func dataSourceDomain() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDomainRead,
		Schema: map[string]*schema.Schema{
			"domain_id_or_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Brand ID",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the Domain",
			},
			"domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Domain name",
			},
			"certificate_source_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate source type that indicates whether the certificate is provided by the user or Okta. Values: MANUAL, OKTA_MANAGED",
			},
			"validation_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the domain. Values: NOT_STARTED, IN_PROGRESS, VERIFIED, COMPLETED",
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
			"public_certificate": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Certificate metadata for the Domain",
			},
		},
		Description: "Get a domain from Okta.",
	}
}

func dataSourceDomainRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	did, _ := d.GetOk("domain_id_or_name")
	domainID := did.(string)

	domainList, _, err := getOktaV5ClientFromMetadata(meta).CustomDomainAPI.ListCustomDomains(ctx).Execute()
	if err != nil {
		return diag.Errorf("failed to get domains: %v", err)
	}

	var domain *oktav5sdk.DomainResponse
	for _, _domain := range domainList.Domains {
		if _domain.GetId() == domainID {
			domain = &_domain
			break
		}
		if _domain.GetDomain() == domainID {
			domain = &_domain
			break
		}
		if strings.EqualFold(_domain.GetDomain(), domainID) {
			domain = &_domain
			break
		}
	}
	if domain == nil {
		return diag.Errorf("failed to find domain by id or name: %q", domainID)
	}

	d.SetId(domain.GetId())
	d.Set("domain", domain.Domain)
	d.Set("validation_status", domain.ValidationStatus)
	d.Set("certificate_source_type", domain.CertificateSourceType)
	// retrieve DNS records by calling api/v1/domains/{domainId} and reassign domain to response of GetCustomDomain()
	domain, _, err = getOktaV5ClientFromMetadata(meta).CustomDomainAPI.GetCustomDomain(ctx, domainID).Execute()
	if err != nil {
		return diag.Errorf("failed to get domain to retrieve DNS records: %v", err)
	}
	arr := make([]map[string]interface{}, len(domain.DnsRecords))
	for i := range domain.GetDnsRecords() {
		arr[i] = map[string]interface{}{
			"expiration":  domain.DnsRecords[i].Expiration,
			"fqdn":        domain.DnsRecords[i].Fqdn,
			"record_type": domain.DnsRecords[i].RecordType,
			"values":      utils.ConvertStringSliceToInterfaceSlice(domain.DnsRecords[i].Values),
		}
	}
	err = utils.SetNonPrimitives(d, map[string]interface{}{"dns_records": arr})
	if err != nil {
		return diag.Errorf("failed to set DNS records: %v", err)
	}

	if domain.PublicCertificate != nil {
		cert := map[string]interface{}{
			"subject":     domain.PublicCertificate.Subject,
			"fingerprint": domain.PublicCertificate.Fingerprint,
			"expiration":  domain.PublicCertificate.Expiration,
		}
		d.Set("publice_certificate", cert)
	}

	return nil
}
