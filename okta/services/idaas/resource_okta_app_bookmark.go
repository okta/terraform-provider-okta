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

func resourceAppBookmark() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppBookmarkCreate,
		ReadContext:   resourceAppBookmarkRead,
		UpdateContext: resourceAppBookmarkUpdate,
		DeleteContext: resourceAppBookmarkDelete,
		Importer: &schema.ResourceImporter{
			StateContext: appImporter,
		},
		// For those familiar with Terraform schemas be sure to check the base application schema and/or
		// the examples in the documentation
		Description: `This resource allows you to create and configure a Bookmark Application.
-> During an apply if there is change in status the app will first be
activated or deactivated in accordance with the status change. Then, all
other arguments that changed will be applied.`,
		Schema: BuildAppSchemaWithVisibility(map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The URL of the bookmark.",
			},
			"request_integration": {
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
				Description: "Would you like Okta to add an integration for this app?",
			},
			"authentication_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The ID of the associated app_signon_policy. If this property is removed from the application the default sign-on-policy will be associated with this application.`,
			},
			"skip_authentication_policy": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Skip authentication policy operations. When set to true, the provider will not attempt to create, update, or delete authentication policies for this application.",
			},
		}),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Hour),
			Read:   schema.DefaultTimeout(1 * time.Hour),
			Update: schema.DefaultTimeout(1 * time.Hour),
		},
	}
}

func resourceAppBookmarkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(meta)
	app := buildAppBookmark(d)
	activate := d.Get("status").(string) == StatusActive
	params := &query.Params{Activate: &activate}
	_, _, err := client.Application.CreateApplication(ctx, app, params)
	if err != nil {
		return diag.Errorf("failed to create bookmark application: %v", err)
	}
	d.SetId(app.Id)
	err = handleAppLogo(ctx, d, meta, app.Id, app.Links)
	if err != nil {
		return diag.Errorf("failed to upload logo for bookmark application: %v", err)
	}
	if !d.Get("skip_authentication_policy").(bool) {
		err = createOrUpdateAuthenticationPolicy(ctx, d, meta, app.Id)
		if err != nil {
			return diag.Errorf("failed to set authentication policy for bookmark application: %v", err)
		}
	}
	return resourceAppBookmarkRead(ctx, d, meta)
}

func resourceAppBookmarkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	app := sdk.NewBookmarkApplication()
	err := fetchApp(ctx, d, meta, app)
	if err != nil {
		return diag.Errorf("failed to get bookmark application: %v", err)
	}
	if app.Id == "" {
		d.SetId("")
		return nil
	}
	if !d.Get("skip_authentication_policy").(bool) {
		setAuthenticationPolicy(ctx, meta, d, app.Links)
	}
	_ = d.Set("url", app.Settings.App.Url)
	_ = d.Set("request_integration", app.Settings.App.RequestIntegration)
	appRead(d, app.Name, app.Status, app.SignOnMode, app.Label, app.Accessibility, app.Visibility, app.Settings.Notes)
	_ = d.Set("logo_url", utils.LinksValue(app.Links, "logo", "href"))
	return nil
}

func resourceAppBookmarkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	additionalChanges, err := AppUpdateStatus(ctx, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	if !additionalChanges {
		return nil
	}

	client := getOktaClientFromMetadata(meta)
	app := buildAppBookmark(d)
	_, _, err = client.Application.UpdateApplication(ctx, d.Id(), app)
	if err != nil {
		return diag.Errorf("failed to update bookmark application: %v", err)
	}
	if d.HasChange("logo") {
		err = handleAppLogo(ctx, d, meta, app.Id, app.Links)
		if err != nil {
			o, _ := d.GetChange("logo")
			_ = d.Set("logo", o)
			return diag.Errorf("failed to upload logo for bookmark application: %v", err)
		}
	}
	if !d.Get("skip_authentication_policy").(bool) {
		err = createOrUpdateAuthenticationPolicy(ctx, d, meta, app.Id)
		if err != nil {
			return diag.Errorf("failed to set authentication policy for bookmark application: %v", err)
		}
	}
	return resourceAppBookmarkRead(ctx, d, meta)
}

func resourceAppBookmarkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := deleteApplication(ctx, d, meta)
	if err != nil {
		return diag.Errorf("failed to delete bookmark application: %v", err)
	}
	return nil
}

func buildAppBookmark(d *schema.ResourceData) *sdk.BookmarkApplication {
	app := sdk.NewBookmarkApplication()
	integration := d.Get("request_integration").(bool)
	app.Label = d.Get("label").(string)
	app.Settings = &sdk.BookmarkApplicationSettings{
		App: &sdk.BookmarkApplicationSettingsApplication{
			RequestIntegration: &integration,
			Url:                d.Get("url").(string),
		},
		Notes: BuildAppNotes(d),
	}
	app.Visibility = BuildAppVisibility(d)
	app.Accessibility = BuildAppAccessibility(d)
	return app
}
