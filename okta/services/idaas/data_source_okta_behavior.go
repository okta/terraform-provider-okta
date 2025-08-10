package idaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/okta/terraform-provider-okta/sdk/query"
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
	var behavior *sdk.Behavior
	behaviorID, ok := d.GetOk("id")
	if ok {
		respBehavior, _, err := getAPISupplementFromMetadata(meta).GetBehavior(ctx, behaviorID.(string))
		if err != nil {
			return diag.Errorf("failed get behavior by ID: %v", err)
		}
		behavior = respBehavior
	} else {
		name := d.Get("name").(string)
		searchParams := &query.Params{Q: name, Limit: 1}
		logger(meta).Info("looking for behavior", "query", searchParams.String())
		behaviors, _, err := getAPISupplementFromMetadata(meta).ListBehaviors(ctx, searchParams)
		switch {
		case err != nil:
			return diag.Errorf("failed to query for behaviors: %v", err)
		case len(behaviors) < 1:
			return diag.Errorf("behavior with name '%s' does not exist", name)
		case behaviors[0].Name != name:
			logger(meta).Warn("behavior with exact name match was not found: using partial match which contains name as a substring", "name", behaviors[0].Name)
		}
		behavior = behaviors[0]
	}
	d.SetId(behavior.ID)
	_ = d.Set("type", behavior.Type)
	_ = d.Set("status", behavior.Status)
	settings := make(map[string]string)
	for k, v := range behavior.Settings {
		settings[k] = fmt.Sprint(v)
	}
	_ = d.Set("settings", settings)
	return nil
}

// func dataSourceBehaviorReadUsingSDK(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
// 	var behavior *v5okta.ListBehaviorDetectionRules200ResponseInner
// 	behaviorID, ok := d.GetOk("id")
// 	if ok {
// 		respBehavior, _, err := getOktaV5ClientFromMetadata(meta).BehaviorAPI.GetBehaviorDetectionRule(ctx, behaviorID.(string)).Execute()
// 		if err != nil {
// 			return diag.Errorf("failed get behavior by ID: %v", err)
// 		}
// 		behavior = respBehavior
// 	} else {
// 		name := d.Get("name").(string)
// 		searchParams := &query.Params{Q: name, Limit: 1}
// 		logger(meta).Info("looking for behavior", "query", searchParams.String())
// 		listBehaviorDetectionRules := getOktaV5ClientFromMetadata(meta).BehaviorAPI.ListBehaviorDetectionRules(ctx)
// 		behaviors, _, err := listBehaviorDetectionRules.Execute()
// 		switch {
// 		case err != nil:
// 			return diag.Errorf("failed to query for behaviors: %v", err)
// 		case len(behaviors) < 1:
// 			return diag.Errorf("behavior with name '%s' does not exist", name)
// 		case behaviors[0].Name != name:
// 			logger(meta).Warn("behavior with exact name match was not found: using partial match which contains name as a substring", "name", behaviors[0].Name)
// 		}
// 		behavior = behaviors[0]
// 	}
// 	d.SetId(behavior.ID)
// 	_ = d.Set("type", behavior.Type)
// 	_ = d.Set("status", behavior.Status)
// 	settings := make(map[string]string)
// 	for k, v := range behavior.Settings {
// 		settings[k] = fmt.Sprint(v)
// 	}
// 	_ = d.Set("settings", settings)
// 	return nil
// }
