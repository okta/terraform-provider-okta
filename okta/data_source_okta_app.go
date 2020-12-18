package okta

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceApp() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"label", "label_prefix"},
			},
			"label": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id", "label_prefix"},
			},
			"label_prefix": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id", "label"},
			},
			"active_only": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Search only ACTIVE applications.",
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAppRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	filters, err := getAppFilters(d)
	if err != nil {
		return diag.Errorf("failed to get filters: %v", err)
	}
	appList, err := listApps(ctx, m.(*Config), filters)
	if err != nil {
		return diag.Errorf("failed to list apps: %v", err)
	}
	if len(appList) < 1 {
		return diag.Errorf("no application found with provided filter: %+v", filters)
	} else if len(appList) > 1 {
		log.Print("found multiple applications with the criteria supplied, using the first one, sorted by creation date.\n")
	}
	app := appList[0]
	d.SetId(app.ID)
	_ = d.Set("label", app.Label)
	_ = d.Set("description", app.Description)
	_ = d.Set("name", app.Name)
	_ = d.Set("status", app.Status)

	return nil
}
