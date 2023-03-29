package okta

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/okta/terraform-provider-okta/sdk/query"
)

func resourceAppBasicAuth() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppBasicAuthCreate,
		ReadContext:   resourceAppBasicAuthRead,
		UpdateContext: resourceAppBasicAuthUpdate,
		DeleteContext: resourceAppBasicAuthDelete,
		Importer: &schema.ResourceImporter{
			StateContext: appImporter,
		},
		Schema: buildAppSchemaWithVisibility(map[string]*schema.Schema{
			"auth_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Login button field",
			},
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Login password field",
			},
		}),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Hour),
			Read:   schema.DefaultTimeout(1 * time.Hour),
			Update: schema.DefaultTimeout(1 * time.Hour),
		},
	}
}

func resourceAppBasicAuthCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(m)
	app := buildAppBasicAuth(d)
	activate := d.Get("status").(string) == statusActive
	params := &query.Params{Activate: &activate}
	_, _, err := client.Application.CreateApplication(ctx, app, params)
	if err != nil {
		return diag.Errorf("failed to create basic auth application: %v", err)
	}
	d.SetId(app.Id)
	err = handleAppGroupsAndUsers(ctx, app.Id, d, m)
	if err != nil {
		return diag.Errorf("failed to handle groups and users for basic auth application: %v", err)
	}
	err = handleAppLogo(ctx, d, m, app.Id, app.Links)
	if err != nil {
		return diag.Errorf("failed to upload logo for basic auth application: %v", err)
	}
	return resourceAppBasicAuthRead(ctx, d, m)
}

func resourceAppBasicAuthRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	app := sdk.NewBasicAuthApplication()
	err := fetchApp(ctx, d, m, app)
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
	_ = d.Set("logo_url", linksValue(app.Links, "logo", "href"))
	err = syncGroupsAndUsers(ctx, app.Id, d, m)
	if err != nil {
		return diag.Errorf("failed to sync groups and users for basic auth application: %v", err)
	}
	return nil
}

func resourceAppBasicAuthUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(m)
	app := buildAppBasicAuth(d)
	_, _, err := client.Application.UpdateApplication(ctx, d.Id(), app)
	if err != nil {
		return diag.Errorf("failed to update basic auth application: %v", err)
	}
	err = setAppStatus(ctx, d, client, app.Status)
	if err != nil {
		return diag.Errorf("failed to set basic auth application status: %v", err)
	}
	err = handleAppGroupsAndUsers(ctx, app.Id, d, m)
	if err != nil {
		return diag.Errorf("failed to handle groups and users for basic auth application: %v", err)
	}
	if d.HasChange("logo") {
		err = handleAppLogo(ctx, d, m, app.Id, app.Links)
		if err != nil {
			o, _ := d.GetChange("logo")
			_ = d.Set("logo", o)
			return diag.Errorf("failed to upload logo for basic auth application: %v", err)
		}
	}
	return resourceAppBasicAuthRead(ctx, d, m)
}

func resourceAppBasicAuthDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := deleteApplication(ctx, d, m)
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
		Notes: buildAppNotes(d),
	}
	app.Visibility = buildAppVisibility(d)
	app.Accessibility = buildAppAccessibility(d)

	return app
}
