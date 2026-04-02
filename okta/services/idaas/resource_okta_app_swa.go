package idaas

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/okta/terraform-provider-okta/sdk/query"
)

func resourceAppSwa() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppSwaCreate,
		ReadContext:   resourceAppSwaRead,
		UpdateContext: resourceAppSwaUpdate,
		DeleteContext: resourceAppSwaDelete,
		Importer: &schema.ResourceImporter{
			StateContext: appImporter,
		},
		Description: `Creates a SWA Application.
		
This resource allows you to create and configure a SWA Application.
-> During an apply if there is change in 'status' the app will first be
activated or deactivated in accordance with the 'status' change. Then, all
other arguments that changed will be applied.`,
		Schema: BuildAppSwaSchema(map[string]*schema.Schema{
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
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Login URL",
			},
			"url_regex": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A regex that further restricts URL to the specified regex",
			},
			"checkbox": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "CSS selector for the checkbox",
			},
			"redirect_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "If going to the login page URL redirects to another page, then enter that URL here",
			},
			"app_settings_json": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Application settings in JSON format. If app_settings_json is defined, the individual fields button_field, password_field, username_field, url, url_regex, checkbox, redirect_url will be ignored",
				ValidateDiagFunc: stringIsJSON,
				StateFunc:        utils.NormalizeDataJSON,
				DiffSuppressFunc: utils.NoChangeInObjectFromUnmarshaledJSON,
			},
		}),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Hour),
			Read:   schema.DefaultTimeout(1 * time.Hour),
			Update: schema.DefaultTimeout(1 * time.Hour),
		},
	}
}

func resourceAppSwaCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(meta)
	_, appSettingsJSONIsUsed := d.GetOk("app_settings_json")
	switch {
	case appSettingsJSONIsUsed: // Need to support app_settings_json while also supporting customers who have existing configuration based on the Settings.App struct
		app := buildAppSwaWithApplicationSettingsJSON(d)
		activate := d.Get("status").(string) == StatusActive
		params := &query.Params{Activate: &activate}
		_, _, err := client.Application.CreateApplication(ctx, app, params)
		if err != nil {
			return diag.Errorf("failed to create SWA application: %v", err)
		}
		d.SetId(app.Id)
		err = handleAppLogo(ctx, d, meta, app.Id, app.Links)
		if err != nil {
			return diag.Errorf("failed to upload logo for SWA application: %v", err)
		}
		return resourceAppSwaRead(ctx, d, meta)
	default:
		app := buildAppSwa(d)
		activate := d.Get("status").(string) == StatusActive
		params := &query.Params{Activate: &activate}
		_, _, err := client.Application.CreateApplication(ctx, app, params)
		if err != nil {
			return diag.Errorf("failed to create SWA application: %v", err)
		}
		d.SetId(app.Id)
		err = handleAppLogo(ctx, d, meta, app.Id, app.Links)
		if err != nil {
			return diag.Errorf("failed to upload logo for SWA application: %v", err)
		}
		return resourceAppSwaRead(ctx, d, meta)
	}

}

func resourceAppSwaRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_, appSettingsJSONIsUsed := d.GetOk("app_settings_json") // boolean value to indicate if app_settings_json is set
	switch {
	case appSettingsJSONIsUsed: // Need to support app_settings_json while also supporting customers who have existing configuration based on the Settings.App struct
		app := sdk.NewSwaApplicationWithApplicationSettingsJSON()
		err := fetchApp(ctx, d, meta, app)
		if err != nil {
			return diag.Errorf("failed to get SWA application: %v", err)
		}
		if app.Id == "" {
			d.SetId("")
			return nil
		}
		err = setAppSettings(d, app.Settings.App)
		if err != nil {
			return diag.Errorf("failed to set SAML app settings: %v", err)
		}
		_ = d.Set("user_name_template", app.Credentials.UserNameTemplate.Template)
		_ = d.Set("user_name_template_type", app.Credentials.UserNameTemplate.Type)
		_ = d.Set("user_name_template_suffix", app.Credentials.UserNameTemplate.Suffix)
		_ = d.Set("user_name_template_push_status", app.Credentials.UserNameTemplate.PushStatus)
		_ = d.Set("logo_url", utils.LinksValue(app.Links, "logo", "href"))
		appRead(d, app.Name, app.Status, app.SignOnMode, app.Label, app.Accessibility, app.Visibility, app.Settings.Notes)
	default:
		app := sdk.NewSwaApplication()
		err := fetchApp(ctx, d, meta, app)
		if err != nil {
			return diag.Errorf("failed to get SWA application: %v", err)
		}
		if app.Id == "" {
			d.SetId("")
			return nil
		}
		_ = d.Set("button_field", app.Settings.App.ButtonField)
		_ = d.Set("password_field", app.Settings.App.PasswordField)
		_ = d.Set("username_field", app.Settings.App.UsernameField)
		_ = d.Set("url", app.Settings.App.Url)
		_ = d.Set("url_regex", app.Settings.App.LoginUrlRegex)
		_ = d.Set("checkbox", app.Settings.App.Checkbox)
		_ = d.Set("redirect_url", app.Settings.App.RedirectUrl)
		_ = d.Set("user_name_template", app.Credentials.UserNameTemplate.Template)
		_ = d.Set("user_name_template_type", app.Credentials.UserNameTemplate.Type)
		_ = d.Set("user_name_template_suffix", app.Credentials.UserNameTemplate.Suffix)
		_ = d.Set("user_name_template_push_status", app.Credentials.UserNameTemplate.PushStatus)
		_ = d.Set("logo_url", utils.LinksValue(app.Links, "logo", "href"))
		appRead(d, app.Name, app.Status, app.SignOnMode, app.Label, app.Accessibility, app.Visibility, app.Settings.Notes)
	}
	return nil

}

func resourceAppSwaUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	additionalChanges, err := AppUpdateStatus(ctx, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	if !additionalChanges {
		return nil
	}
	_, appSettingsJSONIsUsed := d.GetOk("app_settings_json") // boolean value to indicate if app_settings_json is set
	switch {
	case appSettingsJSONIsUsed:
		client := getOktaClientFromMetadata(meta)
		app := buildAppSwaWithApplicationSettingsJSON(d)
		_, _, err = client.Application.UpdateApplication(ctx, d.Id(), app)
		if err != nil {
			return diag.Errorf("failed to update SWA application: %v", err)
		}
		if d.HasChange("logo") {
			err = handleAppLogo(ctx, d, meta, app.Id, app.Links)
			if err != nil {
				o, _ := d.GetChange("logo")
				_ = d.Set("logo", o)
				return diag.Errorf("failed to upload logo for SWA application: %v", err)
			}
		}
	default:
		client := getOktaClientFromMetadata(meta)
		app := buildAppSwa(d)
		_, _, err = client.Application.UpdateApplication(ctx, d.Id(), app)
		if err != nil {
			return diag.Errorf("failed to update SWA application: %v", err)
		}
		if d.HasChange("logo") {
			err = handleAppLogo(ctx, d, meta, app.Id, app.Links)
			if err != nil {
				o, _ := d.GetChange("logo")
				_ = d.Set("logo", o)
				return diag.Errorf("failed to upload logo for SWA application: %v", err)
			}
		}
	}
	return resourceAppSwaRead(ctx, d, meta)
}

func resourceAppSwaDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := deleteApplication(ctx, d, meta)
	if err != nil {
		return diag.Errorf("failed to delete SWA application: %v", err)
	}
	return nil
}

func buildAppSwa(d *schema.ResourceData) *sdk.SwaApplication {
	// Abstracts away name and SignOnMode which are constant for this app type.
	app := sdk.NewSwaApplication()
	app.Label = d.Get("label").(string)
	name := d.Get("preconfigured_app").(string)
	if name != "" {
		app.Name = name
		app.SignOnMode = "AUTO_LOGIN" // in case pre-configured app has more than one sign-on modes
	}
	app.Settings = &sdk.SwaApplicationSettings{
		App: &sdk.SwaApplicationSettingsApplication{
			ButtonField:   d.Get("button_field").(string),
			UsernameField: d.Get("username_field").(string),
			PasswordField: d.Get("password_field").(string),
			Url:           d.Get("url").(string),
			LoginUrlRegex: d.Get("url_regex").(string),
			RedirectUrl:   d.Get("redirect_url").(string),
			Checkbox:      d.Get("checkbox").(string),
		},
		Notes: BuildAppNotes(d),
	}
	app.Visibility = BuildAppVisibility(d)
	app.Accessibility = BuildAppAccessibility(d)
	app.Credentials = &sdk.SchemeApplicationCredentials{
		UserNameTemplate: BuildUserNameTemplate(d),
	}
	return app
}

func buildAppSwaWithApplicationSettingsJSON(d *schema.ResourceData) *sdk.SwaApplicationWithApplicationSettingsJSON {
	app := sdk.NewSwaApplicationWithApplicationSettingsJSON()
	app.Label = d.Get("label").(string)
	name := d.Get("preconfigured_app").(string)
	if name != "" {
		app.Name = name
		app.SignOnMode = "AUTO_LOGIN" // in case pre-configured app has more than one sign-on modes
	}
	app.Settings = &sdk.SwaApplicationSettingsWithJSON{
		App:   BuildAppSettings(d),
		Notes: BuildAppNotes(d),
	}
	app.Visibility = BuildAppVisibility(d)
	app.Accessibility = BuildAppAccessibility(d)
	app.Credentials = &sdk.SchemeApplicationCredentials{
		UserNameTemplate: BuildUserNameTemplate(d),
	}
	return app
}
