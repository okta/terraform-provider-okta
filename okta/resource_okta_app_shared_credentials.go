package okta

import (
	"context"

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
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: buildAppSwaSchema(map[string]*schema.Schema{
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
	_ = d.Set("username_field", flatMap["loginUrlRegex"])
	_ = d.Set("password_field", flatMap["passwordField"])
	_ = d.Set("url", flatMap["url"])
	_ = d.Set("url_regex", flatMap["usernameField"])
	_ = d.Set("redirect_url", flatMap["redirectUrl"])
	_ = d.Set("checkbox", flatMap["checkbox"])
	_ = d.Set("shared_username", app.Credentials.UserName)
	_ = d.Set("user_name_template", app.Credentials.UserNameTemplate.Template)
	_ = d.Set("user_name_template_type", app.Credentials.UserNameTemplate.Type)
	_ = d.Set("user_name_template_suffix", app.Credentials.UserNameTemplate.Suffix)
	appRead(d, app.Name, app.Status, app.SignOnMode, app.Label, app.Accessibility, app.Visibility)
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
	app.Label = d.Get("label").(string)
	app.Settings = &okta.ApplicationSettings{
		App: &okta.ApplicationSettingsApplication{
			"buttonField":   d.Get("button_field").(string),
			"loginUrlRegex": d.Get("username_field").(string),
			"passwordField": d.Get("password_field").(string),
			"url":           d.Get("url").(string),
			"usernameField": d.Get("url_regex").(string),
			"redirectUrl":   d.Get("redirect_url").(string),
			"checkbox":      d.Get("checkbox").(string),
		},
	}
	app.Credentials = &okta.SchemeApplicationCredentials{
		UserNameTemplate: &okta.ApplicationCredentialsUsernameTemplate{
			Template: d.Get("user_name_template").(string),
			Type:     d.Get("user_name_template_type").(string),
			Suffix:   d.Get("user_name_template_suffix").(string),
		},
		Password: &okta.PasswordCredential{
			Value: d.Get("shared_password").(string),
		},
		Scheme:   "SHARED_USERNAME_AND_PASSWORD",
		UserName: d.Get("shared_username").(string),
	}
	app.Accessibility = &okta.ApplicationAccessibility{
		SelfService:      boolPtr(d.Get("accessibility_self_service").(bool)),
		ErrorRedirectUrl: d.Get("accessibility_error_redirect_url").(string),
		LoginRedirectUrl: d.Get("accessibility_login_redirect_url").(string),
	}
	app.Visibility = buildVisibility(d)
	return app
}
