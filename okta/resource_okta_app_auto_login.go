package okta

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/query"
)

func resourceAppAutoLogin() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppAutoLoginCreate,
		Read:   resourceAppAutoLoginRead,
		Update: resourceAppAutoLoginUpdate,
		Delete: resourceAppAutoLoginDelete,
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
			"sign_on_url": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Login URL",
				ValidateFunc: validateIsURL,
			},
			"sign_on_redirect_url": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Post login redirect URL",
				ValidateFunc: validateIsURL,
			},
			"credentials_scheme": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						"EDIT_USERNAME_AND_PASSWORD",
						"ADMIN_SETS_CREDENTIALS",
						"EDIT_PASSWORD_ONLY",
						"EXTERNAL_PASSWORD_SYNC",
						"SHARED_USERNAME_AND_PASSWORD",
					},
					false,
				),
				Default:     "EDIT_USERNAME_AND_PASSWORD",
				Description: "Application credentials scheme",
			},
			"reveal_password": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Allow user to reveal password",
			},
			"shared_username": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Shared username, required for certain schemes.",
			},
			"shared_password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Shared password, required for certain schemes.",
			},
		}),
	}
}

func resourceAppAutoLoginCreate(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	app := buildAppAutoLogin(d, m)
	activate := d.Get("status").(string) == "ACTIVE"
	params := &query.Params{Activate: &activate}
	_, _, err := client.Application.CreateApplication(app, params)

	if err != nil {
		return err
	}

	d.SetId(app.Id)

	return resourceAppAutoLoginRead(d, m)
}

func resourceAppAutoLoginRead(d *schema.ResourceData, m interface{}) error {
	app := okta.NewAutoLoginApplication()
	err := fetchApp(d, m, app)

	if app == nil {
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	if app.Settings.SignOn != nil {
		d.Set("sign_on_url", app.Settings.SignOn.LoginUrl)
		d.Set("sign_on_redirect_url", app.Settings.SignOn.RedirectUrl)
	}

	d.Set("credentials_scheme", app.Credentials.Scheme)
	d.Set("reveal_password", app.Credentials.RevealPassword)

	// We can sync shared username but not password from upstream
	d.Set("shared_username", app.Credentials.UserName)

	d.Set("user_name_template", app.Credentials.UserNameTemplate.Template)
	d.Set("user_name_template_type", app.Credentials.UserNameTemplate.Type)
	appRead(d, app.Name, app.Status, app.SignOnMode, app.Label, app.Accessibility, app.Visibility)

	return nil
}

func resourceAppAutoLoginUpdate(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	app := buildAppAutoLogin(d, m)
	_, _, err := client.Application.UpdateApplication(d.Id(), app)

	if err != nil {
		return err
	}

	desiredStatus := d.Get("status").(string)
	err = setAppStatus(d, client, app.Status, desiredStatus)

	if err != nil {
		return err
	}

	return resourceAppAutoLoginRead(d, m)
}

func resourceAppAutoLoginDelete(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	_, err := client.Application.DeactivateApplication(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Application.DeleteApplication(d.Id())

	return err
}

func buildAppAutoLogin(d *schema.ResourceData, m interface{}) *okta.AutoLoginApplication {
	// Abstracts away name and SignOnMode which are constant for this app type.
	app := okta.NewAutoLoginApplication()
	app.Label = d.Get("label").(string)
	name := d.Get("preconfigured_app").(string)

	if name != "" {
		app.Name = name
	}

	app.Settings = &okta.AutoLoginApplicationSettings{
		SignOn: &okta.AutoLoginApplicationSettingsSignOn{
			LoginUrl:    d.Get("sign_on_url").(string),
			RedirectUrl: d.Get("sign_on_redirect_url").(string),
		},
	}
	app.Visibility = buildVisibility(d)
	app.Credentials = buildSchemeCreds(d)

	return app
}
