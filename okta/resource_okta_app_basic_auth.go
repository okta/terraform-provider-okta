package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

func resourceAppBasicAuth() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppBasicAuthCreate,
		Read:   resourceAppBasicAuthRead,
		Update: resourceAppBasicAuthUpdate,
		Delete: resourceAppBasicAuthDelete,
		Exists: resourceAppExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: buildAppSchemaWithVisibility(map[string]*schema.Schema{
			"auth_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Login button field",
			},
			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Login password field",
			},
		}),
	}
}

func resourceAppBasicAuthCreate(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	app := buildAppBasicAuth(d)
	activate := d.Get("status").(string) == statusActive
	params := &query.Params{Activate: &activate}
	_, _, err := client.Application.CreateApplication(context.Background(), app, params)

	if err != nil {
		return err
	}

	d.SetId(app.Id)

	err = handleAppGroupsAndUsers(app.Id, d, m)

	if err != nil {
		return err
	}

	return resourceAppBasicAuthRead(d, m)
}

func resourceAppBasicAuthRead(d *schema.ResourceData, m interface{}) error {
	app := okta.NewBasicAuthApplication()
	err := fetchApp(d, m, app)

	if app == nil {
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	_ = d.Set("url", app.Settings.App.Url)
	_ = d.Set("auth_url", app.Settings.App.AuthURL)
	appRead(d, app.Name, app.Status, app.SignOnMode, app.Label, app.Accessibility, app.Visibility)

	return syncGroupsAndUsers(app.Id, d, m)
}

func resourceAppBasicAuthUpdate(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	app := buildAppBasicAuth(d)
	_, _, err := client.Application.UpdateApplication(context.Background(), d.Id(), app)

	if err != nil {
		return err
	}

	desiredStatus := d.Get("status").(string)
	err = setAppStatus(d, client, app.Status, desiredStatus)

	if err != nil {
		return err
	}

	err = handleAppGroupsAndUsers(app.Id, d, m)

	if err != nil {
		return err
	}

	return resourceAppBasicAuthRead(d, m)
}

func resourceAppBasicAuthDelete(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	_, err := client.Application.DeactivateApplication(context.Background(), d.Id())
	if err != nil {
		return err
	}

	_, err = client.Application.DeleteApplication(context.Background(), d.Id())

	return err
}

func buildAppBasicAuth(d *schema.ResourceData) *okta.BasicAuthApplication {
	// Abstracts away name and SignOnMode which are constant for this app type.
	app := okta.NewBasicAuthApplication()
	app.Label = d.Get("label").(string)

	app.Settings = &okta.BasicApplicationSettings{
		App: &okta.BasicApplicationSettingsApplication{
			AuthURL: d.Get("auth_url").(string),
			Url:     d.Get("url").(string),
		},
	}
	app.Visibility = buildVisibility(d)

	return app
}
