package okta

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/okta/terraform-provider-okta/sdk/query"
)

func resourceAppAutoLogin() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppAutoLoginCreate,
		ReadContext:   resourceAppAutoLoginRead,
		UpdateContext: resourceAppAutoLoginUpdate,
		DeleteContext: resourceAppAutoLoginDelete,
		Importer: &schema.ResourceImporter{
			StateContext: appImporter,
		},
		Schema: buildAppSwaSchema(map[string]*schema.Schema{
			"preconfigured_app": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Preconfigured app name",
			},
			"sign_on_url": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Login URL",
				ValidateDiagFunc: stringIsURL(validURLSchemes...),
			},
			"sign_on_redirect_url": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Post login redirect URL",
				ValidateDiagFunc: stringIsURL(validURLSchemes...),
			},
			"credentials_scheme": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateDiagFunc: elemInSlice(
					[]string{
						"EDIT_USERNAME_AND_PASSWORD",
						"ADMIN_SETS_CREDENTIALS",
						"EDIT_PASSWORD_ONLY",
						"EXTERNAL_PASSWORD_SYNC",
						"SHARED_USERNAME_AND_PASSWORD",
					}),
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
			"app_settings_json": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Application settings in JSON format",
				ValidateDiagFunc: stringIsJSON,
				StateFunc:        normalizeDataJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == ""
				},
			},
		}),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Hour),
			Read:   schema.DefaultTimeout(1 * time.Hour),
			Update: schema.DefaultTimeout(1 * time.Hour),
		},
	}
}

func resourceAppAutoLoginCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	app := buildAppAutoLogin(d)
	activate := d.Get("status").(string) == statusActive
	params := &query.Params{Activate: &activate}
	_, _, err := getOktaClientFromMetadata(m).Application.CreateApplication(ctx, app, params)
	if err != nil {
		return diag.Errorf("failed to create auto login application: %v", err)
	}
	d.SetId(app.Id)
	err = handleAppGroupsAndUsers(ctx, app.Id, d, m)
	if err != nil {
		return diag.Errorf("failed to handle groups and users for auto login application: %v", err)
	}
	err = handleAppLogo(ctx, d, m, app.Id, app.Links)
	if err != nil {
		return diag.Errorf("failed to upload logo for auto login application: %v", err)
	}
	return resourceAppAutoLoginRead(ctx, d, m)
}

func resourceAppAutoLoginRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	app := sdk.NewAutoLoginApplication()
	err := fetchApp(ctx, d, m, app)
	if err != nil {
		return diag.Errorf("failed to get auto login application: %v", err)
	}
	if app.Id == "" {
		d.SetId("")
		return nil
	}
	if app.Settings != nil {
		if app.Settings.SignOn != nil {
			_ = d.Set("sign_on_url", app.Settings.SignOn.LoginUrl)
			_ = d.Set("sign_on_redirect_url", app.Settings.SignOn.RedirectUrl)
		}
		err = setAppSettings(d, app.Settings.App)
		if err != nil {
			return diag.Errorf("failed to set auto login application settings: %v", err)
		}
	}
	_ = d.Set("credentials_scheme", app.Credentials.Scheme)
	_ = d.Set("reveal_password", app.Credentials.RevealPassword)
	_ = d.Set("shared_username", app.Credentials.UserName) // We can sync shared username but not password from upstream
	_ = d.Set("user_name_template", app.Credentials.UserNameTemplate.Template)
	_ = d.Set("user_name_template_type", app.Credentials.UserNameTemplate.Type)
	_ = d.Set("user_name_template_suffix", app.Credentials.UserNameTemplate.Suffix)
	_ = d.Set("user_name_template_push_status", app.Credentials.UserNameTemplate.PushStatus)
	_ = d.Set("logo_url", linksValue(app.Links, "logo", "href"))
	appRead(d, app.Name, app.Status, app.SignOnMode, app.Label, app.Accessibility, app.Visibility, app.Settings.Notes)
	err = syncGroupsAndUsers(ctx, app.Id, d, m)
	if err != nil {
		return diag.Errorf("failed to sync groups and users for auto login application: %v", err)
	}
	return nil
}

func resourceAppAutoLoginUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(m)
	app := buildAppAutoLogin(d)
	err := updateAppByID(ctx, d.Id(), m, app)
	if err != nil {
		return diag.Errorf("failed to update auto login application: %v", err)
	}
	err = setAppStatus(ctx, d, client, app.Status)
	if err != nil {
		return diag.Errorf("failed to set auto login application status: %v", err)
	}
	err = handleAppGroupsAndUsers(ctx, app.Id, d, m)
	if err != nil {
		return diag.Errorf("failed to handle groups and users for auto login application: %v", err)
	}
	if d.HasChange("logo") {
		err = handleAppLogo(ctx, d, m, app.Id, app.Links)
		if err != nil {
			o, _ := d.GetChange("logo")
			_ = d.Set("logo", o)
			return diag.Errorf("failed to upload logo for auto login application: %v", err)
		}
	}
	return resourceAppAutoLoginRead(ctx, d, m)
}

func resourceAppAutoLoginDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := deleteApplication(ctx, d, m)
	if err != nil {
		return diag.Errorf("failed to delete auto login application: %v", err)
	}
	return nil
}

func buildAppAutoLogin(d *schema.ResourceData) *sdk.AutoLoginApplication {
	// Abstracts away name and SignOnMode which are constant for this app type.
	app := sdk.NewAutoLoginApplication()
	app.Label = d.Get("label").(string)
	name := d.Get("preconfigured_app").(string)

	if name != "" {
		app.Name = name
	}
	app.Settings = &sdk.AutoLoginApplicationSettings{
		SignOn: &sdk.AutoLoginApplicationSettingsSignOn{
			LoginUrl:    d.Get("sign_on_url").(string),
			RedirectUrl: d.Get("sign_on_redirect_url").(string),
		},
		Notes: buildAppNotes(d),
	}
	app.Settings.App = buildAppSettings(d)
	app.Visibility = buildAppVisibility(d)
	app.Credentials = buildSchemeAppCreds(d)
	app.Accessibility = buildAppAccessibility(d)

	return app
}
