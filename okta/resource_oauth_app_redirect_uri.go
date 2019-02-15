package okta

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/okta/okta-sdk-golang/okta"
)

func resourceOAuthAppRedirectUri() *schema.Resource {
	return &schema.Resource{
		Create: resourceOAuthAppRedirectUriCreate,
		Read:   resourceOAuthAppRedirectUriRead,
		Update: resourceOAuthAppRedirectUriUpdate,
		Delete: resourceOAuthAppRedirectUriDelete,
		Exists: resourceOAuthAppRedirectUriExists,
		// The id for this is the uri
		Importer: createCustomNestedResourceImporter([]string{"app_id", "id"}, "Expecting the following format: <app_id>/<uri>"),

		Schema: map[string]*schema.Schema{
			"app_id": &schema.Schema{
				Required: true,
				Type:     schema.TypeString,
				ForceNew: true,
			},
			"uri": &schema.Schema{
				Required:    true,
				Type:        schema.TypeString,
				Description: "Redirect URI to append to Okta OIDC application.",
			},
		},
	}
}

func resourceOAuthAppRedirectUriExists(d *schema.ResourceData, m interface{}) (bool, error) {
	app := okta.NewOpenIdConnectApplication()
	err := fetchAppById(d.Get("app_id").(string), m, app)
	return err == nil && app.Id != "" && contains(app.Settings.OauthClient.RedirectUris, d.Id()), err
}

func resourceOAuthAppRedirectUriCreate(d *schema.ResourceData, m interface{}) error {
	if err := appendRedirectUri(d, m); err != nil {
		return err
	}
	d.SetId(d.Get("uri").(string))

	return resourceOAuthAppRedirectUriRead(d, m)
}

// read does nothing due to the nature of this resource
func resourceOAuthAppRedirectUriRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceOAuthAppRedirectUriUpdate(d *schema.ResourceData, m interface{}) error {
	if err := appendRedirectUri(d, m); err != nil {
		return err
	}
	// Normally not advisable, but ForceNew generated unnecessary calls
	d.SetId(d.Get("uri").(string))

	return resourceOAuthAppRedirectUriRead(d, m)
}

func resourceOAuthAppRedirectUriDelete(d *schema.ResourceData, m interface{}) error {
	appId := d.Get("app_id").(string)
	app := okta.NewOpenIdConnectApplication()
	// Should never hit a 404 due to exists function
	if err := fetchAppById(appId, m, app); err != nil {
		return err
	}
	app.Settings.OauthClient.RedirectUris = remove(app.Settings.OauthClient.RedirectUris, d.Id())

	return updateAppById(appId, m, app)
}

func appendRedirectUri(d *schema.ResourceData, m interface{}) error {
	appId := d.Get("app_id").(string)
	app := okta.NewOpenIdConnectApplication()
	// Should never hit a 404 due to exists function
	if err := fetchAppById(appId, m, app); err != nil {
		return err
	}
	uri := d.Get("uri").(string)
	app.Settings.OauthClient.RedirectUris = append(app.Settings.OauthClient.RedirectUris, uri)
	return updateAppById(appId, m, app)
}
