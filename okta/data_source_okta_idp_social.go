package okta

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/okta/terraform-provider-okta/sdk/query"
)

func dataSourceIdpSocial() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdpSocialRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "The name of the social idp to retrieve, conflicts with `id`.",
				ConflictsWith: []string{"id"},
			},
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "The id of the social idp to retrieve, conflicts with `name`.",
				ConflictsWith: []string{"name"},
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the IdP.",
			},
			"account_link_action": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the account linking action for an IdP user.",
			},
			"account_link_group_include": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "Group memberships to determine link candidates.",
			},
			"provisioning_action": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Provisioning action for an IdP user during authentication.",
			},
			"deprovisioned_action": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Action for a previously deprovisioned IdP user during authentication.",
			},
			"suspended_action": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Action for a previously suspended IdP user during authentication.",
			},
			"groups_action": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Provisioning action for IdP user's group memberships.",
			},
			"groups_attribute": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IdP user profile attribute name for an array value that contains group memberships.",
			},
			"groups_assignment": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "List of Okta Group IDs.",
			},
			"groups_filter": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "Whitelist of Okta Group identifiers.",
			},
			"username_template": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Okta EL Expression to generate or transform a unique username for the IdP user.",
			},
			"subject_match_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Determines the Okta user profile attribute match conditions for account linking and authentication of the transformed IdP username.",
			},
			"subject_match_attribute": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Okta user profile attribute for matching transformed IdP username.",
			},
			"profile_master": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Determines if the IdP should act as a source of truth for user profile attributes.",
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
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of Social IdP. See API docs [Identity Provider Type](https://developer.okta.com/docs/reference/api/idps/#identity-provider-type)",
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
			"max_clock_skew": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Maximum allowable clock-skew when processing messages from the IdP.",
			},
			"issuer_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates whether Okta uses the original Okta org domain URL, or a custom domain URL.",
			},
		},
		Description: "Get a social IdP from Okta.",
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
		idp *sdk.IdentityProvider
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
	if idp.Policy.MaxClockSkewPtr != nil {
		_ = d.Set("max_clock_skew", idp.Policy.MaxClockSkewPtr)
	}
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
		"scopes": convertStringSliceToSet(idp.Protocol.Scopes),
	}
	if idp.Policy.AccountLink != nil {
		_ = d.Set("account_link_action", idp.Policy.AccountLink.Action)
		if idp.Policy.AccountLink.Filter != nil {
			setMap["account_link_group_include"] = convertStringSliceToSet(idp.Policy.AccountLink.Filter.Groups.Include)
		}
	}
	err = setNonPrimitives(d, setMap)
	if err != nil {
		return diag.Errorf("failed to set social identity provider properties: %v", err)
	}
	return nil
}

func getSocialIdPByName(ctx context.Context, m interface{}, name string) (*sdk.IdentityProvider, error) {
	idps, _, err := getOktaClientFromMetadata(m).IdentityProvider.
		ListIdentityProviders(ctx, &query.Params{Q: name, Limit: defaultPaginationLimit})
	if err != nil {
		return nil, fmt.Errorf("failed to get social identity provider with name '%s': %w", name, err)
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
