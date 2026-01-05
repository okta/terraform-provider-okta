package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
)

const oidcIdp = "OIDC"

func dataSourceIdpOidc() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdpOidcRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"name"},
				Description:   "Id of idp.",
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id"},
				Description:   "Name of the idp.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of idp.",
			},
			"authorization_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IdP Authorization Server (AS) endpoint to request consent from the user and obtain an authorization code grant.",
			},
			"authorization_binding": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The method of making an authorization request.",
			},
			"token_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IdP Authorization Server (AS) endpoint to exchange the authorization code grant for an access token.",
			},
			"token_binding": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The method of making a token request.",
			},
			"user_info_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Protected resource endpoint that returns claims about the authenticated user.",
			},
			"user_info_binding": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The method of making a user info request.",
			},
			"jwks_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Endpoint where the keys signer publishes its keys in a JWK Set.",
			},
			"jwks_binding": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: " The method of making a request for the OIDC JWKS.",
			},
			"scopes": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "The scopes of the IdP.",
			},
			"protocol_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of protocol to use.",
			},
			"client_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique identifier issued by AS for the Okta IdP instance.",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Client secret issued by AS for the Okta IdP instance.",
			},
			"issuer_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URI that identifies the issuer.",
			},
			"issuer_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates whether Okta uses the original Okta org domain URL, a custom domain URL, or dynamic.",
			},
			"max_clock_skew": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Maximum allowable clock-skew when processing messages from the IdP.",
			},
		},
		Description: "Get a OIDC IdP from Okta.",
	}
}

func dataSourceIdpOidcRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id := d.Get("id").(string)
	name := d.Get("name").(string)
	if id == "" && name == "" {
		return diag.Errorf("config must provide either 'id' or 'name' to retrieve the IdP")
	}
	var (
		err  error
		oidc *sdk.IdentityProvider
	)
	if id != "" {
		oidc, err = getIdentityProviderByID(ctx, meta, id, oidcIdp)
	} else {
		oidc, err = getIdpByNameAndType(ctx, meta, name, oidcIdp)
	}
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(oidc.Id)
	_ = d.Set("name", oidc.Name)
	_ = d.Set("type", oidc.Type)
	syncEndpoint("authorization", oidc.Protocol.Endpoints.Authorization, d)
	syncEndpoint("token", oidc.Protocol.Endpoints.Token, d)
	syncEndpoint("user_info", oidc.Protocol.Endpoints.UserInfo, d)
	syncEndpoint("jwks", oidc.Protocol.Endpoints.Jwks, d)
	if oidc.Protocol != nil {
		_ = d.Set("protocol_type", oidc.Protocol.Type)
		if oidc.Protocol.Credentials != nil && oidc.Protocol.Credentials.Client != nil {
			_ = d.Set("client_secret", oidc.Protocol.Credentials.Client.ClientSecret)
			_ = d.Set("client_id", oidc.Protocol.Credentials.Client.ClientId)
		}
		if oidc.Protocol.Issuer != nil {
			_ = d.Set("issuer_url", oidc.Protocol.Issuer.Url)
		}
	}
	if oidc.Policy.MaxClockSkewPtr != nil {
		_ = d.Set("max_clock_skew", oidc.Policy.MaxClockSkewPtr)
	}
	_ = d.Set("scopes", utils.ConvertStringSliceToSet(oidc.Protocol.Scopes))
	if oidc.IssuerMode != "" {
		_ = d.Set("issuer_mode", oidc.IssuerMode)
	}
	return nil
}
