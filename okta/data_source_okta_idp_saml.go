package okta

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
)

const (
	saml2Idp = "SAML2"
)

func dataSourceIdpSaml() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIdpSamlRead,

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

func dataSourceIdpSamlRead(d *schema.ResourceData, m interface{}) error {
	var (
		err error
		idp *sdk.SAMLIdentityProvider
	)

	id := d.Get("id").(string)
	name := d.Get("name").(string)

	if id == "" && name == "" {
		return errors.New("config must provide an id or name to retrieve the IdP")
	}

	if id != "" {
		idp, err = getIdentityProviderByID(m, id)
	} else {
		idp, err = getIdpByName(m, name)
	}

	if err != nil {
		return err
	} else if idp == nil {
		return fmt.Errorf("failed to find IdP via filter, id %s, name %s", id, name)
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

	return setNonPrimitives(d, map[string]interface{}{
		"subject_format": convertStringSetToInterface(idp.Policy.Subject.Format),
	})
}

func getIdentityProviderByID(m interface{}, id string) (*sdk.SAMLIdentityProvider, error) {
	var idp sdk.SAMLIdentityProvider
	client := getSupplementFromMetadata(m)
	_, resp, err := client.GetIdentityProvider(id, &idp)

	return &idp, responseErr(resp, err)
}

func getIdpByName(m interface{}, label string) (*sdk.SAMLIdentityProvider, error) {
	var idps []*sdk.SAMLIdentityProvider
	queryParams := query.Params{Limit: 1, Q: label}
	client := getSupplementFromMetadata(m)
	_, resp, err := client.ListIdentityProviders(&idps, &queryParams)

	if len(idps) > 0 {
		return idps[0], nil
	}

	return nil, responseErr(resp, err)
}
