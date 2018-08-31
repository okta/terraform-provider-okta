package okta

import (
	"fmt"
	"github.com/articulate/oktasdk-go/okta"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupCreate,
		Read:   resourceGroupRead,
		Update: resourceGroupUpdate,
		Delete: resourceGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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

func assembleGroup() *okta.Group {
	group := &okta.Group{}
	profile := &okta.GroupProfile{}
	links := &okta.GroupLinks{}

	group.GroupProfile = profile
	group.GroupLinks = links

	return group
}

func populateGroup(group *okta.Group, d *schema.ResourceData) error {
	group.GroupProfile.Name = d.Get("name").(string)
	group.GroupProfile.Description = d.Get("description").(string)

	return nil
}

func resourceGroupCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Config).oktaClient
	
	group := assembleGroup()
	populateGroup(group, d)

	returnedGroup, _, err := client.Groups.Add(group.GroupProfile.Name, group.GroupProfile.Description)

	if err != nil {
		return fmt.Errorf("[ERROR] %v.", err)
	}

	d.SetId(returnedGroup.ID)
	
	return nil
}

func resourceGroupRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceGroupUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Config).oktaClient

	group := assembleGroup()
	populateGroup(group, d)

	_, _, err := client.Groups.Update(d.Id(), group)

	if err != nil {
		return fmt.Errorf("[ERROR] %v.", err)
	}

	return nil
}

func resourceGroupDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Config).oktaClient

	_, err := client.Groups.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("[ERROR] %v.", err)
	}

	d.SetId("")
	
	return nil
}
