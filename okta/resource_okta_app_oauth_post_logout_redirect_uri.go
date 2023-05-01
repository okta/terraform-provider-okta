package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourceAppOAuthPostLogoutRedirectURI() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppOAuthPostLogoutRedirectURICreate,
		ReadContext:   resourceFuncNoOp,
		UpdateContext: resourceAppOAuthPostLogoutRedirectURIUpdate,
		DeleteContext: resourceAppOAuthPostLogoutRedirectURIDelete,
		// The id for this is the uri
		Importer: createCustomNestedResourceImporter([]string{"app_id", "id"}, "Expecting the following format: <app_id>/<uri>"),
		Schema: map[string]*schema.Schema{
			"app_id": {
				Required: true,
				Type:     schema.TypeString,
				ForceNew: true,
			},
			"uri": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Post Logout Redirect URI to append to Okta OIDC application.",
			},
		},
	}
}

func resourceAppOAuthPostLogoutRedirectURICreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := appendPostLogoutRedirectURI(ctx, d, m)
	if err != nil {
		return diag.Errorf("failed to create post logout redirect URI: %v", err)
	}
	d.SetId(d.Get("uri").(string))
	return resourceFuncNoOp(ctx, d, m)
}

func resourceAppOAuthPostLogoutRedirectURIUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if err := appendPostLogoutRedirectURI(ctx, d, m); err != nil {
		return diag.Errorf("failed to update post logout redirect URI: %v", err)
	}
	// Normally not advisable, but ForceNew generated unnecessary calls
	d.SetId(d.Get("uri").(string))
	return resourceFuncNoOp(ctx, d, m)
}

func resourceAppOAuthPostLogoutRedirectURIDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	if !contains(app.Settings.OauthClient.PostLogoutRedirectUris, d.Id()) {
		logger(m).Info(fmt.Sprintf("application with appID %s does not have post logout redirect URI %s", appID, d.Id()))
		return nil
	}
	app.Settings.OauthClient.PostLogoutRedirectUris = remove(app.Settings.OauthClient.PostLogoutRedirectUris, d.Id())
	err = updateAppByID(ctx, appID, m, app)
	if err != nil {
		return diag.Errorf("failed to delete post logout redirect URI: %v", err)
	}
	return nil
}

func appendPostLogoutRedirectURI(ctx context.Context, d *schema.ResourceData, m interface{}) error {
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
	if contains(app.Settings.OauthClient.PostLogoutRedirectUris, d.Id()) {
		logger(m).Info(fmt.Sprintf("application with appID %s already has post logout redirect URI %s", appID, d.Id()))
		return nil
	}
	uri := d.Get("uri").(string)
	app.Settings.OauthClient.PostLogoutRedirectUris = append(app.Settings.OauthClient.PostLogoutRedirectUris, uri)
	return updateAppByID(ctx, appID, m, app)
}
