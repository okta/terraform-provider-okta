package idaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oktav5sdk "github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/okta/terraform-provider-okta/sdk/query"
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

func dataSourceBehaviorReadUsingSDK(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var result any
	behaviorID, ok := d.GetOk("id")
	name := d.Get("name").(string)
	if ok { // check if "id" was provided in the resource's configuration by user
		behavior, _, err := getOktaV5ClientFromMetadata(meta).BehaviorAPI.GetBehaviorDetectionRule(ctx, behaviorID.(string)).Execute()
		if err != nil {
			return diag.Errorf("failed get behavior by ID: %v", err)
		}
		result = behavior.GetActualInstance()
	} else if name != "" { // list all behavior detection rules and search the list for the given name, query params aren't supported by SDK
		listBehaviorDetectionRules := getOktaV5ClientFromMetadata(meta).BehaviorAPI.ListBehaviorDetectionRules(ctx)
		behaviors, _, err := listBehaviorDetectionRules.Execute()
		if err != nil {
			return diag.Errorf("failed to list behaviors: %v", err)
		}
		if len(behaviors) == 0 {
			return diag.Errorf("no behaviors seem to exist")
		}
		found := false
		for _, behavior := range behaviors {
			switch {
			case behavior.BehaviorRuleAnomalousDevice != nil && behavior.BehaviorRuleAnomalousDevice.GetName() == name:
				found = true
			case behavior.BehaviorRuleAnomalousIP != nil && behavior.BehaviorRuleAnomalousIP.GetName() == name:
				found = true
			case behavior.BehaviorRuleAnomalousLocation != nil && behavior.BehaviorRuleAnomalousLocation.GetName() == name:
				found = true
			case behavior.BehaviorRuleVelocity != nil && behavior.BehaviorRuleVelocity.GetName() == name:
				found = true
			}
			if found {
				result = behavior.GetActualInstance()
				break
			}
		}
		if !found {
			return diag.Errorf("behavior with name '%s' does not exist", name)
		}
	} else {
		return diag.Errorf("neither name nor id were provided")
	}

	switch concreteType := result.(type) {
	case oktav5sdk.BehaviorRuleAnomalousDevice:
		d.Set("type", concreteType.GetType())
		d.Set("status", concreteType.GetStatus())
		d.Set("settings", map[string]any{
			"maxEventsUsedForEvaluation": *concreteType.GetSettings().MaxEventsUsedForEvaluation,
			"minEventsUsedForEvaluation": *concreteType.GetSettings().MinEventsNeededForEvaluation,
		})
	case oktav5sdk.BehaviorRuleAnomalousLocation:
		d.Set("type", concreteType.GetType())
		d.Set("status", concreteType.GetStatus())
		d.Set("settings", map[string]any{
			"maxEventsUsedForEvaluation": *concreteType.GetSettings().MaxEventsUsedForEvaluation,
			"minEventsUsedForEvaluation": *concreteType.GetSettings().MinEventsNeededForEvaluation,
			"radiusKilometers":           *concreteType.GetSettings().RadiusKilometers,
			"granularity":                concreteType.GetSettings().Granularity,
		})
	case oktav5sdk.BehaviorRuleAnomalousIP:
		d.Set("type", concreteType.GetType())
		d.Set("status", concreteType.GetStatus())
		d.Set("settings", map[string]any{
			"maxEventsUsedForEvaluation": *concreteType.GetSettings().MaxEventsUsedForEvaluation,
			"minEventsUsedForEvaluation": *concreteType.GetSettings().MinEventsNeededForEvaluation,
		})
	case oktav5sdk.BehaviorRuleVelocity:
		d.Set("type", concreteType.GetType())
		d.Set("status", concreteType.GetStatus())
		d.Set("settings", map[string]any{
			"velocityKph": concreteType.GetSettings().VelocityKph,
		})
	}
	return nil
}
