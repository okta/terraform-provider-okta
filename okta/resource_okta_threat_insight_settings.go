package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func resourceThreatInsightSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceThreatInsightSettingsCreate,
		ReadContext:   resourceThreatInsightSettingsRead,
		UpdateContext: resourceThreatInsightSettingsUpdate,
		DeleteContext: resourceThreatInsightSettingsDelete,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
		Schema: map[string]*schema.Schema{
			"action": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: elemInSlice([]string{"none", "audit", "block"}),
				Description:      "Specifies how Okta responds to authentication requests from suspicious IPs",
			},
			"network_excludes": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of Network Zone IDs to exclude to be not logged or blocked by Okta ThreatInsight and proceed to Sign On rules evaluation",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceThreatInsightSettingsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conf, _, err := getOktaClientFromMetadata(m).ThreatInsightConfiguration.UpdateConfiguration(ctx, buildThreatInsightSettings(d))
	if err != nil {
		return diag.Errorf("failed to update threat insight configuration: %v", err)
	}
	d.SetId("threat_insight_settings")
	_ = d.Set("action", conf.Action)
	_ = d.Set("network_excludes", convertStringSliceToInterfaceSlice(conf.ExcludeZones))
	return nil
}

func resourceThreatInsightSettingsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conf, _, err := getOktaClientFromMetadata(m).ThreatInsightConfiguration.GetCurrentConfiguration(ctx)
	if err != nil {
		return diag.Errorf("failed to get threat insight configuration: %v", err)
	}
	d.SetId("threat_insight_settings")
	_ = d.Set("action", conf.Action)
	_ = d.Set("network_excludes", convertStringSliceToInterfaceSlice(conf.ExcludeZones))
	return nil
}

func resourceThreatInsightSettingsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	_, _, err := getOktaClientFromMetadata(m).ThreatInsightConfiguration.UpdateConfiguration(ctx, buildThreatInsightSettings(d))
	if err != nil {
		return diag.Errorf("failed to update threat insight configuration: %v", err)
	}
	return nil
}

func resourceThreatInsightSettingsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	_, _, err := getOktaClientFromMetadata(m).ThreatInsightConfiguration.UpdateConfiguration(ctx, okta.ThreatInsightConfiguration{Action: "none"})
	if err != nil {
		return diag.Errorf("failed to set default threat insight configuration: %v", err)
	}
	return nil
}

func buildThreatInsightSettings(d *schema.ResourceData) okta.ThreatInsightConfiguration {
	return okta.ThreatInsightConfiguration{
		Action:       d.Get("action").(string),
		ExcludeZones: convertInterfaceToStringArrNullable(d.Get("network_excludes")),
	}
}
