package okta

import (
	"github.com/hashicorp/terraform/helper/schema"
)

// data source to retrieve information on the Everyone Group

func dataSourceEveryoneGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceEveryoneGroupRead,

		Schema: map[string]*schema.Schema{},
	}
}

func dataSourceEveryoneGroupRead(d *schema.ResourceData, m interface{}) error {
	return findGroup("Everyone", d, m)
}
