package okta

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceApp() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAppRead,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"label"},
			},
			"label": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id"},
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
	id := d.Get("id").(string)
	label := d.Get("label").(string)

	if id == "" && label == "" {
		return errors.New("you must provide either an id or label to search with")
	}
	appList, err := listApps(m.(*Config), &appFilters{ID: id, Label: label})
	if err != nil {
		return err
	}
	if len(appList) < 1 {
		return fmt.Errorf("No application found with provided id or label. ID: %s, Label %s", id, label)
	}
	app := appList[0]
	d.SetId(app.ID)
	d.Set("label", app.Label)
	d.Set("description", app.Description)
	d.Set("name", app.Name)
	d.Set("status", app.Status)

	return nil
}
