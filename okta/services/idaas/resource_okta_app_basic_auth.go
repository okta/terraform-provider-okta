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

func ResourceAppBasicAuth() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppBasicAuthCreate,
		ReadContext:   resourceAppBasicAuthRead,
		UpdateContext: resourceAppBasicAuthUpdate,
		DeleteContext: resourceAppBasicAuthDelete,
		Importer: &schema.ResourceImporter{
			StateContext: appImporter,
		},
		Description: `This resource allows you to create and configure an Auto Login Okta Application.
-> During an apply if there is change in status the app will first be
activated or deactivated in accordance with the status change. Then, all
other arguments that changed will be applied.`,
		Schema: BuildAppSchemaWithVisibility(map[string]*schema.Schema{
			"auth_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The URL of the authenticating site for this app.",
			},
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The URL of the sign-in page for this app.",
			},
		}),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Hour),
			Read:   schema.DefaultTimeout(1 * time.Hour),
			Update: schema.DefaultTimeout(1 * time.Hour),
		},
	}
}

func resourceAppBasicAuthCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := GetOktaClientFromMetadata(meta)
	app := buildAppBasicAuth(d)
	activate := d.Get("status").(string) == StatusActive
	params := &query.Params{Activate: &activate}
	_, _, err := client.Application.CreateApplication(ctx, app, params)
	if err != nil {
		return diag.Errorf("failed to create basic auth application: %v", err)
	}
	d.SetId(app.Id)
	err = handleAppLogo(ctx, d, meta, app.Id, app.Links)
	if err != nil {
		return diag.Errorf("failed to upload logo for basic auth application: %v", err)
	}
	return resourceAppBasicAuthRead(ctx, d, meta)
}

func resourceAppBasicAuthRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	app := sdk.NewBasicAuthApplication()
	err := fetchApp(ctx, d, meta, app)
	if err != nil {
		return diag.Errorf("failed to get basic auth application: %v", err)
	}
	if app.Id == "" {
		d.SetId("")
		return nil
	}
	_ = d.Set("url", app.Settings.App.Url)
	_ = d.Set("auth_url", app.Settings.App.AuthURL)
	appRead(d, app.Name, app.Status, app.SignOnMode, app.Label, app.Accessibility, app.Visibility, app.Settings.Notes)
	_ = d.Set("logo_url", utils.LinksValue(app.Links, "logo", "href"))
	return nil
}

func resourceAppBasicAuthUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	additionalChanges, err := AppUpdateStatus(ctx, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	if !additionalChanges {
		return nil
	}

	client := GetOktaClientFromMetadata(meta)
	app := buildAppBasicAuth(d)
	_, _, err = client.Application.UpdateApplication(ctx, d.Id(), app)
	if err != nil {
		return diag.Errorf("failed to update basic auth application: %v", err)
	}
	if d.HasChange("logo") {
		err = handleAppLogo(ctx, d, meta, app.Id, app.Links)
		if err != nil {
			o, _ := d.GetChange("logo")
			_ = d.Set("logo", o)
			return diag.Errorf("failed to upload logo for basic auth application: %v", err)
		}
	}
	return resourceAppBasicAuthRead(ctx, d, meta)
}

func resourceAppBasicAuthDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := deleteApplication(ctx, d, meta)
	if err != nil {
		return diag.Errorf("failed to delete basic auth application: %v", err)
	}
	return nil
}

func buildAppBasicAuth(d *schema.ResourceData) *sdk.BasicAuthApplication {
	// Abstracts away name and SignOnMode which are constant for this app type.
	app := sdk.NewBasicAuthApplication()
	app.Label = d.Get("label").(string)

	app.Settings = &sdk.BasicApplicationSettings{
		App: &sdk.BasicApplicationSettingsApplication{
			AuthURL: d.Get("auth_url").(string),
			Url:     d.Get("url").(string),
		},
		Notes: BuildAppNotes(d),
	}
	app.Visibility = BuildAppVisibility(d)
	app.Accessibility = BuildAppAccessibility(d)

	return app
}
