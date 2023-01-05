package okta

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func resourceAppWSFedAppSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppWSFedSettingsCreate,
		ReadContext:   resourceAppWSFedSettingsRead,
		UpdateContext: resourceAppWSFedSettingsUpdate,
		DeleteContext: resourceFuncNoOp,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Application ID",
				ForceNew:    true,
			},
			"settings": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Application settings in JSON format",
				ValidateDiagFunc: stringIsJSON,
				StateFunc:        normalizeDataJSON,
			},
		},
	}
}

func resourceAppWSFedSettingsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id, err := updateOrCreateAppSettings(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)
	return resourceAppWSFedSettingsRead(ctx, d, m)
}

func resourceAppWSFedSettingsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	app := okta.NewWsFederationApplication()
	err := fetchApp(ctx, d, m, app)
	if err != nil {
		return diag.Errorf("failed to get WS Federated application: %v", err)
	}
	if app.Id == "" {
		d.SetId("")
		return nil
	}
	flatMap := map[string]interface{}{}
	for key, val := range *app.Settings.App {
		if str, ok := val.(string); ok {
			if str != "" {
				flatMap[key] = str
			}
		} else if val != nil {
			flatMap[key] = val
		}
	}
	payload, _ := json.Marshal(flatMap)
	_ = d.Set("settings", string(payload))
	_ = d.Set("app_id", app.Id)
	return nil
}

func resourceAppWSFedSettingsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	_, err := updateOrCreateAppSettings(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceAppWSFedSettingsRead(ctx, d, m)
}

func updateOrCreateWSFedAppSettings(ctx context.Context, d *schema.ResourceData, m interface{}) (string, error) {
	app := okta.NewWsFederationApplication()
	appID := d.Get("app_id").(string)
	err := fetchAppByID(ctx, appID, m, app)
	if err != nil {
		return "", fmt.Errorf("failed to get WS Federated application: %v", err)
	}
	if app.Id == "" {
		return "", fmt.Errorf("application with id %s does not exist", appID)
	}
	settings := make(okta.ApplicationSettingsApplication)
	_ = json.Unmarshal([]byte(d.Get("settings").(string)), &settings)
	app.Settings.App = &settings
	_, _, err = getOktaClientFromMetadata(m).Application.UpdateApplication(ctx, appID, app)
	if err != nil {
		return "", fmt.Errorf("failed to update WS Federated application's settings: %v", err)
	}
	return app.Id, nil
}
