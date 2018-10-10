package okta

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/query"
)

func resourceSwaApp() *schema.Resource {
	return &schema.Resource{
		CustomizeDiff: func(d *schema.ResourceDiff, v interface{}) error {
			return nil
		},
		Create: resourceSwaAppCreate,
		Read:   resourceSwaAppRead,
		Update: resourceSwaAppUpdate,
		Delete: resourceSwaAppDelete,
		Exists: resourceAppExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: buildAppSchema(map[string]*schema.Schema{
			"preconfigured_app": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Preconfigured app name",
			},
			"button_field": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Login button field",
			},
			"password_field": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Login password field",
			},
			"username_field": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Login username field",
			},
			"url": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Login URL",
				ValidateFunc: validateIsURL,
			},
			"url_regex": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A regex that further restricts URL to the specified regex",
			},
		}),
	}
}

func resourceSwaAppCreate(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	app := buildSwaApp(d, m)
	activate := d.Get("status").(string) == "ACTIVE"
	params := &query.Params{Activate: &activate}
	_, _, err := client.Application.CreateApplication(app, params)

	if err != nil {
		return err
	}

	d.SetId(app.Id)

	return resourceSwaAppRead(d, m)
}

func resourceSwaAppRead(d *schema.ResourceData, m interface{}) error {
	app := okta.NewSwaApplication()
	err := fetchApp(d, m, app)

	if err != nil {
		return err
	}

	d.Set("button_field", app.Settings.App.ButtonField)
	d.Set("password_field", app.Settings.App.PasswordField)
	d.Set("username_field", app.Settings.App.UsernameField)
	d.Set("url", app.Settings.App.Url)
	d.Set("url_regex", app.Settings.App.LoginUrlRegex)
	d.Set("user_name_template", app.Credentials.UserNameTemplate.Template)
	d.Set("user_name_template_type", app.Credentials.UserNameTemplate.Type)
	appRead(d, app.Name, app.Status, app.SignOnMode, app.Label, app.Accessibility, app.Visibility)

	return nil
}

func resourceSwaAppUpdate(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	app := buildSwaApp(d, m)
	_, _, err := client.Application.UpdateApplication(d.Id(), app)

	if err != nil {
		return err
	}

	desiredStatus := d.Get("status").(string)
	err = setAppStatus(d, client, app.Status, desiredStatus)

	if err != nil {
		return err
	}

	return resourceSwaAppRead(d, m)
}

func resourceSwaAppDelete(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	_, err := client.Application.DeactivateApplication(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Application.DeleteApplication(d.Id())

	return err
}

func buildSwaApp(d *schema.ResourceData, m interface{}) *okta.SwaApplication {
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
