package idaas

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	v6okta "github.com/okta/okta-sdk-golang/v6/okta"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func dataSourceAppOauth() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppOauthRead,
		Schema: utils.BuildSchema(skipUsersAndGroupsSchema, map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"label", "label_prefix"},
				Description:   "Id of application to retrieve, conflicts with label and label_prefix.",
			},
			"label": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id", "label_prefix"},
				Description: `The label of the app to retrieve, conflicts with
				label_prefix and id. Label uses the ?q=<label> query parameter exposed by
				Okta's List Apps API. The API will search both name and label using that
				query. Therefore similarly named and labeled apps may be returned in the query
				and have the unitended result of associating the wrong app with this data
				source. See:
				https://developer.okta.com/docs/reference/api/apps/#list-applications`,
			},
			"label_prefix": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id", "label"},
				Description: `Label prefix of the app to retrieve, conflicts with label and id. This will tell the
				provider to do a starts with query as opposed to an equals query.`,
			},
			"active_only": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Search only ACTIVE applications.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of OAuth application.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of application.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of application.",
			},
			"auto_submit_toolbar": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Display auto submit toolbar",
			},
			"hide_ios": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Do not display application icon on mobile app",
			},
			"hide_web": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Do not display application icon to users",
			},
			"grant_types": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "List of OAuth 2.0 grant types",
			},
			"response_types": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "List of OAuth 2.0 response type strings.",
			},
			"redirect_uris": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "List of URIs for use in the redirect-based flow.",
			},
			"post_logout_redirect_uris": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "List of URIs for redirection after logout",
			},
			"logo_uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URI that references a logo for the client.",
			},
			"login_uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URI that initiates login.",
			},
			"login_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of Idp-Initiated login that the client supports, if any",
			},
			"login_scopes": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of scopes to use for the request when 'login_mode' == OKTA",
			},
			"client_uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URI to a web page providing information about the client.",
			},
			"client_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "OAuth client ID",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "OAuth client secret",
			},
			"policy_uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URI to web page providing client policy document.",
			},
			"links": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Discoverable resources related to the app",
			},
			"wildcard_redirect": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates if the client is allowed to use wildcard matching of redirect_uris. Some valid values include: \"SUBDOMAIN\", \"DISABLED\".",
			},
			"dpop_bound_access_tokens": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates that the client application uses Demonstrating Proof-of-Possession (DPoP) for token requests. If true, the authorization server rejects token requests from this client that don't contain the DPoP header.",
			},
			"network": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Network restrictions for the application client.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connection": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The network connection type. Can be `ANYWHERE` or `ZONE`.",
						},
						"include": {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "IP zones to include when `connection` is `ZONE`. Can be `ALL_IP_ZONES` or specific zone IDs.",
						},
						"exclude": {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "IP zones to exclude when `connection` is `ZONE`. Can be `ALL_IP_ZONES` or specific zone IDs.",
						},
					},
				},
			},
		}),
		Description: "Get a OIDC application from Okta.",
	}
}

func dataSourceAppOauthRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	filters, err := getAppFilters(d)
	if err != nil {
		return diag.Errorf("invalid OAuth app filters: %v", err)
	}
	client := getOktaV6ClientFromMetadata(meta)

	var app *v6okta.OpenIdConnectApplication
	if filters.ID != "" {
		respApp, _, err := client.ApplicationAPI.GetApplication(ctx, filters.ID).Execute()
		if err != nil {
			return diag.Errorf("failed get app by ID: %v", err)
		}
		app, err = verifyOidcAppTypeV6(*respApp)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		req := client.ApplicationAPI.ListApplications(ctx).Limit(1)
		if q := filters.GetQ(); q != "" {
			req = req.Q(q)
		}
		if filters.Status != "" {
			req = req.Filter(filters.Status)
		}
		appList, _, err := req.Execute()
		if err != nil {
			return diag.Errorf("failed to list OAuth apps: %v", err)
		}
		if len(appList) < 1 {
			return diag.Errorf("no OAuth application found with provided filter: %s", filters)
		}
		app, err = verifyOidcAppTypeV6(appList[0])
		if err != nil {
			return diag.FromErr(err)
		}
		if filters.Label != "" && app.GetLabel() != filters.Label {
			return diag.Errorf("no OAuth application found with the provided label: %s", filters.Label)
		}
		logger(meta).Info("found multiple OAuth applications with the criteria supplied, using the first one, sorted by creation date")
	}

	d.SetId(app.GetId())
	_ = d.Set("label", app.GetLabel())
	_ = d.Set("name", app.GetName())
	_ = d.Set("status", app.GetStatus())

	visibility := app.GetVisibility()
	_ = d.Set("auto_submit_toolbar", visibility.GetAutoSubmitToolbar())
	hide := visibility.GetHide()
	_ = d.Set("hide_ios", hide.GetIOS())
	_ = d.Set("hide_web", hide.GetWeb())

	respTypes := []string{}
	grantTypes := []string{}
	redirectUris := []string{}
	postLogoutRedirectUris := []string{}

	settings := app.GetSettings()
	oauthClient := settings.OauthClient
	if oauthClient != nil {
		_ = d.Set("type", oauthClient.GetApplicationType())
		_ = d.Set("client_uri", oauthClient.GetClientUri())
		_ = d.Set("logo_uri", oauthClient.GetLogoUri())
		_ = d.Set("login_uri", oauthClient.GetInitiateLoginUri())

		credentials := app.GetCredentials()
		credentialsOauthClient := credentials.GetOauthClient()
		_ = d.Set("client_id", credentialsOauthClient.GetClientId())

		secret, err := getCurrentlyActiveClientSecret(ctx, meta, app.GetId())
		if err != nil {
			return diag.Errorf("failed to fetch OAuth client secret: %v", err)
		}
		_ = d.Set("client_secret", secret)

		_ = d.Set("policy_uri", oauthClient.GetPolicyUri())
		_ = d.Set("wildcard_redirect", oauthClient.GetWildcardRedirect())
		_ = d.Set("dpop_bound_access_tokens", oauthClient.GetDpopBoundAccessTokens())

		respTypes = append(respTypes, oauthClient.ResponseTypes...)
		grantTypes = append(grantTypes, oauthClient.GrantTypes...)
		redirectUris = append(redirectUris, oauthClient.RedirectUris...)
		postLogoutRedirectUris = append(postLogoutRedirectUris, oauthClient.PostLogoutRedirectUris...)
	}

	aggMap := map[string]interface{}{
		"redirect_uris":             utils.ConvertStringSliceToSet(redirectUris),
		"response_types":            utils.ConvertStringSliceToSet(respTypes),
		"grant_types":               utils.ConvertStringSliceToSet(grantTypes),
		"post_logout_redirect_uris": utils.ConvertStringSliceToSet(postLogoutRedirectUris),
	}
	if oauthClient != nil {
		if idpLogin, ok := oauthClient.GetIdpInitiatedLoginOk(); ok && idpLogin != nil {
			_ = d.Set("login_mode", idpLogin.GetMode())
			aggMap["login_scopes"] = utils.ConvertStringSliceToSet(idpLogin.DefaultScope)
		}
	}

	err = utils.SetNonPrimitives(d, aggMap)
	if err != nil {
		return diag.Errorf("failed to set OAuth application properties: %v", err)
	}
	if oauthClient != nil {
		if network, ok := oauthClient.GetNetworkOk(); ok && network != nil {
			networkMap := map[string]interface{}{
				"connection": network.GetConnection(),
				"include":    utils.ConvertStringSliceToSet(network.GetInclude()),
				"exclude":    utils.ConvertStringSliceToSet(network.GetExclude()),
			}
			if err := utils.SetNonPrimitives(d, map[string]interface{}{"network": []interface{}{networkMap}}); err != nil {
				return diag.Errorf("failed to set OAuth application network properties: %v", err)
			}
		}
	}
	p, _ := json.Marshal(app.Links)
	_ = d.Set("links", string(p))
	return nil
}

// getCurrentlyActiveClientSecret See: https://developer.okta.com/docs/reference/api/apps/#list-client-secrets
func getCurrentlyActiveClientSecret(ctx context.Context, meta interface{}, appID string) (string, error) {
	secrets, _, err := getOktaV6ClientFromMetadata(meta).ApplicationSSOPublicKeysAPI.ListOAuth2ClientSecrets(ctx, appID).Execute()
	if err != nil {
		return "", err
	}

	// There can only be two client secrets. Regardless, choose the latest created active secret.
	// v6 reports `created` as an RFC3339 string rather than a time.Time, so parse before comparing;
	// an unparseable timestamp sorts as the zero time and simply loses to any parseable one.
	var secretValue string
	var newest time.Time
	found := false
	for i := range secrets {
		s := secrets[i]
		if s.GetStatus() != StatusActive {
			continue
		}
		created, err := time.Parse(time.RFC3339, s.GetCreated())
		if err != nil {
			created = time.Time{}
		}
		if !found || created.After(newest) {
			secretValue = s.GetClientSecret()
			newest = created
			found = true
		}
	}

	return secretValue, nil
}
