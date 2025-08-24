package idaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBehavior() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBehaviorReadUsingSDK,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"name"},
				Description:   "Behavior ID.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Behavior name.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Behavior status.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Behavior type.",
			},
			"settings": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "Map of behavior settings.",
			},
		},
		Description: "Get a behavior by name or ID.",
	}
}

func dataSourceBehaviorReadUsingSDK(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_, idExists := d.GetOk("id")
	if idExists {
		return resourceBehaviorReadUsingSDK(ctx, d, meta)
	}
	name, nameExists := d.GetOk("name")
	if !nameExists {
		return diag.Errorf("either id or name must be specified")
	}
	behaviorRulesFromRawResp, err := getBehaviorRules(ctx, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	for _, behaviorRule := range behaviorRulesFromRawResp {
		settingsMap := make(map[string]any)
		for k, v := range behaviorRule["settings"].(map[string]any) {
			settingsMap[k] = fmt.Sprint(v)
		}
		if name == behaviorRule["name"] {
			d.Set("name", behaviorRule["name"])
			d.Set("type", behaviorRule["type"])
			d.Set("status", behaviorRule["status"])
			d.Set("settings", settingsMap)
		}
	}
	return nil
}
