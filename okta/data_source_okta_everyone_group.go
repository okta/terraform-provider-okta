package okta

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// data source to retrieve information on the Everyone Group

func dataSourceEveryoneGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceEveryoneGroupRead,

		Schema: map[string]*schema.Schema{
			"include_users": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Fetch group users, having default off cuts down on API calls.",
			},
		},
	}
}

func dataSourceEveryoneGroupRead(d *schema.ResourceData, m interface{}) error {
	return findGroup("Everyone", d, m)
}
