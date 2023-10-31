package okta

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/okta/terraform-provider-okta/sdk/query"
)

func dataSourceGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGroupRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeString,
				Computed:      true,
				Optional:      true,
				ConflictsWith: []string{"name", "type"},
				Description:   "ID of group.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Name of group.",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Type of the group. When specified in the terraform resource, will act as a filter when searching for the group",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of group.",
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
			"delay_read_seconds": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Force delay of the group read by N seconds. Useful when eventual consistency of group information needs to be allowed for; for instance, when group rules are known to have been applied.",
			},
		},
		Description: "Get a group from Okta.",
	}
}

func dataSourceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if n, ok := d.GetOk("delay_read_seconds"); ok {
		delay, err := strconv.Atoi(n.(string))
		if err == nil {
			logger(m).Info("delaying group read by ", delay, " seconds")
			m.(*Config).timeOperations.Sleep(time.Duration(delay) * time.Second)
		} else {
			logger(m).Warn("group read delay value ", n, " is not an integer")
		}
	}

	return findGroup(ctx, d.Get("name").(string), d, m, false)
}

func findGroup(ctx context.Context, name string, d *schema.ResourceData, m interface{}, isEveryone bool) diag.Diagnostics {
	var group *sdk.Group
	groupID, ok := d.GetOk("id")
	if ok {
		respGroup, _, err := getOktaClientFromMetadata(m).Group.GetGroup(ctx, groupID.(string))
		if err != nil {
			return diag.Errorf("failed get group by ID: %v", err)
		}
		group = respGroup
	} else {
		// NOTE: Okta API query on name is effectively a starts with query, not
		// an exact match query.
		searchParams := &query.Params{Q: name}

		// NOTE: Uncertain of the OKTA /api/v1/groups API drifted during Classic
		// (when this data source was originally created) to OIE migration.
		// Currently, Okta API enforces unique names on all groups regardless of
		// type so type is essentially a meaningless parameter in OIE.
		t, okType := d.GetOk("type")
		if okType {
			searchParams.Filter = fmt.Sprintf("type eq \"%s\"", t.(string))
		}

		logger(m).Info("looking for data source group", "query", searchParams.String())
		groups, _, err := getOktaClientFromMetadata(m).Group.ListGroups(ctx, searchParams)
		switch {
		case err != nil:
			return diag.Errorf("failed to query for groups: %v", err)
		case len(groups) > 1:
			if okType {
				return diag.Errorf("group starting with name %q and type %q matches %d groups, select a more precise name parameter", name, d.Get("type").(string), len(groups))
			}
			return diag.Errorf("group starting with name %q matches %d groups, select a more precise name parameter", name, len(groups))
		case len(groups) < 1:
			if okType {
				return diag.Errorf("group with name %q and type %q does not exist", name, d.Get("type").(string))
			}
			return diag.Errorf("group with name %q does not exist", name)
		case groups[0].Profile.Name != name:
			logger(m).Warn("group with exact name match was not found: using partial match which contains name as a substring", "name", groups[0].Profile.Name)
		}
		group = groups[0]
	}
	d.SetId(group.Id)
	_ = d.Set("description", group.Profile.Description)
	if !isEveryone {
		_ = d.Set("type", group.Type)
		_ = d.Set("name", group.Profile.Name)
	}
	if !d.Get("include_users").(bool) {
		return nil
	}
	userIDList, err := listGroupUserIDs(ctx, m, d.Id())
	if err != nil {
		return diag.Errorf("failed to list group user IDs: %v", err)
	}
	_ = d.Set("users", convertStringSliceToSet(userIDList))
	return nil
}
