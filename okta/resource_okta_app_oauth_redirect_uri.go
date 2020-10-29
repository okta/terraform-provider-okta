package okta

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func resourceAppOAuthRedirectUri() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppOAuthRedirectUriCreate,
		Read:   resourceAppOAuthRedirectUriRead,
		Update: resourceAppOAuthRedirectUriUpdate,
		Delete: resourceAppOAuthRedirectUriDelete,
		Exists: resourceAppOAuthRedirectUriExists,
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

func resourceAppOAuthRedirectUriExists(d *schema.ResourceData, m interface{}) (bool, error) {
	app := okta.NewOpenIdConnectApplication()
	err := fetchAppById(d.Get("app_id").(string), m, app)
	return err == nil && app.Id != "" && contains(app.Settings.OauthClient.RedirectUris, d.Id()), err
}

func resourceAppOAuthRedirectUriCreate(d *schema.ResourceData, m interface{}) error {
	if err := appendRedirectUri(d, m); err != nil {
		return err
	}
	d.SetId(d.Get("uri").(string))

	return resourceAppOAuthRedirectUriRead(d, m)
}

// read does nothing due to the nature of this resource
func resourceAppOAuthRedirectUriRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceAppOAuthRedirectUriUpdate(d *schema.ResourceData, m interface{}) error {
	if err := appendRedirectUri(d, m); err != nil {
		return err
	}
	// Normally not advisable, but ForceNew generated unnecessary calls
	d.SetId(d.Get("uri").(string))

	return resourceAppOAuthRedirectUriRead(d, m)
}

func resourceAppOAuthRedirectUriDelete(d *schema.ResourceData, m interface{}) error {
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
