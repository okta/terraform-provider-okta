package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

func resourceAppSecurePasswordStore() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppSecurePasswordStoreCreate,
		ReadContext:   resourceAppSecurePasswordStoreRead,
		UpdateContext: resourceAppSecurePasswordStoreUpdate,
		DeleteContext: resourceAppSecurePasswordStoreDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		// For those familiar with Terraform schemas be sure to check the base application schema and/or
		// the examples in the documentation
		Schema: buildAppSwaSchema(map[string]*schema.Schema{
			"password_field": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Login password field",
			},
			"username_field": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Login username field",
			},
			"url": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Login URL",
				ValidateDiagFunc: stringIsURL(validURLSchemes...),
			},
			"optional_field1": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of optional param in the login form",
			},
			"optional_field1_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of optional value in login form",
			},
			"optional_field2": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of optional param in the login form",
			},
			"optional_field2_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of optional value in login form",
			},
			"optional_field3": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of optional param in the login form",
			},
			"optional_field3_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of optional value in login form",
			},
			"credentials_scheme": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "EDIT_USERNAME_AND_PASSWORD",
				ValidateDiagFunc: stringInSlice(
					[]string{
						"EDIT_USERNAME_AND_PASSWORD",
						"ADMIN_SETS_CREDENTIALS",
						"EDIT_PASSWORD_ONLY",
						"EXTERNAL_PASSWORD_SYNC",
						"SHARED_USERNAME_AND_PASSWORD",
					}),
				Description: "Application credentials scheme",
			},
			"reveal_password": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Allow user to reveal password",
			},
			"shared_username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Shared username, required for certain schemes.",
			},
			"shared_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Shared password, required for certain schemes.",
			},
		}),
	}
}

func resourceAppSecurePasswordStoreCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	app := buildAppSecurePasswordStore(d)
	activate := d.Get("status").(string) == statusActive
	params := &query.Params{Activate: &activate}
	_, _, err := getOktaClientFromMetadata(m).Application.CreateApplication(ctx, app, params)
	if err != nil {
		return diag.Errorf("failed to create secure password store application: %v", err)
	}
	d.SetId(app.Id)
	return resourceAppSecurePasswordStoreRead(ctx, d, m)
}

func resourceAppSecurePasswordStoreRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	app := okta.NewSecurePasswordStoreApplication()
	err := fetchApp(ctx, d, m, app)
	if err != nil {
		return diag.Errorf("failed to get secure password store application: %v", err)
	}
	if app.Id == "" {
		d.SetId("")
		return nil
	}
	_ = d.Set("password_field", app.Settings.App.PasswordField)
	_ = d.Set("username_field", app.Settings.App.UsernameField)
	_ = d.Set("url", app.Settings.App.Url)
	_ = d.Set("optional_field1", app.Settings.App.OptionalField1)
	_ = d.Set("optional_field1_value", app.Settings.App.OptionalField1Value)
	_ = d.Set("optional_field2", app.Settings.App.OptionalField2)
	_ = d.Set("optional_field2_value", app.Settings.App.OptionalField2Value)
	_ = d.Set("optional_field3", app.Settings.App.OptionalField3)
	_ = d.Set("optional_field3_value", app.Settings.App.OptionalField3Value)
	_ = d.Set("credentials_scheme", app.Credentials.Scheme)
	_ = d.Set("reveal_password", app.Credentials.RevealPassword)
	_ = d.Set("shared_username", app.Credentials.UserName)
	_ = d.Set("user_name_template", app.Credentials.UserNameTemplate.Template)
	_ = d.Set("user_name_template_type", app.Credentials.UserNameTemplate.Type)
	_ = d.Set("user_name_template_suffix", app.Credentials.UserNameTemplate.Suffix)
	appRead(d, app.Name, app.Status, app.SignOnMode, app.Label, app.Accessibility, app.Visibility)
	return nil
}

func resourceAppSecurePasswordStoreUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(m)
	app := buildAppSecurePasswordStore(d)
	_, _, err := client.Application.UpdateApplication(ctx, d.Id(), app)
	if err != nil {
		return diag.Errorf("failed to update secure password store application: %v", err)
	}
	err = setAppStatus(ctx, d, client, app.Status)
	if err != nil {
		return diag.Errorf("failed to set secure password store application status: %v", err)
	}
	return resourceAppSecurePasswordStoreRead(ctx, d, m)
}

func resourceAppSecurePasswordStoreDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := deleteApplication(ctx, d, m)
	if err != nil {
		return diag.Errorf("failed to delete secure password store application: %v", err)
	}
	return nil
}

func buildAppSecurePasswordStore(d *schema.ResourceData) *okta.SecurePasswordStoreApplication {
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
