package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
)

const saml2Idp = "SAML2"

func dataSourceIdpSaml() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdpSamlRead,
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
			"acs_binding": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"acs_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sso_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sso_binding": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sso_destination": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subject_format": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"subject_filter": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"issuer": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"issuer_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"audience": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceIdpSamlRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Get("id").(string)
	name := d.Get("name").(string)
	if id == "" && name == "" {
		return diag.Errorf("config must provide either 'id' or 'name' to retrieve the IdP")
	}
	var (
		err error
		idp *sdk.SAMLIdentityProvider
	)
	if id != "" {
		idp, err = getIdentityProviderByID(ctx, m, id)
	} else {
		idp, err = getIdpByName(ctx, m, name)
	}
	if err != nil {
		return diag.Errorf("failed to get identity provider: %v", err)
	}
	if idp == nil && id != "" {
		return diag.Errorf("identity provider with id '%s' does not exist", id)
	}
	if idp == nil && name != "" {
		return diag.Errorf("identity provider with name '%s' does not exist", name)
	}

	d.SetId(idp.ID)
	_ = d.Set("name", idp.Name)
	_ = d.Set("type", idp.Type)
	_ = d.Set("acs_binding", idp.Protocol.Endpoints.Acs.Binding)
	_ = d.Set("acs_type", idp.Protocol.Endpoints.Acs.Type)
	_ = d.Set("sso_url", idp.Protocol.Endpoints.Sso.URL)
	_ = d.Set("sso_binding", idp.Protocol.Endpoints.Sso.Binding)
	_ = d.Set("sso_destination", idp.Protocol.Endpoints.Sso.Destination)
	_ = d.Set("subject_filter", idp.Policy.Subject.Filter)
	_ = d.Set("kid", idp.Protocol.Credentials.Trust.Kid)
	_ = d.Set("issuer", idp.Protocol.Credentials.Trust.Issuer)
	_ = d.Set("audience", idp.Protocol.Credentials.Trust.Audience)
	err = setNonPrimitives(d, map[string]interface{}{
		"subject_format": convertStringSetToInterface(idp.Policy.Subject.Format),
	})
	if err != nil {
		return diag.Errorf("failed to set SAML identity provider properties: %v", err)
	}
	return nil
}

func getIdentityProviderByID(ctx context.Context, m interface{}, id string) (*sdk.SAMLIdentityProvider, error) {
	var idp sdk.SAMLIdentityProvider
	_, resp, err := getSupplementFromMetadata(m).GetIdentityProvider(ctx, id, &idp)
	return &idp, responseErr(resp, err)
}

func getIdpByName(ctx context.Context, m interface{}, label string) (*sdk.SAMLIdentityProvider, error) {
	var idps []*sdk.SAMLIdentityProvider
	queryParams := query.Params{Limit: 1, Q: label}
	_, resp, err := getSupplementFromMetadata(m).ListIdentityProviders(ctx, &idps, &queryParams)
	if len(idps) > 0 {
		return idps[0], nil
	}
	return nil, responseErr(resp, err)
}
