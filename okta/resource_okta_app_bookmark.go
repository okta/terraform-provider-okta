package okta

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/query"
)

func resourceAppBookmark() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppBookmarkCreate,
		Read:   resourceAppBookmarkRead,
		Update: resourceAppBookmarkUpdate,
		Delete: resourceAppBookmarkDelete,
		Exists: resourceAppExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		// For those familiar with Terraform schemas be sure to check the base application schema and/or
		// the examples in the documentation
		Schema: buildAppSchemaWithVisibility(map[string]*schema.Schema{
			"label": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"request_integration": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
		}),
	}
}

func resourceAppBookmarkCreate(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	app := buildAppBookmark(d, m)
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

	return resourceAppBookmarkRead(d, m)
}

func resourceAppBookmarkRead(d *schema.ResourceData, m interface{}) error {
	app := okta.NewBookmarkApplication()
	err := fetchApp(d, m, app)

	if app == nil {
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	d.Set("url", app.Settings.App.Url)
	d.Set("label", app.Label)
	d.Set("request_integration", app.Settings.App.RequestIntegration)

	if err = syncGroupsAndUsers(app.Id, d, m); err != nil {
		return err
	}

	appRead(d, app.Name, app.Status, app.SignOnMode, app.Label, app.Accessibility, app.Visibility)

	return nil
}

func resourceAppBookmarkUpdate(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	app := buildAppBookmark(d, m)
	_, _, err := client.Application.UpdateApplication(d.Id(), app)

	if err != nil {
		return err
	}

	desiredStatus := d.Get("status").(string)
	err = setAppStatus(d, client, app.Status, desiredStatus)
	if err != nil {
		return err
	}

	if err := handleAppGroupsAndUsers(app.Id, d, m); err != nil {
		return err
	}

	return resourceAppBookmarkRead(d, m)
}

func resourceAppBookmarkDelete(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	_, err := client.Application.DeactivateApplication(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Application.DeleteApplication(d.Id())

	return err
}

func buildAppBookmark(d *schema.ResourceData, m interface{}) *okta.BookmarkApplication {
	app := okta.NewBookmarkApplication()
	integration := d.Get("request_integration").(bool)
	app.Label = d.Get("label").(string)
	app.Settings = &okta.BookmarkApplicationSettings{
		App: &okta.BookmarkApplicationSettingsApplication{
			RequestIntegration: &integration,
			Url:                d.Get("url").(string),
		},
	}
	app.Visibility = buildVisibility(d)

	return app
}
