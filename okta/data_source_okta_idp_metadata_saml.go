package okta

import (
	"encoding/xml"
	"fmt"

	"github.com/crewjam/saml"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIdpMetadataSaml() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIdpSamlMetadataRead,

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

func dataSourceIdpSamlMetadataRead(d *schema.ResourceData, m interface{}) error {
	id := d.Get("idp_id").(string)
	d.SetId(fmt.Sprintf("%s_metadata", id))
	client := getSupplementFromMetadata(m)
	metadata, _, err := client.GetSAMLIdpMetdata(id)
	if err != nil {
		return err
	}

	d.Set("metadata", string(metadata))
	metadataRoot := &saml.EntityDescriptor{}
	err = xml.Unmarshal(metadata, metadataRoot)
	if err != nil {
		return fmt.Errorf("Could not parse SAML app metadata, error: %s", err)
	}

	desc := metadataRoot.SPSSODescriptors[0]
	syncSamlIndexEndpointBinding(d, desc.AssertionConsumerServices)
	d.Set("entity_id", metadataRoot.EntityID)
	d.Set("authn_request_signed", desc.AuthnRequestsSigned)
	d.Set("assertions_signed", desc.WantAssertionsSigned)
	syncSamlCertificates(d, desc.KeyDescriptors)
	return nil
}
