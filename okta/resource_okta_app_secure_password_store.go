package okta

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/query"
)

func resourceAppSecurePasswordStore() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppSecurePasswordStoreCreate,
		Read:   resourceAppSecurePasswordStoreRead,
		Update: resourceAppSecurePasswordStoreUpdate,
		Delete: resourceAppSecurePasswordStoreDelete,
		Exists: resourceAppExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		// For those familiar with Terraform schemas be sure to check the base application schema and/or
		// the examples in the documentation
		Schema: buildAppSwaSchema(map[string]*schema.Schema{
			"password_field": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Login password field",
			},
			"username_field": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Login username field",
			},
			"url": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Login URL",
				ValidateFunc: validateIsURL,
			},
			"optional_field1": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of optional param in the login form",
			},
			"optional_field1_value": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of optional value in login form",
			},
			"optional_field2": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of optional param in the login form",
			},
			"optional_field2_value": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of optional value in login form",
			},
			"optional_field3": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of optional param in the login form",
			},
			"optional_field3_value": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of optional value in login form",
			},
			"credentials_scheme": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "EDIT_USERNAME_AND_PASSWORD",
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

func resourceAppSecurePasswordStoreCreate(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	app := buildAppSecurePasswordStore(d, m)
	activate := d.Get("status").(string) == "ACTIVE"
	params := &query.Params{Activate: &activate}
	_, _, err := client.Application.CreateApplication(app, params)

	if err != nil {
		return err
	}

	d.SetId(app.Id)

	return resourceAppSecurePasswordStoreRead(d, m)
}

func resourceAppSecurePasswordStoreRead(d *schema.ResourceData, m interface{}) error {
	app := okta.NewSecurePasswordStoreApplication()
	err := fetchApp(d, m, app)

	if app == nil {
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	d.Set("password_field", app.Settings.App.PasswordField)
	d.Set("username_field", app.Settings.App.UsernameField)
	d.Set("url", app.Settings.App.Url)
	d.Set("optional_field1", app.Settings.App.OptionalField1)
	d.Set("optional_field1_value", app.Settings.App.OptionalField1Value)
	d.Set("optional_field2", app.Settings.App.OptionalField2)
	d.Set("optional_field2_value", app.Settings.App.OptionalField2Value)
	d.Set("optional_field3", app.Settings.App.OptionalField3)
	d.Set("optional_field3_value", app.Settings.App.OptionalField3Value)
	d.Set("credentials_scheme", app.Credentials.Scheme)
	d.Set("reveal_password", app.Credentials.RevealPassword)
	d.Set("shared_username", app.Credentials.UserName)
	d.Set("user_name_template", app.Credentials.UserNameTemplate.Template)
	d.Set("user_name_template_type", app.Credentials.UserNameTemplate.Type)
	appRead(d, app.Name, app.Status, app.SignOnMode, app.Label, app.Accessibility, app.Visibility)

	return nil
}

func resourceAppSecurePasswordStoreUpdate(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	app := buildAppSecurePasswordStore(d, m)
	_, _, err := client.Application.UpdateApplication(d.Id(), app)

	if err != nil {
		return err
	}

	desiredStatus := d.Get("status").(string)
	err = setAppStatus(d, client, app.Status, desiredStatus)

	if err != nil {
		return err
	}

	return resourceAppSecurePasswordStoreRead(d, m)
}

func resourceAppSecurePasswordStoreDelete(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	_, err := client.Application.DeactivateApplication(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Application.DeleteApplication(d.Id())

	return err
}

func buildAppSecurePasswordStore(d *schema.ResourceData, m interface{}) *okta.SecurePasswordStoreApplication {
	// Abstracts away name and SignOnMode which are constant for this app type.
	app := okta.NewSecurePasswordStoreApplication()
	app.Label = d.Get("label").(string)

	app.Settings = &okta.SecurePasswordStoreApplicationSettings{
		App: &okta.SecurePasswordStoreApplicationSettingsApplication{
			Url:                 d.Get("url").(string),
			PasswordField:       d.Get("password_field").(string),
			UsernameField:       d.Get("username_field").(string),
			OptionalField1:      d.Get("optional_field1").(string),
			OptionalField2:      d.Get("optional_field2").(string),
			OptionalField3:      d.Get("optional_field3").(string),
			OptionalField1Value: d.Get("optional_field1_value").(string),
			OptionalField2Value: d.Get("optional_field2_value").(string),
			OptionalField3Value: d.Get("optional_field3_value").(string),
		},
	}
	app.Credentials = buildSchemeCreds(d)
	app.Visibility = buildVisibility(d)

	return app
}
