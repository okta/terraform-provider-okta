package okta

import (
	"context"
	"crypto/sha512"
	"encoding/base64"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func resourceDomainCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainCertificateCreate,
		ReadContext:   resourceDomainCertificateRead,
		UpdateContext: resourceDomainCertificateUpdate,
		DeleteContext: resourceDomainCertificateDelete,
		Importer:      nil,
		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Domain's ID",
				ForceNew:    true,
			},
			"type": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Certificate type",
				ValidateDiagFunc: elemInSlice([]string{"PEM"}),
				DefaultFunc: func() (interface{}, error) {
					return "PEM", nil
				},
			},
			"certificate": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Certificate content",
				StateFunc: func(val interface{}) string {
					h := sha512.New()
					h.Write([]byte(val.(string)))
					return base64.URLEncoding.EncodeToString(h.Sum(nil))
				},
			},
			"private_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Certificate private key",
				StateFunc: func(val interface{}) string {
					h := sha512.New()
					h.Write([]byte(val.(string)))
					return base64.URLEncoding.EncodeToString(h.Sum(nil))
				},
			},
			"certificate_chain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Certificate chain",
				StateFunc: func(val interface{}) string {
					if val.(string) == "" {
						return ""
					}
					h := sha512.New()
					h.Write([]byte(val.(string)))
					return base64.URLEncoding.EncodeToString(h.Sum(nil))
				},
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

func resourceDomainCertificateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

// nothing to do here, since domain's certificate can be deleted
func resourceDomainCertificateDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func buildDomainCertificate(d *schema.ResourceData) okta.DomainCertificate {
	return okta.DomainCertificate{
		Certificate:      d.Get("certificate").(string),
		CertificateChain: d.Get("certificate_chain").(string),
		PrivateKey:       d.Get("private_key").(string),
		Type:             d.Get("type").(string),
	}
}
