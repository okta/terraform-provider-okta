package idaas

import (
	"context"

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
	return resourceBehaviorCreateUsingSDK(ctx, d, meta)
}
