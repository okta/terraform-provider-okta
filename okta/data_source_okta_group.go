package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

func dataSourceGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGroupRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Type of the group. When specified in the terraform resource, will act as a filter when searching for the group",
				ValidateFunc: validation.StringInSlice([]string{"OKTA_GROUP", "APP_GROUP", "BUILT_IN"}, false),
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
	searchParams := &query.Params{Q: name}
	if d.Get("type") != nil && d.Get("type").(string) != "" {
		searchParams.Filter = fmt.Sprintf("type eq \"%s\"", d.Get("type").(string))
	}

	groups, _, err := client.Group.ListGroups(context.Background(), searchParams)
	if err != nil {
		return fmt.Errorf("failed to query for groups: %v", err)
	}

	if len(groups) < 1 {
		if d.Get("type") != nil {
			return fmt.Errorf("group \"%s\" was not found with type \"%s\"", name, d.Get("type").(string))
		}
		return fmt.Errorf("group \"%s\" was not found", name)
	}

	d.SetId(groups[0].Id)
	_ = d.Set("description", groups[0].Profile.Description)
	_ = d.Set("type", groups[0].Type)

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
