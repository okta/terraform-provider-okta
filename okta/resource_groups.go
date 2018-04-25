package okta

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGroups() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupCreate,
		Read:   resourceGroupRead,
		Update: resourceGroupUpdate,
		Delete: resourceGroupDelete,

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

func resourceGroupCreate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceGroupRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceGroupUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceGroupDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
