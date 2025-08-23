package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oktav5sdk "github.com/okta/okta-sdk-golang/v5/okta"
)

func dataSourceBehaviors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBehaviorsReadUsingSDK,
		Schema: map[string]*schema.Schema{
			"q": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Searches the name property of behaviors for matching value",
			},
			"behaviors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Behavior ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
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
							Type:        schema.TypeMap,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "Map of behavior settings.",
						},
					},
				},
			},
		},
		Description: "Get a behaviors by search criteria.",
	}
}

func dataSourceBehaviorsReadUsingSDK(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	listBehaviorDetectionRules := getOktaV5ClientFromMetadata(meta).BehaviorAPI.ListBehaviorDetectionRules(ctx)
	behaviors, _, err := listBehaviorDetectionRules.Execute()
	if err != nil {
		return diag.Errorf("failed to list behaviors: %v", err)
	}
	arr := make([]map[string]any, len(behaviors))
	for i, behavior := range behaviors {

		switch concreteType := behavior.GetActualInstance().(type) {
		case oktav5sdk.BehaviorRuleAnomalousDevice:
			arr[i] = map[string]any{
				"type":   concreteType.GetType(),
				"status": concreteType.GetStatus(),
				"settings": map[string]any{
					"maxEventsUsedForEvaluation": *concreteType.GetSettings().MaxEventsUsedForEvaluation,
					"minEventsUsedForEvaluation": *concreteType.GetSettings().MinEventsNeededForEvaluation,
				},
			}

		case oktav5sdk.BehaviorRuleAnomalousLocation:
			arr[i] = map[string]any{
				"type":   concreteType.GetType(),
				"status": concreteType.GetStatus(),
				"settings": map[string]any{
					"maxEventsUsedForEvaluation": *concreteType.GetSettings().MaxEventsUsedForEvaluation,
					"minEventsUsedForEvaluation": *concreteType.GetSettings().MinEventsNeededForEvaluation,
					"radiusKilometers":           *concreteType.GetSettings().RadiusKilometers,
					"granularity":                concreteType.GetSettings().Granularity,
				},
			}
		case oktav5sdk.BehaviorRuleAnomalousIP:
			arr[i] = map[string]any{
				"type":   concreteType.GetType(),
				"status": concreteType.GetStatus(),
				"settings": map[string]any{
					"maxEventsUsedForEvaluation": *concreteType.GetSettings().MaxEventsUsedForEvaluation,
					"minEventsUsedForEvaluation": *concreteType.GetSettings().MinEventsNeededForEvaluation,
				},
			}
		case oktav5sdk.BehaviorRuleVelocity:
			arr[i] = map[string]any{
				"type":   concreteType.GetType(),
				"status": concreteType.GetStatus(),
				"settings": map[string]any{
					"velocityKph": concreteType.GetSettings().VelocityKph,
				},
			}
		}
		err = d.Set("behaviors", arr)
	}
	return diag.FromErr(err)
}
