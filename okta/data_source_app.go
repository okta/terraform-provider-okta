package okta

import (
	"errors"

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
		},
	}
}

func dataSourceAppRead(d *schema.ResourceData, m interface{}) error {
	id := d.Get("id")
	label := d.Get("label")

	if id == "" && label == "" {
		return errors.New("you must provide either an id or label to search with")
	}
	appList, err := listApps(m.(*Config))
	if err != nil {
		return err
	}
}
