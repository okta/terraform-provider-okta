package okta

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceApp() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAppRead,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"label", "label_prefix"},
			},
			"label": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id", "label_prefix"},
			},
			"label_prefix": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id", "label"},
			},
			"active_only": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Search only ACTIVE applications.",
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": &schema.Schema{
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
		return fmt.Errorf("No application found with provided filter: %s", filters)
	} else if len(appList) > 1 {
		fmt.Println("Found multiple applications with the criteria supplied, using the first one, sorted by creation date.")
	}
	app := appList[0]
	d.SetId(app.ID)
	d.Set("label", app.Label)
	d.Set("description", app.Description)
	d.Set("name", app.Name)
	d.Set("status", app.Status)

	return nil
}
