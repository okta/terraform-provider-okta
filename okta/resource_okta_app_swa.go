package okta

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/query"
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

func resourceAppSwaCreate(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	app := buildAppSwa(d, m)
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

	d.Set("button_field", app.Settings.App.ButtonField)
	d.Set("password_field", app.Settings.App.PasswordField)
	d.Set("username_field", app.Settings.App.UsernameField)
	d.Set("url", app.Settings.App.Url)
	d.Set("url_regex", app.Settings.App.LoginUrlRegex)
	d.Set("user_name_template", app.Credentials.UserNameTemplate.Template)
	d.Set("user_name_template_type", app.Credentials.UserNameTemplate.Type)
	appRead(d, app.Name, app.Status, app.SignOnMode, app.Label, app.Accessibility, app.Visibility)

	return syncGroupsAndUsers(app.Id, d, m)
}

func resourceAppSwaUpdate(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	app := buildAppSwa(d, m)
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

	return resourceAppSwaRead(d, m)
}

func resourceAppSwaDelete(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	_, err := client.Application.DeactivateApplication(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Application.DeleteApplication(d.Id())

	return err
}

func buildAppSwa(d *schema.ResourceData, m interface{}) *okta.SwaApplication {
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
