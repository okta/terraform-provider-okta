package idaas

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/okta/terraform-provider-okta/sdk/query"
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
			"authentication_policy": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the app's authentication policy",
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
	var app *sdk.OpenIdConnectApplication
	if filters.ID != "" {
		respApp, _, err := getOktaClientFromMetadata(meta).Application.GetApplication(ctx, filters.ID, sdk.NewOpenIdConnectApplication(), nil)
		if err != nil {
			return diag.Errorf("failed get app by ID: %v", err)
		}
		app = respApp.(*sdk.OpenIdConnectApplication)
	} else {
		re := getOktaClientFromMetadata(meta).GetRequestExecutor()
		qp := &query.Params{Limit: 1, Filter: filters.Status, Q: filters.GetQ()}
		req, err := re.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/apps%s", qp.String()), nil)
		if err != nil {
			return diag.Errorf("failed to list OAuth apps: %v", err)
		}
		var appList []*sdk.OpenIdConnectApplication
		_, err = re.Do(ctx, req, &appList)
		if err != nil {
			return diag.Errorf("failed to list OAuth apps: %v", err)
		}
		if len(appList) < 1 {
			return diag.Errorf("no OAuth application found with provided filter: %s", filters)
		}
		if filters.Label != "" && appList[0].Label != filters.Label {
			return diag.Errorf("no OAuth application found with the provided label: %s", filters.Label)
		}
		logger(meta).Info("found multiple OAuth applications with the criteria supplied, using the first one, sorted by creation date")
		app = appList[0]
	}

	d.SetId(app.Id)
	_ = d.Set("label", app.Label)
	_ = d.Set("name", app.Name)
	_ = d.Set("status", app.Status)
	_ = d.Set("auto_submit_toolbar", app.Visibility.AutoSubmitToolbar)
	_ = d.Set("hide_ios", app.Visibility.Hide.IOS)
	_ = d.Set("hide_web", app.Visibility.Hide.Web)

	respTypes := []string{}
	grantTypes := []string{}
	redirectUris := []string{}
	postLogoutRedirectUris := []string{}

	if app.Settings.OauthClient != nil {
		_ = d.Set("type", app.Settings.OauthClient.ApplicationType)
		_ = d.Set("client_uri", app.Settings.OauthClient.ClientUri)
		_ = d.Set("logo_uri", app.Settings.OauthClient.LogoUri)
		_ = d.Set("login_uri", app.Settings.OauthClient.InitiateLoginUri)
		_ = d.Set("client_id", app.Credentials.OauthClient.ClientId)

		secret, err := getCurrentlyActiveClientSecret(ctx, meta, app.Id)
		if err != nil {
			return diag.Errorf("failed to fetch OAuth client secret: %v", err)
		}
		_ = d.Set("client_secret", secret)

		_ = d.Set("policy_uri", app.Settings.OauthClient.PolicyUri)
		_ = d.Set("wildcard_redirect", app.Settings.OauthClient.WildcardRedirect)
		for i := range app.Settings.OauthClient.ResponseTypes {
			respTypes = append(respTypes, string(*app.Settings.OauthClient.ResponseTypes[i]))
		}
		for i := range app.Settings.OauthClient.GrantTypes {
			grantTypes = append(grantTypes, string(*app.Settings.OauthClient.GrantTypes[i]))
		}
		redirectUris = append(redirectUris, app.Settings.OauthClient.RedirectUris...)
		postLogoutRedirectUris = append(postLogoutRedirectUris, app.Settings.OauthClient.PostLogoutRedirectUris...)
	}

	aggMap := map[string]interface{}{
		"redirect_uris":             utils.ConvertStringSliceToSet(redirectUris),
		"response_types":            utils.ConvertStringSliceToSet(respTypes),
		"grant_types":               utils.ConvertStringSliceToSet(grantTypes),
		"post_logout_redirect_uris": utils.ConvertStringSliceToSet(postLogoutRedirectUris),
	}
	if app.Settings.OauthClient != nil &&
		app.Settings.OauthClient.IdpInitiatedLogin != nil {
		_ = d.Set("login_mode", app.Settings.OauthClient.IdpInitiatedLogin.Mode)
		aggMap["login_scopes"] = utils.ConvertStringSliceToSet(app.Settings.OauthClient.IdpInitiatedLogin.DefaultScope)
	}

	err = utils.SetNonPrimitives(d, aggMap)
	if err != nil {
		return diag.Errorf("failed to set OAuth application properties: %v", err)
	}
	p, _ := json.Marshal(app.Links)
	_ = d.Set("links", string(p))
	setAuthenticationPolicy(ctx, meta, d, app.Links)
	return nil
}

// getCurrentlyActiveClientSecret See: https://developer.okta.com/docs/reference/api/apps/#list-client-secrets
func getCurrentlyActiveClientSecret(ctx context.Context, meta interface{}, appId string) (string, error) {
	secrets, _, err := getOktaClientFromMetadata(meta).Application.ListClientSecretsForApplication(ctx, appId)
	if err != nil {
		return "", err
	}

	// There can only be two client secrets. Regardless, choose the latest created active secret.
	var secretValue string
	var secret *sdk.ClientSecret
	for _, s := range secrets {
		if secret == nil && s.Status == "ACTIVE" {
			secret = s
		}
		if secret != nil && s.Status == "ACTIVE" && secret.Created.Before(*s.Created) {
			secret = s
		}
	}
	if secret != nil {
		secretValue = secret.ClientSecret
	}

	return secretValue, nil
}
