package idaas

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/okta/terraform-provider-okta/sdk/query"
)

func resourceAppSecurePasswordStore() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppSecurePasswordStoreCreate,
		ReadContext:   resourceAppSecurePasswordStoreRead,
		UpdateContext: resourceAppSecurePasswordStoreUpdate,
		DeleteContext: resourceAppSecurePasswordStoreDelete,
		Importer: &schema.ResourceImporter{
			StateContext: appImporter,
		},
		Description: `Creates a Secure Password Store Application.
	
		This resource allows you to create and configure a Secure Password Store Application.
		-> During an apply if there is change in 'status' the app will first be
		activated or deactivated in accordance with the 'status' change. Then, all
		other arguments that changed will be applied.`,
		// For those familiar with Terraform schemas be sure to check the base application schema and/or
		// the examples in the documentation
		Schema: BuildAppSwaSchema(map[string]*schema.Schema{
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
				Type:        schema.TypeString,
				Required:    true,
				Description: "Login URL",
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
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "EDIT_USERNAME_AND_PASSWORD",
				Description: "Application credentials scheme. One of: `EDIT_USERNAME_AND_PASSWORD`, `ADMIN_SETS_CREDENTIALS`, `EDIT_PASSWORD_ONLY`, `EXTERNAL_PASSWORD_SYNC`, or `SHARED_USERNAME_AND_PASSWORD`",
			},
			"reveal_password": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Allow user to reveal password. It can not be set to `true` if `credentials_scheme` is `ADMIN_SETS_CREDENTIALS`, `SHARED_USERNAME_AND_PASSWORD` or `EXTERNAL_PASSWORD_SYNC`.",
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
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Hour),
			Read:   schema.DefaultTimeout(1 * time.Hour),
			Update: schema.DefaultTimeout(1 * time.Hour),
		},
	}
}

func resourceAppSecurePasswordStoreCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	app := buildAppSecurePasswordStore(d)
	activate := d.Get("status").(string) == StatusActive
	params := &query.Params{Activate: &activate}
	_, _, err := getOktaClientFromMetadata(meta).Application.CreateApplication(ctx, app, params)
	if err != nil {
		return diag.Errorf("failed to create secure password store application: %v", err)
	}
	d.SetId(app.Id)
	err = handleAppLogo(ctx, d, meta, app.Id, app.Links)
	if err != nil {
		return diag.Errorf("failed to upload logo for secure password store application: %v", err)
	}
	return resourceAppSecurePasswordStoreRead(ctx, d, meta)
}

func resourceAppSecurePasswordStoreRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	app := sdk.NewSecurePasswordStoreApplication()
	err := fetchApp(ctx, d, meta, app)
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
	_ = d.Set("user_name_template_push_status", app.Credentials.UserNameTemplate.PushStatus)
	appRead(d, app.Name, app.Status, app.SignOnMode, app.Label, app.Accessibility, app.Visibility, app.Settings.Notes)
	return nil
}

func resourceAppSecurePasswordStoreUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	additionalChanges, err := AppUpdateStatus(ctx, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	if !additionalChanges {
		return nil
	}

	client := getOktaClientFromMetadata(meta)
	app := buildAppSecurePasswordStore(d)
	_, _, err = client.Application.UpdateApplication(ctx, d.Id(), app)
	if err != nil {
		return diag.Errorf("failed to update secure password store application: %v", err)
	}
	if d.HasChange("logo") {
		err = handleAppLogo(ctx, d, meta, app.Id, app.Links)
		if err != nil {
			o, _ := d.GetChange("logo")
			_ = d.Set("logo", o)
			return diag.Errorf("failed to upload logo for secure password store application: %v", err)
		}
	}
	return resourceAppSecurePasswordStoreRead(ctx, d, meta)
}

func resourceAppSecurePasswordStoreDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := deleteApplication(ctx, d, meta)
	if err != nil {
		return diag.Errorf("failed to delete secure password store application: %v", err)
	}
	return nil
}

func buildAppSecurePasswordStore(d *schema.ResourceData) *sdk.SecurePasswordStoreApplication {
	// Abstracts away name and SignOnMode which are constant for this app type.
	app := sdk.NewSecurePasswordStoreApplication()
	app.Label = d.Get("label").(string)

	app.Settings = &sdk.SecurePasswordStoreApplicationSettings{
		App: &sdk.SecurePasswordStoreApplicationSettingsApplication{
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
		Notes: BuildAppNotes(d),
	}
	app.Credentials = BuildSchemeAppCreds(d)
	app.Visibility = BuildAppVisibility(d)
	app.Accessibility = BuildAppAccessibility(d)

	return app
}
