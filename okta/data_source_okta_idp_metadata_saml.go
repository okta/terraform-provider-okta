package okta

import (
	"context"
	"encoding/xml"
	"fmt"

	"github.com/crewjam/saml"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIdpMetadataSaml() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdpSamlMetadataRead,
		Schema: map[string]*schema.Schema{
			"idp_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"metadata": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"http_post_binding": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Post location from the SAML metadata.",
			},
			"http_redirect_binding": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect location from the SAML metadata.",
			},
			"signing_certificate": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"encryption_certificate": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"authn_request_signed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"assertions_signed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"entity_id": {
				Type:        schema.TypeString,
				Description: "Entity URL for instance https://www.okta.com/saml2/service-provider/sposcfdmlybtwkdcgtuf",
				Computed:    true,
			},
		},
	}
}

func dataSourceIdpSamlMetadataRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Get("idp_id").(string)
	d.SetId(fmt.Sprintf("%s_metadata", id))
	metadata, err := getSupplementFromMetadata(m).GetSAMLIdpMetdata(id)
	if err != nil {
		return diag.Errorf("failed to get SAML idp metadata: %v", err)
	}

	_ = d.Set("metadata", string(metadata))
	metadataRoot := &saml.EntityDescriptor{}
	err = xml.Unmarshal(metadata, metadataRoot)
	if err != nil {
		return diag.Errorf("could not parse SAML app metadata, %v", err)
	}

	desc := metadataRoot.SPSSODescriptors[0]
	syncSamlIndexEndpointBinding(d, desc.AssertionConsumerServices)
	_ = d.Set("entity_id", metadataRoot.EntityID)
	_ = d.Set("authn_request_signed", desc.AuthnRequestsSigned)
	_ = d.Set("assertions_signed", desc.WantAssertionsSigned)
	syncSamlCertificates(d, desc.KeyDescriptors)
	return nil
}
