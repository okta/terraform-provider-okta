package okta

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func Groups() *schema.Resource {
	return &schema.Resource{
		Create: GroupCreate,
		Read:   GroupRead,
		Update: GroupUpdate,
		Delete: GroupDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Group name",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Group description",
			},

			"group_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Group ID (generated)",
			},
		},
	}
}

func GroupCreate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func GroupRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func GroupUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func GroupDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
