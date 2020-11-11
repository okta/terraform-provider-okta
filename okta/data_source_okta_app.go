package okta

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceApp() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAppRead,

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

func dataSourceAppRead(d *schema.ResourceData, m interface{}) error {
	filters, err := getAppFilters(d)
	if err != nil {
		return err
	}
	appList, err := listApps(m.(*Config), filters)
	if err != nil {
		return err
	}
	if len(appList) < 1 {
		return fmt.Errorf("no application found with provided filter: %s", filters)
	} else if len(appList) > 1 {
		fmt.Println("found multiple applications with the criteria supplied, using the first one, sorted by creation date.")
	}
	app := appList[0]
	d.SetId(app.ID)
	_ = d.Set("label", app.Label)
	_ = d.Set("description", app.Description)
	_ = d.Set("name", app.Name)
	_ = d.Set("status", app.Status)

	return nil
}
