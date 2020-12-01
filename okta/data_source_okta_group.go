package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

func dataSourceGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGroupRead,
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

func dataSourceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return findGroup(ctx, d.Get("name").(string), d, m)
}

func findGroup(ctx context.Context, name string, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(m)
	groups, _, err := client.Group.ListGroups(ctx, &query.Params{Q: name})
	if err != nil {
		return diag.Errorf("failed to query for groups: %v", err)
	} else if len(groups) < 1 {
		return diag.Errorf("group with name '%s' does not exist", name)
	}

	d.SetId(groups[0].Id)
	_ = d.Set("description", groups[0].Profile.Description)

	if d.Get("include_users").(bool) {
		userIDList, err := listGroupUserIDs(ctx, m, d.Id())
		if err != nil {
			return diag.Errorf("failed to list group user IDs: %v", err)
		}
		// just user ids for now
		err = d.Set("users", convertStringSetToInterface(userIDList))
		if err != nil {
			return diag.FromErr(err)
		}
		return nil
	}
	return nil
}
