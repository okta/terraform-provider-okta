package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourceDomainCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainCertificateCreate,
		ReadContext:   resourceFuncNoOp,
		UpdateContext: resourceDomainCertificateUpdate,
		DeleteContext: resourceFuncNoOp,
		Importer:      nil,
		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Domain's ID",
				ForceNew:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Certificate type",
				DefaultFunc: func() (interface{}, error) {
					return "PEM", nil
				},
			},
			"certificate": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Certificate content",
			},
			"private_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Certificate private key",
			},
			"certificate_chain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Certificate chain",
			},
		},
	}
}

func resourceDomainCertificateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := buildDomainCertificate(d)
	_, err := getOktaClientFromMetadata(m).Domain.CreateCertificate(ctx, d.Get("domain_id").(string), c)
	if err != nil {
		return diag.Errorf("failed to create domain's certificate: %v", err)
	}
	d.SetId(d.Get("domain_id").(string))
	return nil
}

func resourceDomainCertificateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := buildDomainCertificate(d)
	_, err := getOktaClientFromMetadata(m).Domain.CreateCertificate(ctx, d.Get("domain_id").(string), c)
	if err != nil {
		return diag.Errorf("failed to update domain's certificate: %v", err)
	}
	return nil
}

func buildDomainCertificate(d *schema.ResourceData) sdk.DomainCertificate {
	return sdk.DomainCertificate{
		Certificate:      d.Get("certificate").(string),
		CertificateChain: d.Get("certificate_chain").(string),
		PrivateKey:       d.Get("private_key").(string),
		Type:             d.Get("type").(string),
	}
}
