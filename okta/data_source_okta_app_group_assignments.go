package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/okta/terraform-provider-okta/sdk/query"
)

func dataSourceAppGroupAssignments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppGroupAssignmentsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the Okta App being queried for groups",
				ForceNew:    true,
			},
			"groups": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of groups IDs assigned to the app",
			},
		},
	}
}

func dataSourceAppGroupAssignmentsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(m)
	id := d.Get("id").(string)

	groupAssignments, resp, err := client.Application.ListApplicationGroupAssignments(ctx, id, &query.Params{})
	if err != nil {
		return diag.Errorf("unable to query for groups from app (%s): %s", id, err)
	}

	for {
		var moreAssignments []*sdk.ApplicationGroupAssignment
		if resp.HasNextPage() {
			resp, err = resp.Next(ctx, &moreAssignments)
			if err != nil {
				return diag.Errorf("unable to query for groups from app (%s): %s", id, err)
			}
			groupAssignments = append(groupAssignments, moreAssignments...)
		} else {
			break
		}
	}

	var groups []string
	for _, assignment := range groupAssignments {
		groups = append(groups, assignment.Id)
	}
	_ = d.Set("groups", convertStringSliceToSet(groups))
	d.SetId(id)
	return nil
}
