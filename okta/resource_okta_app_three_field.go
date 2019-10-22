package okta

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/query"
)

func resourceAppThreeField() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppThreeFieldCreate,
		Read:   resourceAppThreeFieldRead,
		Update: resourceAppThreeFieldUpdate,
		Delete: resourceAppThreeFieldDelete,
		Exists: resourceAppExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		// For those familiar with Terraform schemas be sure to check the base application schema and/or
		// the examples in the documentation
		Schema: buildAppSwaSchema(map[string]*schema.Schema{
			"button_selector": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Login button field CSS selector",
			},
			"password_selector": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Login password field CSS selector",
			},
			"username_selector": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Login username field CSS selector",
			},
			"extra_field_selector": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Extra field CSS selector",
			},
			"extra_field_value": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Value for extra form field",
			},
			"url": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
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

func resourceAppThreeFieldCreate(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	app := buildAppThreeField(d, m)
	activate := d.Get("status").(string) == "ACTIVE"
	params := &query.Params{Activate: &activate}
	_, _, err := client.Application.CreateApplication(app, params)

	if err != nil {
		return err
	}

	d.SetId(app.Id)

	return resourceAppThreeFieldRead(d, m)
}

func resourceAppThreeFieldRead(d *schema.ResourceData, m interface{}) error {
	app := okta.NewSwaThreeFieldApplication()
	err := fetchApp(d, m, app)

	if app == nil {
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	d.Set("button_selector", app.Settings.App.ButtonSelector)
	d.Set("password_selector", app.Settings.App.PasswordSelector)
	d.Set("username_selector", app.Settings.App.UserNameSelector)
	d.Set("extra_field_selector", app.Settings.App.ExtraFieldSelector)
	d.Set("extra_field_value", app.Settings.App.ExtraFieldValue)
	d.Set("url", app.Settings.App.TargetURL)
	d.Set("url_regex", app.Settings.App.LoginUrlRegex)
	d.Set("user_name_template", app.Credentials.UserNameTemplate.Template)
	d.Set("user_name_template_type", app.Credentials.UserNameTemplate.Type)
	appRead(d, app.Name, app.Status, app.SignOnMode, app.Label, app.Accessibility, app.Visibility)

	return nil
}

func resourceAppThreeFieldUpdate(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	app := buildAppThreeField(d, m)
	_, resp, err := client.Application.UpdateApplication(d.Id(), app)

	if err != nil {
		return responseErr(resp, err)
	}

	desiredStatus := d.Get("status").(string)
	err = setAppStatus(d, client, app.Status, desiredStatus)

	if err != nil {
		return err
	}

	return resourceAppThreeFieldRead(d, m)
}

func resourceAppThreeFieldDelete(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	resp, err := client.Application.DeactivateApplication(d.Id())
	if err != nil {
		return responseErr(resp, err)
	}

	return responseErr(client.Application.DeleteApplication(d.Id()))
}

func buildAppThreeField(d *schema.ResourceData, m interface{}) *okta.SwaThreeFieldApplication {
	app := okta.NewSwaThreeFieldApplication()
	app.Label = d.Get("label").(string)

	app.Settings = &okta.SwaThreeFieldApplicationSettings{
		App: &okta.SwaThreeFieldApplicationSettingsApplication{
			TargetURL:          d.Get("url").(string),
			ButtonSelector:     d.Get("button_selector").(string),
			UserNameSelector:   d.Get("username_selector").(string),
			PasswordSelector:   d.Get("password_selector").(string),
			ExtraFieldSelector: d.Get("extra_field_selector").(string),
			ExtraFieldValue:    d.Get("extra_field_value").(string),
			LoginUrlRegex:      d.Get("url_regex").(string),
		},
	}
	app.Visibility = buildVisibility(d)

	return app
}
