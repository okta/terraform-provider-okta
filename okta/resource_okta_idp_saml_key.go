package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func resourceIdpSigningKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdpSigningKeyCreate,
		ReadContext:   resourceIdpSigningKeyRead,
		DeleteContext: resourceIdpSigningKeyDelete,
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

func resourceIdpSigningKeyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cert := okta.JsonWebKey{
		X5c: convertInterfaceToStringSet(d.Get("x5c")),
	}
	key, _, err := getOktaClientFromMetadata(m).IdentityProvider.CreateIdentityProviderKey(ctx, cert)
	if err != nil {
		return diag.Errorf("failed to create identity provider signing key: %v", err)
	}
	d.SetId(key.Kid)
	return resourceIdpSigningKeyRead(ctx, d, m)
}

func resourceIdpSigningKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	key, resp, err := getOktaClientFromMetadata(m).IdentityProvider.GetIdentityProviderKey(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get identity provider signing key: %v", err)
	}
	if key == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("created", key.Created.UTC().String())
	_ = d.Set("expires_at", key.ExpiresAt.UTC().String())
	_ = d.Set("kid", key.Kid)
	_ = d.Set("kty", key.Kty)
	_ = d.Set("use", key.Use)
	_ = d.Set("x5t_s256", key.X5tS256)
	err = setNonPrimitives(d, map[string]interface{}{
		"x5c": convertStringSetToInterface(key.X5c),
	})
	if err != nil {
		return diag.Errorf("failed to set identity provider signing key properties: %v", err)
	}
	return nil
}

func resourceIdpSigningKeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	_, err := getOktaClientFromMetadata(m).IdentityProvider.DeleteIdentityProviderKey(ctx, d.Id())
	if err != nil {
		return diag.Errorf("failed to delete identity provider signing key: %v", err)
	}
	return nil
}
