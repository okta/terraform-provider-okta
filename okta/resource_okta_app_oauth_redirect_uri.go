package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourceAppOAuthRedirectURI() *schema.Resource {
	return &schema.Resource{
		// NOTE: These CRUD contexts are flexible for use with resources
		// okta_app_oauth_redirect_uri and
		// okta_app_oauth_post_logout_redirect_uri
		CreateContext: resourceAppOAuthRedirectURICreate("okta_app_oauth_redirect_uri"),
		ReadContext:   resourceAppOAuthRedirectURIRead("okta_app_oauth_redirect_uri"),
		UpdateContext: resourceAppOAuthRedirectURIUpdate("okta_app_oauth_redirect_uri"),
		DeleteContext: resourceAppOAuthRedirectURIDelete("okta_app_oauth_redirect_uri"),
		Importer:      createCustomNestedResourceImporter([]string{"app_id", "id"}, "Expecting the following format: <app_id>/<uri>"),
		Schema: map[string]*schema.Schema{
			"app_id": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "OAuth application ID.",
			},
			"uri": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Redirect URI to append to Okta OIDC application.",
			},
		},
	}
}

func resourceAppOAuthRedirectURICreate(kind string) func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		d.SetId(d.Get("uri").(string))
		err := appendRedirectURI(ctx, d, m, kind)
		if err != nil {
			return diag.Errorf("failed to create %q: %v", kind, err)
		}
		return resourceFuncNoOp(ctx, d, m)
	}
}

func resourceAppOAuthRedirectURIRead(kind string) func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		aid, ok := d.GetOk("app_id")
		if !ok || aid.(string) == "" {
			return diag.Errorf("app_id not set on resource")

		}
		appID := aid.(string)

		app := sdk.NewOpenIdConnectApplication()
		if err := fetchAppByID(ctx, appID, m, app); err != nil {
			return diag.Errorf("application %q not found: %q", appID, err)
		}
		if app.Id == "" {
			return diag.Errorf("application with id %s does not exist", appID)
		}

		switch kind {
		case "okta_app_oauth_redirect_uri":
			if !contains(app.Settings.OauthClient.RedirectUris, d.Id()) {
				return diag.Errorf("application %q does not have redirect uri %q", appID, d.Id())
			}
		case "okta_app_oauth_post_logout_redirect_uri":
			if !contains(app.Settings.OauthClient.PostLogoutRedirectUris, d.Id()) {
				return diag.Errorf("application %q does not have post logout redirect uri %q", appID, d.Id())
			}
		default:
			return diag.Errorf("unknown resource type %q", kind)
		}

		return resourceFuncNoOp(ctx, d, m)
	}
}

func resourceAppOAuthRedirectURIUpdate(kind string) func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		if d.HasChange("app_id") {
			return diag.Errorf("it is invalid to change the app_id of this resource once set")
		}
		if !d.HasChange("uri") {
			return resourceFuncNoOp(ctx, d, m)
		}

		o, n := d.GetChange("uri")
		oldURI := o.(string)
		newURI := n.(string)

		if newURI == "" {
			return diag.Errorf("it is invalid to change uri to a blank value")
		}
		if err := changeOauthAppRedirectURI(ctx, d, m, kind, oldURI, newURI); err != nil {
			return diag.Errorf("failed to update %q's uri: %v", kind, err)
		}
		d.SetId(newURI)
		return resourceFuncNoOp(ctx, d, m)
	}
}

func resourceAppOAuthRedirectURIDelete(kind string) func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		appID := d.Get("app_id").(string)

		oktaMutexKV.Lock(appID)
		defer oktaMutexKV.Unlock(appID)

		app := sdk.NewOpenIdConnectApplication()
		err := fetchAppByID(ctx, appID, m, app)
		if err != nil {
			return diag.Errorf("failed to get application: %v", err)
		}
		if app.Id == "" {
			return diag.Errorf("application with id %s does not exist", appID)
		}

		switch kind {
		case "okta_app_oauth_redirect_uri":
			if !contains(app.Settings.OauthClient.RedirectUris, d.Id()) {
				logger(m).Info(fmt.Sprintf("application with appID %s does not have redirect URI %s", appID, d.Id()))
				return resourceFuncNoOp(ctx, d, m)
			}
			app.Settings.OauthClient.RedirectUris = remove(app.Settings.OauthClient.RedirectUris, d.Id())
		case "okta_app_oauth_post_logout_redirect_uri":
			if !contains(app.Settings.OauthClient.PostLogoutRedirectUris, d.Id()) {
				logger(m).Info(fmt.Sprintf("application with appID %s does not have post logout redirect URI %s", appID, d.Id()))
				return resourceFuncNoOp(ctx, d, m)
			}
			app.Settings.OauthClient.PostLogoutRedirectUris = remove(app.Settings.OauthClient.PostLogoutRedirectUris, d.Id())
		default:
			return diag.Errorf("unknown resource type %q", kind)
		}

		err = updateAppByID(ctx, appID, m, app)
		if err != nil {
			return diag.Errorf("failed to delete uri for %q: %v", kind, err)
		}
		return resourceFuncNoOp(ctx, d, m)
	}
}

func appendRedirectURI(ctx context.Context, d *schema.ResourceData, m interface{}, uriType string) error {
	appID := d.Get("app_id").(string)

	oktaMutexKV.Lock(appID)
	defer oktaMutexKV.Unlock(appID)

	app := sdk.NewOpenIdConnectApplication()
	if err := fetchAppByID(ctx, appID, m, app); err != nil {
		return err
	}
	if app.Id == "" {
		return fmt.Errorf("application with id %s does not exist", appID)
	}

	uri := d.Get("uri").(string)
	switch uriType {
	case "okta_app_oauth_redirect_uri":
		if contains(app.Settings.OauthClient.RedirectUris, d.Id()) {
			logger(m).Info(fmt.Sprintf("application with appID %s already has redirect URI %s", appID, d.Id()))
			return nil
		}
		app.Settings.OauthClient.RedirectUris = append(app.Settings.OauthClient.RedirectUris, uri)
	case "okta_app_oauth_post_logout_redirect_uri":
		if contains(app.Settings.OauthClient.PostLogoutRedirectUris, d.Id()) {
			logger(m).Info(fmt.Sprintf("application with appID %s already has post logout redirect URI %s", appID, d.Id()))
			return nil
		}
		app.Settings.OauthClient.PostLogoutRedirectUris = append(app.Settings.OauthClient.PostLogoutRedirectUris, d.Id())
	default:
		return fmt.Errorf("unknown resource type %q", uriType)
	}

	return updateAppByID(ctx, appID, m, app)
}

// changeOauthAppRedirectURI will update the redirect uris on the given
// application. toRemoveURI will be removed if it exists as a redirect uri on
// the app and add toAddURI if it doesn't already exist as a redirect URI on the
// app. Blank values are ignored. This function is intended for resources
// okta_app_oauth_redirect_uri and okta_app_oauth_post_logout_redirect_uri
func changeOauthAppRedirectURI(ctx context.Context, d *schema.ResourceData, m interface{}, uriType, toRemoveURI, toAddURI string) error {
	appID := d.Get("app_id").(string)

	oktaMutexKV.Lock(appID)
	defer oktaMutexKV.Unlock(appID)

	app := sdk.NewOpenIdConnectApplication()
	if err := fetchAppByID(ctx, appID, m, app); err != nil {
		return err
	}
	if app.Id == "" {
		return fmt.Errorf("application with id %s does not exist", appID)
	}

	switch uriType {
	case "okta_app_oauth_redirect_uri":
		app.Settings.OauthClient.RedirectUris = appendUnique(app.Settings.OauthClient.RedirectUris, toAddURI)
		app.Settings.OauthClient.RedirectUris = remove(app.Settings.OauthClient.RedirectUris, toRemoveURI)
	case "okta_app_oauth_post_logout_redirect_uri":
		app.Settings.OauthClient.PostLogoutRedirectUris = appendUnique(app.Settings.OauthClient.PostLogoutRedirectUris, toAddURI)
		app.Settings.OauthClient.PostLogoutRedirectUris = remove(app.Settings.OauthClient.PostLogoutRedirectUris, toRemoveURI)
	default:
		return fmt.Errorf("unknown resource type %q", uriType)
	}

	return updateAppByID(ctx, appID, m, app)
}
