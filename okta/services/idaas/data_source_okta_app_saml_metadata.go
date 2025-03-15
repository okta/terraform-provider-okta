package idaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppMetadataSaml() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppMetadataSamlRead,
		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The application ID.",
			},
			"key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Certificate Key ID.",
			},
			"metadata": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Raw metadata of application.",
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
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public certificate from application metadata.",
			},
			"want_authn_requests_signed": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether authn requests are signed.",
			},
			"entity_id": {
				Type:        schema.TypeString,
				Description: "Entity URL for instance https://www.okta.com/saml2/service-provider/sposcfdmlybtwkdcgtuf",
				Computed:    true,
			},
		},
		Description: "Get a SAML application's metadata from Okta.",
	}
}

func dataSourceAppMetadataSamlRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id := d.Get("app_id").(string)
	kid := d.Get("key_id").(string)
	metadata, metadataRoot, err := getAPISupplementFromMetadata(meta).GetSAMLMetadata(ctx, id, kid)
	if err != nil {
		return diag.Errorf("failed to get app's SAML metadata: %v", err)
	}
	d.SetId(fmt.Sprintf("%s/%s_metadata", id, kid))
	_ = d.Set("metadata", string(metadata))
	desc := metadataRoot.IDPSSODescriptors[0]
	syncSamlEndpointBinding(d, desc.SingleSignOnServices)
	_ = d.Set("entity_id", metadataRoot.EntityID)
	_ = d.Set("want_authn_requests_signed", desc.WantAuthnRequestsSigned)
	_ = d.Set("certificate", desc.KeyDescriptors[0].KeyInfo.X509Data.X509Certificates[0].Data)
	return nil
}
