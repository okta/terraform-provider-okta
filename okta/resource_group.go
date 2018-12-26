package okta

import (
	"net/http"

	"github.com/okta/okta-sdk-golang/okta"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/okta/okta-sdk-golang/okta/query"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupCreate,
		Exists: resourceGroupExists,
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
		},
	}
}

func buildGroup(d *schema.ResourceData) *okta.Group {
	return &okta.Group{
		Profile: &okta.GroupProfile{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
		},
	}
}

func resourceGroupCreate(d *schema.ResourceData, m interface{}) error {
	group := buildGroup(d)
	responseGroup, _, err := getOktaClientFromMetadata(m).Group.CreateGroup(*group)
	if err != nil {
		return err
	}

	d.SetId(responseGroup.Id)

	return resourceGroupRead(d, m)
}

func resourceGroupExists(d *schema.ResourceData, m interface{}) (bool, error) {
	g, err := fetchGroup(d, m)

	return err == nil && g != nil, err
}

func resourceGroupRead(d *schema.ResourceData, m interface{}) error {
	g, err := fetchGroup(d, m)
	if err != nil {
		return err
	}

	d.Set("name", g.Profile.Name)
	d.Set("description", g.Profile.Description)

	return nil
}

func resourceGroupUpdate(d *schema.ResourceData, m interface{}) error {
	group := buildGroup(d)
	_, _, err := getOktaClientFromMetadata(m).Group.UpdateGroup(d.Id(), *group)
	if err != nil {
		return err
	}

	return resourceGroupRead(d, m)
}

func resourceGroupDelete(d *schema.ResourceData, m interface{}) error {
	_, err := getOktaClientFromMetadata(m).Group.DeleteGroup(d.Id())

	return err
}

func fetchGroup(d *schema.ResourceData, m interface{}) (*okta.Group, error) {
	g, resp, err := getOktaClientFromMetadata(m).Group.GetGroup(d.Id(), &query.Params{})

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	return g, err
}
