package okta

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/okta/terraform-provider-okta/sdk/query"
)

func dataSourceAppGroupAssignments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppGroupAssignmentsRead,
		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the Okta App being queried for groups",
				ForceNew:    true,
			},
			"groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of groups IDs assigned to the app",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of the group associated with the application",
						},
						"priority": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Priority of group assignment",
						},
						"profile": {
							Type:        schema.TypeString,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "JSON document containing the assigned group's [profile](https://developer.okta.com/docs/reference/api/apps/#profile-object)",
						},
					},
				},
			},
		},
		Description: "Get a set of groups assigned to an Okta application.",
	}
}

func dataSourceAppGroupAssignmentsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(m)
	appId := d.Get("app_id").(string)

	groupAssignments, resp, err := client.Application.ListApplicationGroupAssignments(ctx, appId, &query.Params{})
	if err != nil {
		return diag.Errorf("unable to query for groups from app (%s): %s", appId, err)
	}

	for resp.HasNextPage() {
		var additionalGroups []*sdk.ApplicationGroupAssignment
		resp, err = resp.Next(ctx, &additionalGroups)
		if err != nil {
			return diag.Errorf("unable to query for groups from app (%s): %s", appId, err)
		}
		groupAssignments = append(groupAssignments, additionalGroups...)
	}

	groups := make([]map[string]interface{}, len(groupAssignments))
	for i, group := range groupAssignments {
		profile, err := json.Marshal(group.Profile)
		if err != nil {
			return diag.Errorf("unable to marshal app group profile: %s", err)
		}
		groups[i] = map[string]interface{}{
			"id":       group.Id,
			"priority": group.Priority,
			// "profile":  group.Profile,
			"profile": string(profile),
		}
	}

	d.Set("groups", groups)
	d.SetId(appId)
	return nil
}
