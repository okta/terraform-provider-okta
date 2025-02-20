package idaas

import (
	"context"
	"reflect"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
)

func ResourceThreatInsightSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceThreatInsightSettingsCreate,
		ReadContext:   resourceThreatInsightSettingsRead,
		UpdateContext: resourceThreatInsightSettingsUpdate,
		DeleteContext: resourceThreatInsightSettingsDelete,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
		Description:   "Manages Okta Threat Insight Settings. This resource allows you to configure Threat Insight Settings.",
		Schema: map[string]*schema.Schema{
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies how Okta responds to authentication requests from suspicious IPs. Valid values are `none`, `audit`, or `block`. A value of `none` indicates that ThreatInsight is disabled. A value of `audit` indicates that Okta logs suspicious requests in the System Log. A value of `block` indicates that Okta logs suspicious requests in the System Log and blocks the requests.",
			},
			"network_excludes": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Accepts a list of Network Zone IDs. Can only accept zones of `IP` type. IPs in the excluded Network Zones aren't logged or blocked by Okta ThreatInsight and proceed to Sign On rules evaluation. This ensures that traffic from known, trusted IPs isn't accidentally logged or blocked. The ordering of the network zone is not guarantee from the API sides",
				Elem:        &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// special case before the resource has been created
					if old == "" || k == "network_excludes.#" {
						// old value of an item is blank OR
						// "network_excludes.#"" will be seen during the set up before create
						// always return false or the TF plan will be in an illogical state
						return false
					}
					_o, _n := d.GetChange("network_excludes")
					oldNE, ok := _o.([]interface{})
					if !ok {
						return false
					}
					newNE, ok := _n.([]interface{})
					if !ok {
						return false
					}
					// length of new/old network_exclude slices has changed
					if len(oldNE) != len(newNE) {
						return false
					}
					// sort and check with deepequal
					sort.Slice(oldNE, func(i, j int) bool {
						a := oldNE[i].(string)
						b := oldNE[j].(string)
						return a < b
					})
					sort.Slice(newNE, func(i, j int) bool {
						a := newNE[i].(string)
						b := newNE[j].(string)
						return a < b
					})
					return reflect.DeepEqual(oldNE, newNE)
				},
			},
		},
	}
}

func resourceThreatInsightSettingsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, _, err := GetOktaClientFromMetadata(meta).ThreatInsightConfiguration.UpdateConfiguration(ctx, buildThreatInsightSettings(d))
	if err != nil {
		return diag.Errorf("failed to update threat insight configuration: %v", err)
	}
	d.SetId("threat_insight_settings")
	_ = d.Set("action", conf.Action)
	_ = d.Set("network_excludes", utils.ConvertStringSliceToInterfaceSlice(conf.ExcludeZones))
	return nil
}

func resourceThreatInsightSettingsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, _, err := GetOktaClientFromMetadata(meta).ThreatInsightConfiguration.GetCurrentConfiguration(ctx)
	if err != nil {
		return diag.Errorf("failed to get threat insight configuration: %v", err)
	}
	d.SetId("threat_insight_settings")
	_ = d.Set("action", conf.Action)
	_ = d.Set("network_excludes", utils.ConvertStringSliceToInterfaceSlice(conf.ExcludeZones))
	return nil
}

func resourceThreatInsightSettingsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_, _, err := GetOktaClientFromMetadata(meta).ThreatInsightConfiguration.UpdateConfiguration(ctx, buildThreatInsightSettings(d))
	if err != nil {
		return diag.Errorf("failed to update threat insight configuration: %v", err)
	}
	return nil
}

func resourceThreatInsightSettingsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_, _, err := GetOktaClientFromMetadata(meta).ThreatInsightConfiguration.UpdateConfiguration(ctx, sdk.ThreatInsightConfiguration{Action: "none"})
	if err != nil {
		return diag.Errorf("failed to set default threat insight configuration: %v", err)
	}
	return nil
}

func buildThreatInsightSettings(d *schema.ResourceData) sdk.ThreatInsightConfiguration {
	return sdk.ThreatInsightConfiguration{
		Action:       d.Get("action").(string),
		ExcludeZones: utils.ConvertInterfaceToStringArrNullable(d.Get("network_excludes")),
	}
}
