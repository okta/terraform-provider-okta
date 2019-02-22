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
	labelPrefix := d.Get("label_prefix").(string)

	if id == "" && label == "" && labelPrefix == "" {
		return errors.New("you must provide either an label_prefix, id, or label to search with")
	}
	appList, err := listApps(m.(*Config), &appFilters{ID: id, Label: label, LabelPrefix: labelPrefix})
	if err != nil {
		return err
	}
	if len(appList) < 1 {
		return fmt.Errorf(`No application found with provided filter. id: "%s", label: "%s", label_prefix: "%s"`, id, label, labelPrefix)
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
