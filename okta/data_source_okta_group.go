package okta

import (
	"errors"
	"fmt"

	"github.com/okta/okta-sdk-golang/okta/query"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGroupRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"include_users": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Fetch group users, having default off cuts down on API calls.",
			},
			"users": &schema.Schema{
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Users associated with the group. This can also be done per user.",
			},
		},
	}
}

func dataSourceGroupRead(d *schema.ResourceData, m interface{}) error {
	return findGroup(d.Get("name").(string), d, m)
}

func findGroup(name string, d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	groups, _, err := client.Group.ListGroups(&query.Params{Q: name})
	if err != nil {
		return fmt.Errorf("failed to query for groups: %v", err)
	} else if len(groups) < 1 {
		return errors.New("Group not found")
	}

	d.SetId(groups[0].Id)
	d.Set("description", groups[0].Profile.Description)

	if d.Get("include_users").(bool) {
		userIdList, err := listGroupUserIds(m, d.Id())
		if err != nil {
			return err
		}

		// just user ids for now
		return d.Set("users", convertStringSetToInterface(userIdList))
	}

	return nil
}
