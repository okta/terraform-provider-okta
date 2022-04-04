package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/okta/terraform-provider-okta/sdk"
)

func dataSourceBehavior() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBehaviorRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"settings": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
		},
	}
}

func dataSourceBehaviorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var behavior *sdk.Behavior
	behaviorID, ok := d.GetOk("id")
	if ok {
		respBehavior, _, err := getSupplementFromMetadata(m).GetBehavior(ctx, behaviorID.(string))
		if err != nil {
			return diag.Errorf("failed get behavior by ID: %v", err)
		}
		behavior = respBehavior
	} else {
		name := d.Get("name").(string)
		searchParams := &query.Params{Q: name, Limit: 1}
		logger(m).Info("looking for behavior", "query", searchParams.String())
		behaviors, _, err := getSupplementFromMetadata(m).ListBehaviors(ctx, searchParams)
		switch {
		case err != nil:
			return diag.Errorf("failed to query for behaviors: %v", err)
		case len(behaviors) < 1:
			return diag.Errorf("behavior with name '%s' does not exist", name)
		case behaviors[0].Name != name:
			logger(m).Warn("behavior with exact name match was not found: using partial match which contains name as a substring", "name", behaviors[0].Name)
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
