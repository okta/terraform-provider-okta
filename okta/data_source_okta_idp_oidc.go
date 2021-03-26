package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
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
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"authorization_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"authorization_binding": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"token_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"token_binding": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_info_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_info_binding": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"jwks_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"jwks_binding": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scopes": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"protocol_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_secret": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"issuer_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"issuer_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"max_clock_skew": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceIdpOidcRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Get("id").(string)
	name := d.Get("name").(string)
	if id == "" && name == "" {
		return diag.Errorf("config must provide either 'id' or 'name' to retrieve the IdP")
	}
	var (
		err  error
		oidc *okta.IdentityProvider
	)
	if id != "" {
		oidc, err = getIdentityProviderByID(ctx, m, id, oidcIdp)
	} else {
		oidc, err = getIdpByNameAndType(ctx, m, name, oidcIdp)
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
	_ = d.Set("protocol_type", oidc.Protocol.Type)
	_ = d.Set("client_secret", oidc.Protocol.Credentials.Client.ClientSecret)
	_ = d.Set("client_id", oidc.Protocol.Credentials.Client.ClientId)
	_ = d.Set("issuer_url", oidc.Protocol.Issuer.Url)
	_ = d.Set("max_clock_skew", oidc.Policy.MaxClockSkew)
	_ = d.Set("scopes", convertStringSetToInterface(oidc.Protocol.Scopes))
	if oidc.IssuerMode != "" {
		_ = d.Set("issuer_mode", oidc.IssuerMode)
	}
	return nil
}
