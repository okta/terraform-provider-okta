package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
)

const Saml2Idp = "SAML2"

func dataSourceIdpSaml() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdpSamlRead,
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
			"acs_binding": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ACS binding",
			},
			"acs_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Determines whether to publish an instance-specific (trust) or organization (shared) ACS endpoint in the SAML metadata.",
			},
			"sso_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Single sign-on url.",
			},
			"sso_binding": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Single sign-on binding.",
			},
			"sso_destination": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SSO request binding, HTTP-POST or HTTP-REDIRECT.",
			},
			"subject_format": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "Expression to generate or transform a unique username for the IdP user.",
			},
			"subject_filter": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Regular expression pattern used to filter untrusted IdP usernames.",
			},
			"issuer": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URI that identifies the issuer (IdP).",
			},
			"issuer_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates whether Okta uses the original Okta org domain URL, or a custom domain URL in the request to the IdP.",
			},
			"audience": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URI that identifies the target Okta IdP instance (SP)",
			},
			"kid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Key ID reference to the IdP's X.509 signature certificate.",
			},
		},
		Description: "Get a SAML IdP from Okta.",
	}
}

func dataSourceIdpSamlRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id := d.Get("id").(string)
	name := d.Get("name").(string)
	if id == "" && name == "" {
		return diag.Errorf("config must provide either 'id' or 'name' to retrieve the IdP")
	}
	var (
		err error
		idp *sdk.IdentityProvider
	)
	if id != "" {
		idp, err = getIdentityProviderByID(ctx, meta, id, Saml2Idp)
	} else {
		idp, err = getIdpByNameAndType(ctx, meta, name, Saml2Idp)
	}
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(idp.Id)
	_ = d.Set("name", idp.Name)
	_ = d.Set("type", idp.Type)
	e := idp.Protocol.Endpoints
	if e != nil {
		if e.Acs != nil {
			_ = d.Set("acs_binding", idp.Protocol.Endpoints.Acs.Binding)
			_ = d.Set("acs_type", idp.Protocol.Endpoints.Acs.Type)
		}
		if e.Sso != nil {
			_ = d.Set("sso_url", idp.Protocol.Endpoints.Sso.Url)
			_ = d.Set("sso_binding", idp.Protocol.Endpoints.Sso.Binding)
			_ = d.Set("sso_destination", idp.Protocol.Endpoints.Sso.Destination)
		}
	}
	t := idp.Protocol.Credentials.Trust
	if t != nil {
		_ = d.Set("kid", t.Kid)
		_ = d.Set("issuer", t.Issuer)
		_ = d.Set("audience", t.Audience)
	}
	if idp.Policy.Subject != nil {
		_ = d.Set("subject_filter", idp.Policy.Subject.Filter)
	}
	if idp.IssuerMode != "" {
		_ = d.Set("issuer_mode", idp.IssuerMode)
	}
	err = utils.SetNonPrimitives(d, map[string]interface{}{
		"subject_format": utils.ConvertStringSliceToSet(idp.Policy.Subject.Format),
	})
	if err != nil {
		return diag.Errorf("failed to set SAML identity provider properties: %v", err)
	}
	return nil
}
