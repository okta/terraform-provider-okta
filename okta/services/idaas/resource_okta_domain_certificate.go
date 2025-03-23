package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourceDomainCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainCertificateCreate,
		ReadContext:   utils.ResourceFuncNoOp,
		UpdateContext: resourceDomainCertificateUpdate,
		DeleteContext: utils.ResourceFuncNoOp,
		Importer:      nil,
		Description: `Manages certificate for the domain.
This resource's 'certificate', 'private_key', and 'certificate_chain' attributes
hold actual PEM values and can be referred to by other configs requiring
certificate and private key inputs. This is inline with TF's [best
practices](https://developer.hashicorp.com/terraform/plugin/sdkv2/best-practices/sensitive-state#don-t-encrypt-state)
of not encrypting state.
See [Let's Encrypt Certbot notes](#lets-encrypt-certbot) at the end of this
documentation for notes on how to generate a domain certificate with Let's Encrypt Certbot`,
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
				Description: "Certificate type. Valid value is `PEM`",
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

func resourceDomainCertificateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := buildDomainCertificate(d)
	_, err := getOktaClientFromMetadata(meta).Domain.CreateCertificate(ctx, d.Get("domain_id").(string), c)
	if err != nil {
		return diag.Errorf("failed to create domain's certificate: %v", err)
	}
	d.SetId(d.Get("domain_id").(string))
	return nil
}

func resourceDomainCertificateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := buildDomainCertificate(d)
	_, err := getOktaClientFromMetadata(meta).Domain.CreateCertificate(ctx, d.Get("domain_id").(string), c)
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
