package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/okta/terraform-provider-okta/sdk/query"
)

func resourceIdpSigningKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdpSigningKeyCreate,
		ReadContext:   resourceIdpSigningKeyRead,
		UpdateContext: resourceIdpSigningKeyUpdate,
		DeleteContext: resourceIdpSigningKeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "This resource allows you to create and configure a SAML Identity Provider Signing Key.",
		Schema: map[string]*schema.Schema{
			"x5c": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: "base64-encoded X.509 certificate chain with DER encoding",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date created.",
			},
			"expires_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date the cert expires.",
			},
			"kid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Key ID.",
			},
			"kty": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identifies the cryptographic algorithm family used with the key.",
			},
			"use": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Intended use of the public key.",
			},
			"x5t_s256": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "base64url-encoded SHA-256 thumbprint of the DER encoding of an X.509 certificate.",
			},
		},
	}
}

func resourceIdpSigningKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cert := sdk.JsonWebKey{
		X5c: convertInterfaceToStringSet(d.Get("x5c")),
	}
	key, _, err := getOktaClientFromMetadata(meta).IdentityProvider.CreateIdentityProviderKey(ctx, cert)
	if err != nil {
		return diag.Errorf("failed to create identity provider signing key: %v", err)
	}
	d.SetId(key.Kid)
	return resourceIdpSigningKeyRead(ctx, d, meta)
}

func resourceIdpSigningKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	key, resp, err := getOktaClientFromMetadata(meta).IdentityProvider.GetIdentityProviderKey(ctx, d.Id())
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
		"x5c": convertStringSliceToSet(key.X5c),
	})
	if err != nil {
		return diag.Errorf("failed to set identity provider signing key properties: %v", err)
	}
	return nil
}

func resourceIdpSigningKeyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cert := sdk.JsonWebKey{
		X5c: convertInterfaceToStringSet(d.Get("x5c")),
	}
	client := getOktaClientFromMetadata(meta)
	newKey, _, err := client.IdentityProvider.CreateIdentityProviderKey(ctx, cert)
	if err != nil {
		return diag.Errorf("failed to create identity provider signing key: %v", err)
	}
	idps, _, err := getOktaClientFromMetadata(meta).IdentityProvider.
		ListIdentityProviders(ctx, &query.Params{Limit: defaultPaginationLimit, Type: saml2Idp})
	if err != nil {
		return diag.Errorf("failed to list identity providers: %v", err)
	}
	for i := range idps {
		if idps[i].Protocol.Credentials.Trust.Kid != d.Id() {
			// only update IdPs that are using old key
			continue
		}
		idps[i].Protocol.Credentials.Trust.Kid = newKey.Kid
		_, _, err = client.IdentityProvider.UpdateIdentityProvider(ctx, idps[i].Id, *idps[i])
		if err != nil {
			return diag.Errorf("failed to update identity provider using new key: %v", err)
		}
	}
	_, err = getOktaClientFromMetadata(meta).IdentityProvider.DeleteIdentityProviderKey(ctx, d.Id())
	if err != nil {
		return diag.Errorf("failed to delete identity provider signing key: %v", err)
	}
	d.SetId(newKey.Kid)
	return resourceIdpSigningKeyRead(ctx, d, meta)
}

func resourceIdpSigningKeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_, err := getOktaClientFromMetadata(meta).IdentityProvider.DeleteIdentityProviderKey(ctx, d.Id())
	if err != nil {
		return diag.Errorf("failed to delete identity provider signing key: %v", err)
	}
	return nil
}
