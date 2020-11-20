package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

func resourceAppSwa() *schema.Resource {
	return &schema.Resource{
		CustomizeDiff: func(d *schema.ResourceDiff, v interface{}) error {
			return nil
		},
		Create: resourceAppSwaCreate,
		Read:   resourceAppSwaRead,
		Update: resourceAppSwaUpdate,
		Delete: resourceAppSwaDelete,
		Exists: resourceAppExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: buildAppSwaSchema(map[string]*schema.Schema{
			"preconfigured_app": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Preconfigured app name",
			},
			"button_field": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Login button field",
			},
			"password_field": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Login password field",
			},
			"username_field": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Login username field",
			},
			"url": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Login URL",
				ValidateFunc: validateIsURL,
			},
			"url_regex": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A regex that further restricts URL to the specified regex",
			},
		}),
	}
}

func resourceAppSwaCreate(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	app := buildAppSwa(d)
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

	return resourceAppSwaRead(d, m)
}

func resourceAppSwaRead(d *schema.ResourceData, m interface{}) error {
	app := okta.NewSwaApplication()
	err := fetchApp(d, m, app)

	if app == nil {
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	_ = d.Set("button_field", app.Settings.App.ButtonField)
	_ = d.Set("password_field", app.Settings.App.PasswordField)
	_ = d.Set("username_field", app.Settings.App.UsernameField)
	_ = d.Set("url", app.Settings.App.Url)
	_ = d.Set("url_regex", app.Settings.App.LoginUrlRegex)
	_ = d.Set("user_name_template", app.Credentials.UserNameTemplate.Template)
	_ = d.Set("user_name_template_type", app.Credentials.UserNameTemplate.Type)
	_ = d.Set("user_name_template_suffix", app.Credentials.UserNameTemplate.Suffix)
	appRead(d, app.Name, app.Status, app.SignOnMode, app.Label, app.Accessibility, app.Visibility)

	return syncGroupsAndUsers(app.Id, d, m)
}

func resourceAppSwaUpdate(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	app := buildAppSwa(d)
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

	return resourceAppSwaRead(d, m)
}

func resourceAppSwaDelete(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	_, err := client.Application.DeactivateApplication(context.Background(), d.Id())
	if err != nil {
		return err
	}

	_, err = client.Application.DeleteApplication(context.Background(), d.Id())

	return err
}

func buildAppSwa(d *schema.ResourceData) *okta.SwaApplication {
	// Abstracts away name and SignOnMode which are constant for this app type.
	app := okta.NewSwaApplication()
	app.Label = d.Get("label").(string)
	name := d.Get("preconfigured_app").(string)

	if name != "" {
		app.Name = name
	}

	app.Settings = &okta.SwaApplicationSettings{
		App: &okta.SwaApplicationSettingsApplication{
			ButtonField:   d.Get("button_field").(string),
			UsernameField: d.Get("username_field").(string),
			PasswordField: d.Get("password_field").(string),
			Url:           d.Get("url").(string),
			LoginUrlRegex: d.Get("url_regex").(string),
		},
	}
	app.Visibility = buildVisibility(d)

	return app
}
