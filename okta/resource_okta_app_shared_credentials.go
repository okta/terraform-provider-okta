package okta

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

func resourceAppSharedCredentials() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppSharedCredentialsCreate,
		ReadContext:   resourceAppSharedCredentialsRead,
		UpdateContext: resourceAppSharedCredentialsUpdate,
		DeleteContext: resourceAppSharedCredentialsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: appImporter,
		},
		Schema: buildAppSwaSchema(map[string]*schema.Schema{
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
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Login URL",
				ValidateDiagFunc: stringIsURL(validURLSchemes...),
			},
			"url_regex": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A regex that further restricts URL to the specified regex",
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

func resourceAppSharedCredentialsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	app := buildAppSharedCredentials(d)
	activate := d.Get("status").(string) == statusActive
	params := &query.Params{Activate: &activate}
	_, _, err := getOktaClientFromMetadata(m).Application.CreateApplication(ctx, app, params)
	if err != nil {
		return diag.Errorf("failed to create SWA shared credentials application: %v", err)
	}
	d.SetId(app.Id)
	err = handleAppGroupsAndUsers(ctx, app.Id, d, m)
	if err != nil {
		return diag.Errorf("failed to handle groups and users for SWA shared credentials application: %v", err)
	}
	err = handleAppLogo(ctx, d, m, app.Id, app.Links)
	if err != nil {
		return diag.Errorf("failed to upload logo for SWA shared credentials application: %v", err)
	}
	return resourceAppSharedCredentialsRead(ctx, d, m)
}

func resourceAppSharedCredentialsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	app := okta.NewBrowserPluginApplication()
	err := fetchApp(ctx, d, m, app)
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
	_ = d.Set("logo_url", linksValue(app.Links, "logo", "href"))
	_ = d.Set("accessibility_login_redirect_url", app.Accessibility.LoginRedirectUrl)
	appRead(d, app.Name, app.Status, app.SignOnMode, app.Label, app.Accessibility, app.Visibility, app.Settings.Notes)
	err = syncGroupsAndUsers(ctx, app.Id, d, m)
	if err != nil {
		return diag.Errorf("failed to sync groups and users for SWA shared credentials application: %v", err)
	}
	return nil
}

func resourceAppSharedCredentialsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(m)
	app := buildAppSharedCredentials(d)
	_, _, err := client.Application.UpdateApplication(ctx, d.Id(), app)
	if err != nil {
		return diag.Errorf("failed to update SWA shared credentials application: %v", err)
	}
	err = setAppStatus(ctx, d, client, app.Status)
	if err != nil {
		return diag.Errorf("failed to set SWA shared credentials application status: %v", err)
	}
	err = handleAppGroupsAndUsers(ctx, app.Id, d, m)
	if err != nil {
		return diag.Errorf("failed to handle groups and users for SWA shared credentials application: %v", err)
	}
	if d.HasChange("logo") {
		err = handleAppLogo(ctx, d, m, app.Id, app.Links)
		if err != nil {
			o, _ := d.GetChange("logo")
			_ = d.Set("logo", o)
			return diag.Errorf("failed to upload logo for SWA shared credentials application: %v", err)
		}
	}
	return resourceAppSharedCredentialsRead(ctx, d, m)
}

func resourceAppSharedCredentialsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := deleteApplication(ctx, d, m)
	if err != nil {
		return diag.Errorf("failed to delete SWA shared credentials application: %v", err)
	}
	return nil
}

func buildAppSharedCredentials(d *schema.ResourceData) *okta.BrowserPluginApplication {
	app := okta.NewBrowserPluginApplication()
	app.Name = "template_swa"
	name := d.Get("preconfigured_app").(string)
	if name != "" {
		app.Name = name
	}
	app.Label = d.Get("label").(string)
	app.Settings = &okta.ApplicationSettings{
		App: &okta.ApplicationSettingsApplication{
			"buttonField":   d.Get("button_field").(string),
			"loginUrlRegex": d.Get("url_regex").(string),
			"passwordField": d.Get("password_field").(string),
			"url":           d.Get("url").(string),
			"usernameField": d.Get("username_field").(string),
			"redirectUrl":   d.Get("redirect_url").(string),
			"checkbox":      d.Get("checkbox").(string),
		},
		Notes: buildAppNotes(d),
	}
	app.Credentials = &okta.SchemeApplicationCredentials{
		UserNameTemplate: buildUserNameTemplate(d),
		Password: &okta.PasswordCredential{
			Value: d.Get("shared_password").(string),
		},
		Scheme:   "SHARED_USERNAME_AND_PASSWORD",
		UserName: d.Get("shared_username").(string),
	}
	app.Visibility = buildAppVisibility(d)
	app.Accessibility = buildAppAccessibility(d)
	return app
}
