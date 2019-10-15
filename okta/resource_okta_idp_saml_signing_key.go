package okta

import (
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-okta/sdk"
)

func resourceIdpSigningKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceIdpSigningKeyCreate,
		Read:   resourceIdpSigningKeyRead,
		Delete: resourceIdpSigningKeyDelete,
		Exists: resourceIdpSigningKeyExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"x5c": &schema.Schema{
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				ForceNew:    true,
				Description: "base64-encoded X.509 certificate chain with DER encoding",
			},
			"created": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"expires_at": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			// Just an alias for id
			"kid": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"kty": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"use": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"x5t_s256": &schema.Schema{
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

	if is404(resp.StatusCode) {
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	d.Set("created", key.Created)
	d.Set("expires_at", key.ExpiresAt)
	d.Set("kid", key.Kid)
	d.Set("kty", key.Kty)
	d.Set("use", key.Use)
	d.Set("x5t_s256", key.X5T256)

	return setNonPrimitives(d, map[string]interface{}{
		"x5c": convertStringSetToInterface(key.X5C),
	})
}

func resourceIdpSigningKeyDelete(d *schema.ResourceData, m interface{}) error {
	resp, err := getSupplementFromMetadata(m).DeleteIdentityProviderCertificate(d.Id())
	return suppressErrorOn404(resp, err)
}
