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

func ResourceAppSharedCredentials() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppSharedCredentialsCreate,
		ReadContext:   resourceAppSharedCredentialsRead,
		UpdateContext: resourceAppSharedCredentialsUpdate,
		DeleteContext: resourceAppSharedCredentialsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: appImporter,
		},
		Description: `Creates a SWA shared credentials app.
This resource allows you to create and configure SWA shared credentials app.
-> During an apply if there is change in 'status' the app will first be
activated or deactivated in accordance with the 'status' change. Then, all
other arguments that changed will be applied.`,
		Schema: BuildAppSwaSchema(map[string]*schema.Schema{
			"preconfigured_app": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of application from the Okta Integration Network, if not included a custom app will be created.",
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
				Description: "The URL of the sign-in page for this app.",
			},
			"url_regex": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A regular expression that further restricts url to the specified regular expression.",
			},
			"redirect_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Secondary URL of the sign-in page for this app",
			},
			"checkbox": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "CSS selector for the checkbox",
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

func resourceAppSharedCredentialsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	app := buildAppSharedCredentials(d)
	activate := d.Get("status").(string) == StatusActive
	params := &query.Params{Activate: &activate}
	_, _, err := GetOktaClientFromMetadata(meta).Application.CreateApplication(ctx, app, params)
	if err != nil {
		return diag.Errorf("failed to create SWA shared credentials application: %v", err)
	}
	d.SetId(app.Id)
	err = handleAppLogo(ctx, d, meta, app.Id, app.Links)
	if err != nil {
		return diag.Errorf("failed to upload logo for SWA shared credentials application: %v", err)
	}
	return resourceAppSharedCredentialsRead(ctx, d, meta)
}

func resourceAppSharedCredentialsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	app := sdk.NewBrowserPluginApplication()
	err := fetchApp(ctx, d, meta, app)
	if err != nil {
		return diag.Errorf("failed to get SWA shared credentials application: %v", err)
	}
	if app.Id == "" {
		d.SetId("")
		return nil
	}
	flatMap := map[string]interface{}(*app.Settings.App)
	_ = d.Set("button_field", flatMap["buttonField"])
	_ = d.Set("url_regex", flatMap["loginUrlRegex"])
	_ = d.Set("password_field", flatMap["passwordField"])
	_ = d.Set("url", flatMap["url"])
	_ = d.Set("username_field", flatMap["usernameField"])
	_ = d.Set("redirect_url", flatMap["redirectUrl"])
	_ = d.Set("checkbox", flatMap["checkbox"])
	_ = d.Set("shared_username", app.Credentials.UserName)
	_ = d.Set("user_name_template", app.Credentials.UserNameTemplate.Template)
	_ = d.Set("user_name_template_type", app.Credentials.UserNameTemplate.Type)
	_ = d.Set("user_name_template_suffix", app.Credentials.UserNameTemplate.Suffix)
	_ = d.Set("user_name_template_push_status", app.Credentials.UserNameTemplate.PushStatus)
	_ = d.Set("logo_url", utils.LinksValue(app.Links, "logo", "href"))
	_ = d.Set("accessibility_login_redirect_url", app.Accessibility.LoginRedirectUrl)
	appRead(d, app.Name, app.Status, app.SignOnMode, app.Label, app.Accessibility, app.Visibility, app.Settings.Notes)
	return nil
}

func resourceAppSharedCredentialsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	additionalChanges, err := AppUpdateStatus(ctx, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	if !additionalChanges {
		return nil
	}

	client := GetOktaClientFromMetadata(meta)
	app := buildAppSharedCredentials(d)
	_, _, err = client.Application.UpdateApplication(ctx, d.Id(), app)
	if err != nil {
		return diag.Errorf("failed to update SWA shared credentials application: %v", err)
	}
	if d.HasChange("logo") {
		err = handleAppLogo(ctx, d, meta, app.Id, app.Links)
		if err != nil {
			o, _ := d.GetChange("logo")
			_ = d.Set("logo", o)
			return diag.Errorf("failed to upload logo for SWA shared credentials application: %v", err)
		}
	}
	return resourceAppSharedCredentialsRead(ctx, d, meta)
}

func resourceAppSharedCredentialsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := deleteApplication(ctx, d, meta)
	if err != nil {
		return diag.Errorf("failed to delete SWA shared credentials application: %v", err)
	}
	return nil
}

func buildAppSharedCredentials(d *schema.ResourceData) *sdk.BrowserPluginApplication {
	app := sdk.NewBrowserPluginApplication()
	app.Name = "template_swa"
	name := d.Get("preconfigured_app").(string)
	if name != "" {
		app.Name = name
	}
	app.Label = d.Get("label").(string)
	app.Settings = &sdk.ApplicationSettings{
		App: &sdk.ApplicationSettingsApplication{
			"buttonField":   d.Get("button_field").(string),
			"loginUrlRegex": d.Get("url_regex").(string),
			"passwordField": d.Get("password_field").(string),
			"url":           d.Get("url").(string),
			"usernameField": d.Get("username_field").(string),
			"redirectUrl":   d.Get("redirect_url").(string),
			"checkbox":      d.Get("checkbox").(string),
		},
		Notes: BuildAppNotes(d),
	}
	app.Credentials = &sdk.SchemeApplicationCredentials{
		UserNameTemplate: BuildUserNameTemplate(d),
		Password: &sdk.PasswordCredential{
			Value: d.Get("shared_password").(string),
		},
		Scheme:   "SHARED_USERNAME_AND_PASSWORD",
		UserName: d.Get("shared_username").(string),
	}
	app.Visibility = BuildAppVisibility(d)
	app.Accessibility = BuildAppAccessibility(d)
	return app
}
