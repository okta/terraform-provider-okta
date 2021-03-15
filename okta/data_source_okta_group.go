package okta

import (
	"context"
	"fmt"

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
			"type": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Type of the group. When specified in the terraform resource, will act as a filter when searching for the group",
				ValidateDiagFunc: stringInSlice([]string{"OKTA_GROUP", "APP_GROUP", "BUILT_IN"}),
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
	return findGroup(ctx, d.Get("name").(string), d, m, false)
}

func findGroup(ctx context.Context, name string, d *schema.ResourceData, m interface{}, isEveryone bool) diag.Diagnostics {
	client := getOktaClientFromMetadata(m)
	searchParams := &query.Params{Q: name}
	t, okType := d.GetOk("type")
	if okType {
		searchParams.Filter = fmt.Sprintf("type eq \"%s\"", t.(string))
	}
	logger(m).Info("looking for data source group", "query", searchParams.String())
	groups, _, err := client.Group.ListGroups(ctx, searchParams)
	switch {
	case err != nil:
		return diag.Errorf("failed to query for groups: %v", err)
	case len(groups) < 1:
		if okType {
			return diag.Errorf("group with name '%s' and type '%s' does not exist", name, d.Get("type").(string))
		}
		return diag.Errorf("group with name '%s' does not exist", name)
	case len(groups) > 1:
		// TODO try to find exact match
		logger(m).Warn("Found multiple groups with the supplied parameters: using the first one which may only be a partial match", "name", groups[0].Profile.Name)
	case len(groups[0].Profile.Name) != len(name):
		logger(m).Warn("The group with an exact match to the supplied name was not found: using partial match which contains name as a substring", "name", groups[0].Profile.Name)
	}
	d.SetId(groups[0].Id)
	_ = d.Set("description", groups[0].Profile.Description)
	if !isEveryone {
		_ = d.Set("type", groups[0].Type)
	}
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
