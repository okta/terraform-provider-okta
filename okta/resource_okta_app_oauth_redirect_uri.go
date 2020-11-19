package okta

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func resourceAppOAuthRedirectURI() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppOAuthRedirectURICreate,
		Read:   resourceAppOAuthRedirectURIRead,
		Update: resourceAppOAuthRedirectURIUpdate,
		Delete: resourceAppOAuthRedirectURIDelete,
		Exists: resourceAppOAuthRedirectURIExists,
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
				Description: "Redirect URI to append to Okta OIDC application.",
			},
		},
	}
}

func resourceAppOAuthRedirectURIExists(d *schema.ResourceData, m interface{}) (bool, error) {
	app := okta.NewOpenIdConnectApplication()
	err := fetchAppByID(d.Get("app_id").(string), m, app)
	return err == nil && app.Id != "" && contains(app.Settings.OauthClient.RedirectUris, d.Id()), err
}

func resourceAppOAuthRedirectURICreate(d *schema.ResourceData, m interface{}) error {
	if err := appendRedirectURI(d, m); err != nil {
		return err
	}
	d.SetId(d.Get("uri").(string))

	return resourceAppOAuthRedirectURIRead(d, m)
}

// read does nothing due to the nature of this resource
func resourceAppOAuthRedirectURIRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceAppOAuthRedirectURIUpdate(d *schema.ResourceData, m interface{}) error {
	if err := appendRedirectURI(d, m); err != nil {
		return err
	}
	// Normally not advisable, but ForceNew generated unnecessary calls
	d.SetId(d.Get("uri").(string))

	return resourceAppOAuthRedirectURIRead(d, m)
}

func resourceAppOAuthRedirectURIDelete(d *schema.ResourceData, m interface{}) error {
	appID := d.Get("app_id").(string)
	app := okta.NewOpenIdConnectApplication()
	// Should never hit a 404 due to exists function
	if err := fetchAppByID(appID, m, app); err != nil {
		return err
	}
	app.Settings.OauthClient.RedirectUris = remove(app.Settings.OauthClient.RedirectUris, d.Id())

	return updateAppByID(appID, m, app)
}

func appendRedirectURI(d *schema.ResourceData, m interface{}) error {
	appID := d.Get("app_id").(string)
	app := okta.NewOpenIdConnectApplication()
	// Should never hit a 404 due to exists function
	if err := fetchAppByID(appID, m, app); err != nil {
		return err
	}
	uri := d.Get("uri").(string)
	app.Settings.OauthClient.RedirectUris = append(app.Settings.OauthClient.RedirectUris, uri)
	return updateAppByID(appID, m, app)
}
