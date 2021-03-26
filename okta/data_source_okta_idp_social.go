package okta

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

func dataSourceIdpSocial() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdpSocialRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "name of the IdP",
				ConflictsWith: []string{"id"},
			},
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "ID of the IdP",
				ConflictsWith: []string{"name"},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_link_action": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_link_group_include": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"provisioning_action": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deprovisioned_action": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"suspended_action": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"groups_action": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"groups_attribute": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"groups_assignment": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"groups_filter": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"username_template": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subject_match_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subject_match_attribute": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"profile_master": {
				Type:     schema.TypeBool,
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
			"type": {
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
			"max_clock_skew": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"issuer_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceIdpSocialRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Get("id").(string)
	name := d.Get("name").(string)
	if id == "" && name == "" {
		return diag.Errorf("config must provide either 'id' or 'name' to retrieve the social IdP")
	}
	var (
		err error
		idp *okta.IdentityProvider
	)
	if id != "" {
		idp, _, err = getOktaClientFromMetadata(m).IdentityProvider.GetIdentityProvider(ctx, id)
		if err != nil {
			return diag.Errorf("failed to get social identity provider with id '%s': %v", id, err)
		}
		if !contains([]string{"APPLE", "FACEBOOK", "LINKEDIN", "MICROSOFT", "GOOGLE"}, idp.Type) {
			return diag.Errorf("social identity provider with id '%s' does not exist", id)
		}
	} else {
		idp, err = getSocialIdPByName(ctx, m, name)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId(idp.Id)
	_ = d.Set("name", idp.Name)
	_ = d.Set("status", idp.Status)
	_ = d.Set("max_clock_skew", idp.Policy.MaxClockSkew)
	_ = d.Set("provisioning_action", idp.Policy.Provisioning.Action)
	_ = d.Set("deprovisioned_action", idp.Policy.Provisioning.Conditions.Deprovisioned.Action)
	_ = d.Set("suspended_action", idp.Policy.Provisioning.Conditions.Suspended.Action)
	_ = d.Set("profile_master", idp.Policy.Provisioning.ProfileMaster)
	_ = d.Set("subject_match_type", idp.Policy.Subject.MatchType)
	_ = d.Set("subject_match_attribute", idp.Policy.Subject.MatchAttribute)
	_ = d.Set("username_template", idp.Policy.Subject.UserNameTemplate.Template)
	_ = d.Set("client_id", idp.Protocol.Credentials.Client.ClientId)
	_ = d.Set("client_secret", idp.Protocol.Credentials.Client.ClientSecret)
	_ = d.Set("issuer_mode", idp.IssuerMode)
	_ = d.Set("protocol_type", idp.Protocol.Type)
	_ = d.Set("type", idp.Type)
	syncEndpoint("authorization", idp.Protocol.Endpoints.Authorization, d)
	syncEndpoint("token", idp.Protocol.Endpoints.Authorization, d)

	err = syncGroupActions(d, idp.Policy.Provisioning.Groups)
	if err != nil {
		return diag.Errorf("failed to set social identity provider properties: %v", err)
	}
	setMap := map[string]interface{}{
		"scopes": convertStringSetToInterface(idp.Protocol.Scopes),
	}
	if idp.Policy.AccountLink != nil {
		_ = d.Set("account_link_action", idp.Policy.AccountLink.Action)
		if idp.Policy.AccountLink.Filter != nil {
			setMap["account_link_group_include"] = convertStringSetToInterface(idp.Policy.AccountLink.Filter.Groups.Include)
		}
	}
	err = setNonPrimitives(d, setMap)
	if err != nil {
		return diag.Errorf("failed to set social identity provider properties: %v", err)
	}
	return nil
}

func getSocialIdPByName(ctx context.Context, m interface{}, name string) (*okta.IdentityProvider, error) {
	idps, _, err := getOktaClientFromMetadata(m).IdentityProvider.
		ListIdentityProviders(ctx, &query.Params{Q: name, Limit: 200})
	if err != nil {
		return nil, fmt.Errorf("failed to get social identity provider with name '%s': %v", name, err)
	}
	if len(idps) < 1 || !contains([]string{"APPLE", "FACEBOOK", "LINKEDIN", "MICROSOFT", "GOOGLE"}, idps[0].Type) {
		return nil, fmt.Errorf("social identity provider with name '%s' does not exist", name)
	}
	k := 0
	for i, n := range idps {
		if idps[i].Name == name {
			return idps[i], nil
		}
		if strings.Contains(idps[i].Name, name) {
			if i != k {
				idps[k] = n
			}
			k++
		}
	}
	idps = idps[:k]
	if len(idps) == 0 {
		return nil, fmt.Errorf("social identity provider with name '%s' does not exist", name)
	}
	if len(idps) > 1 {
		logger(m).Warn(fmt.Sprintf("found multiple social IdPs with name '%s': "+
			"using the first one which may only be a partial match", name))
	}
	return idps[0], nil
}
