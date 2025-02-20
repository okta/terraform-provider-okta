package idaas

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
)

func ResourceAppSamlAppSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppSamlSettingsCreate,
		ReadContext:   resourceAppSamlSettingsRead,
		UpdateContext: resourceAppSamlSettingsUpdate,
		DeleteContext: utils.ResourceFuncNoOp,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: `Manages app settings of the SAML application.
This resource allows you to manage app settings of the SAML Application . It's basically the same as
app_settings_json field in okta_app_saml resource and can be used in cases where settings require to be managed separately.`,
		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the application.",
				ForceNew:    true,
			},
			"settings": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Application settings in JSON format",
				ValidateDiagFunc: stringIsJSON,
				StateFunc:        utils.NormalizeDataJSON,
				DiffSuppressFunc: utils.NoChangeInObjectFromUnmarshaledJSON,
			},
		},
	}
}

func resourceAppSamlSettingsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id, err := updateOrCreateAppSettings(ctx, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)
	return resourceAppSamlSettingsRead(ctx, d, meta)
}

func resourceAppSamlSettingsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	app := sdk.NewSamlApplication()
	err := fetchApp(ctx, d, meta, app)
	if err != nil {
		return diag.Errorf("failed to get SAML application: %v", err)
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

func resourceAppSamlSettingsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_, err := updateOrCreateAppSettings(ctx, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceAppSamlSettingsRead(ctx, d, meta)
}

func updateOrCreateAppSettings(ctx context.Context, d *schema.ResourceData, meta interface{}) (string, error) {
	app := sdk.NewSamlApplication()
	appID := d.Get("app_id").(string)
	err := fetchAppByID(ctx, appID, meta, app)
	if err != nil {
		return "", fmt.Errorf("failed to get SAML application: %v", err)
	}
	if app.Id == "" {
		return "", fmt.Errorf("application with id %s does not exist", appID)
	}
	settings := make(sdk.ApplicationSettingsApplication)
	_ = json.Unmarshal([]byte(d.Get("settings").(string)), &settings)
	app.Settings.App = &settings
	_, _, err = GetOktaClientFromMetadata(meta).Application.UpdateApplication(ctx, appID, app)
	if err != nil {
		return "", fmt.Errorf("failed to update SAML application's settings: %v", err)
	}
	return app.Id, nil
}
