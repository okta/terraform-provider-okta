package okta

import (
	"context"
	"fmt"

	"github.com/okta/okta-sdk-golang/v2/okta/query"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGroupRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"include_users": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Fetch group users, having default off cuts down on API calls.",
			},
			"users": {
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
	groups, _, err := client.Group.ListGroups(context.Background(), &query.Params{Q: name})
	if err != nil {
		return fmt.Errorf("failed to query for groups: %v", err)
	} else if len(groups) < 1 {
		return fmt.Errorf("group \"%s\" not found", name)
	}

	d.SetId(groups[0].Id)
	_ = d.Set("description", groups[0].Profile.Description)

	if d.Get("include_users").(bool) {
		userIDList, err := listGroupUserIDs(m, d.Id())
		if err != nil {
			return err
		}

		// just user ids for now
		return d.Set("users", convertStringSetToInterface(userIDList))
	}

	return nil
}
