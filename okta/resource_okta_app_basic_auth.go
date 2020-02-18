package okta

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/query"
)

func resourceAppBasicAuth() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppBasicAuthCreate,
		Read:   resourceAppBasicAuthRead,
		Update: resourceAppBasicAuthUpdate,
		Delete: resourceAppBasicAuthDelete,
		Exists: resourceAppExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: buildAppSchemaWithVisibility(map[string]*schema.Schema{
			"auth_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Login button field",
			},
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Login password field",
			},
		}),
	}
}

func resourceAppBasicAuthCreate(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	app := buildAppBasicAuth(d, m)
	activate := d.Get("status").(string) == "ACTIVE"
	params := &query.Params{Activate: &activate}
	_, _, err := client.Application.CreateApplication(app, params)

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

	d.Set("url", app.Settings.App.Url)
	d.Set("auth_url", app.Settings.App.AuthURL)
	appRead(d, app.Name, app.Status, app.SignOnMode, app.Label, app.Accessibility, app.Visibility)

	return syncGroupsAndUsers(app.Id, d, m)
}

func resourceAppBasicAuthUpdate(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	app := buildAppBasicAuth(d, m)
	_, _, err := client.Application.UpdateApplication(d.Id(), app)

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
	_, err := client.Application.DeactivateApplication(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Application.DeleteApplication(d.Id())

	return err
}

func buildAppBasicAuth(d *schema.ResourceData, m interface{}) *okta.BasicAuthApplication {
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
