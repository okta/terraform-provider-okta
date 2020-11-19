package okta

import (
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
)

func resourceIdpSigningKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceIdpSigningKeyCreate,
		Read:   resourceIdpSigningKeyRead,
		Delete: resourceIdpSigningKeyDelete,
		Exists: resourceIdpSigningKeyExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"x5c": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				ForceNew:    true,
				Description: "base64-encoded X.509 certificate chain with DER encoding",
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expires_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// Just an alias for id
			"kid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kty": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"use": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"x5t_s256": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIdpSigningKeyCreate(d *schema.ResourceData, m interface{}) error {
	cert := &sdk.Certificate{
		X5C: convertInterfaceToStringSet(d.Get("x5c")),
	}
	key, resp, err := getSupplementFromMetadata(m).AddIdentityProviderCertificate(cert)
	if err != nil {
		return responseErr(resp, err)
	}

	d.SetId(key.Kid)
	return resourceIdpSigningKeyRead(d, m)
}

func resourceIdpSigningKeyExists(d *schema.ResourceData, m interface{}) (bool, error) {
	_, resp, err := getSupplementFromMetadata(m).GetIdentityProviderCertificate(d.Id())
	return resp.StatusCode == http.StatusOK, err
}

func resourceIdpSigningKeyRead(d *schema.ResourceData, m interface{}) error {
	key, resp, err := getSupplementFromMetadata(m).GetIdentityProviderCertificate(d.Id())

	if resp != nil && is404(resp.StatusCode) {
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	_ = d.Set("created", key.Created)
	_ = d.Set("expires_at", key.ExpiresAt)
	_ = d.Set("kid", key.Kid)
	_ = d.Set("kty", key.Kty)
	_ = d.Set("use", key.Use)
	_ = d.Set("x5t_s256", key.X5T256)

	return setNonPrimitives(d, map[string]interface{}{
		"x5c": convertStringSetToInterface(key.X5C),
	})
}

func resourceIdpSigningKeyDelete(d *schema.ResourceData, m interface{}) error {
	resp, err := getSupplementFromMetadata(m).DeleteIdentityProviderCertificate(d.Id())
	return suppressErrorOn404(resp, err)
}
