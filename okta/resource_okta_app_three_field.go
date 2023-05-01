package okta

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/okta/terraform-provider-okta/sdk/query"
)

func resourceAppThreeField() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppThreeFieldCreate,
		ReadContext:   resourceAppThreeFieldRead,
		UpdateContext: resourceAppThreeFieldUpdate,
		DeleteContext: resourceAppThreeFieldDelete,
		Importer: &schema.ResourceImporter{
			StateContext: appImporter,
		},

		// For those familiar with Terraform schemas be sure to check the base application schema and/or
		// the examples in the documentation
		Schema: buildAppSwaSchema(map[string]*schema.Schema{
			"button_selector": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Login button field CSS selector",
			},
			"password_selector": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Login password field CSS selector",
			},
			"username_selector": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Login username field CSS selector",
			},
			"extra_field_selector": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Extra field CSS selector",
			},
			"extra_field_value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Value for extra form field",
			},
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Login URL",
			},
			"url_regex": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A regex that further restricts URL to the specified regex",
			},
			"credentials_scheme": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "EDIT_USERNAME_AND_PASSWORD",
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
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Hour),
			Read:   schema.DefaultTimeout(1 * time.Hour),
			Update: schema.DefaultTimeout(1 * time.Hour),
		},
	}
}

func resourceAppThreeFieldCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(m)
	app := buildAppThreeField(d)
	activate := d.Get("status").(string) == statusActive
	params := &query.Params{Activate: &activate}
	_, _, err := client.Application.CreateApplication(ctx, app, params)
	if err != nil {
		return diag.Errorf("failed to create three field application: %v", err)
	}
	d.SetId(app.Id)
	err = handleAppLogo(ctx, d, m, app.Id, app.Links)
	if err != nil {
		return diag.Errorf("failed to upload logo for three field application: %v", err)
	}
	return resourceAppThreeFieldRead(ctx, d, m)
}

func resourceAppThreeFieldRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	app := sdk.NewSwaThreeFieldApplication()
	err := fetchApp(ctx, d, m, app)
	if err != nil {
		return diag.Errorf("failed to get three field application: %v", err)
	}
	if app.Id == "" {
		d.SetId("")
		return nil
	}
	_ = d.Set("button_selector", app.Settings.App.ButtonSelector)
	_ = d.Set("password_selector", app.Settings.App.PasswordSelector)
	_ = d.Set("username_selector", app.Settings.App.UserNameSelector)
	_ = d.Set("extra_field_selector", app.Settings.App.ExtraFieldSelector)
	_ = d.Set("extra_field_value", app.Settings.App.ExtraFieldValue)
	_ = d.Set("url", app.Settings.App.TargetURL)
	_ = d.Set("url_regex", app.Settings.App.LoginUrlRegex)
	_ = d.Set("credentials_scheme", app.Credentials.Scheme)
	_ = d.Set("reveal_password", app.Credentials.RevealPassword)
	_ = d.Set("shared_username", app.Credentials.UserName)
	_ = d.Set("user_name_template", app.Credentials.UserNameTemplate.Template)
	_ = d.Set("user_name_template_type", app.Credentials.UserNameTemplate.Type)
	_ = d.Set("user_name_template_suffix", app.Credentials.UserNameTemplate.Suffix)
	_ = d.Set("user_name_template_push_status", app.Credentials.UserNameTemplate.PushStatus)
	_ = d.Set("logo_url", linksValue(app.Links, "logo", "href"))
	appRead(d, app.Name, app.Status, app.SignOnMode, app.Label, app.Accessibility, app.Visibility, app.Settings.Notes)
	return nil
}

func resourceAppThreeFieldUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(m)
	app := buildAppThreeField(d)
	_, _, err := client.Application.UpdateApplication(ctx, d.Id(), app)
	if err != nil {
		return diag.Errorf("failed to update three field application: %v", err)
	}
	err = setAppStatus(ctx, d, client, app.Status)
	if err != nil {
		return diag.Errorf("failed to set three field application status: %v", err)
	}
	if d.HasChange("logo") {
		err = handleAppLogo(ctx, d, m, app.Id, app.Links)
		if err != nil {
			o, _ := d.GetChange("logo")
			_ = d.Set("logo", o)
			return diag.Errorf("failed to upload logo for three field application: %v", err)
		}
	}
	return resourceAppThreeFieldRead(ctx, d, m)
}

func resourceAppThreeFieldDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := deleteApplication(ctx, d, m)
	if err != nil {
		return diag.Errorf("failed to delete three field application: %v", err)
	}
	return nil
}

func buildAppThreeField(d *schema.ResourceData) *sdk.SwaThreeFieldApplication {
	app := sdk.NewSwaThreeFieldApplication()
	app.Label = d.Get("label").(string)

	app.Settings = &sdk.SwaThreeFieldApplicationSettings{
		App: &sdk.SwaThreeFieldApplicationSettingsApplication{
			TargetURL:          d.Get("url").(string),
			ButtonSelector:     d.Get("button_selector").(string),
			UserNameSelector:   d.Get("username_selector").(string),
			PasswordSelector:   d.Get("password_selector").(string),
			ExtraFieldSelector: d.Get("extra_field_selector").(string),
			ExtraFieldValue:    d.Get("extra_field_value").(string),
			LoginUrlRegex:      d.Get("url_regex").(string),
		},
		Notes: buildAppNotes(d),
	}
	app.Visibility = buildAppVisibility(d)
	app.Credentials = buildSchemeAppCreds(d)
	app.Accessibility = buildAppAccessibility(d)

	return app
}
