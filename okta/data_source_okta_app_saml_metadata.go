package okta

import (
	"context"
	"encoding/xml"
	"fmt"

	"github.com/crewjam/saml"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppMetadataSaml() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppMetadataSamlRead,
		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key_id": {
				Type:     schema.TypeString,
				Required: true,
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
			"certificate": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"want_authn_requests_signed": {
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

func dataSourceAppMetadataSamlRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Get("app_id").(string)
	kid := d.Get("key_id").(string)
	metadata, err := getSupplementFromMetadata(m).GetSAMLMetdata(id, kid)
	if err != nil {
		return diag.Errorf("failed to get SAML metadata: %v", err)
	}
	d.SetId(fmt.Sprintf("%s/%s_metadata", id, kid))

	_ = d.Set("metadata", string(metadata))
	metadataRoot := &saml.EntityDescriptor{}
	err = xml.Unmarshal(metadata, metadataRoot)
	if err != nil {
		return diag.Errorf("could not parse SAML app metadata: %v", err)
	}

	desc := metadataRoot.IDPSSODescriptors[0]
	syncSamlEndpointBinding(d, desc.SingleSignOnServices)
	_ = d.Set("entity_id", metadataRoot.EntityID)
	_ = d.Set("want_authn_requests_signed", desc.WantAuthnRequestsSigned)
	_ = d.Set("certificate", desc.KeyDescriptors[0].KeyInfo.Certificate)
	return nil
}
