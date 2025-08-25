package idaas

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBehavior() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBehaviorRead,
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

func dataSourceBehaviorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_, idExists := d.GetOk("id")
	if idExists {
		d.SetId(fmt.Sprint(d.Get("id")))
		return resourceBehaviorRead(ctx, d, meta)
	}
	partialRes, partialFound := make(map[string]any), false
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
			d.SetId(behaviorRule["id"].(string))
			return nil // we already found our behavior
		}
		if !partialFound && strings.Contains(behaviorRule["name"].(string), name.(string)) {
			partialFound = true
			partialRes = map[string]any{
				"id":       behaviorRule["id"],
				"name":     behaviorRule["name"],
				"type":     behaviorRule["type"],
				"status":   behaviorRule["status"],
				"settings": settingsMap,
			}
			// cannot break since we might find an exact match later on.
		}
	}
	if !partialFound {
		return diag.Errorf("behavior with name '%s' does not exist", name)
	}
	logger(meta).Warn("behavior with exact name match was not found: using partial match which contains name as a substring", "name", partialRes["name"])
	d.Set("id", partialRes["id"])
	d.Set("name", partialRes["name"])
	d.Set("type", partialRes["type"])
	d.Set("status", partialRes["status"])
	d.Set("settings", partialRes["settings"])
	return nil
}
