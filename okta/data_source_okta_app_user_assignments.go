package okta

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

func dataSourceAppUserAssignments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppUserAssignmentsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the Okta App being queried for groups",
				ForceNew:    true,
			},
			"users": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of user IDs assigned to the app",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(1 * time.Hour),
		},
	}
}

func dataSourceAppUserAssignmentsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(m)
	id := d.Get("id").(string)

	userAssignments, resp, err := client.Application.ListApplicationUsers(ctx, id, &query.Params{})
	if err != nil {
		return diag.Errorf("unable to query for users from app (%s): %s", id, err)
	}

	for resp.HasNextPage() {
		var moreAssignments []*okta.AppUser
		resp, err = resp.Next(ctx, &moreAssignments)
		if err != nil {
			return diag.Errorf("unable to query for users from app (%s): %s", id, err)
		}
		userAssignments = append(userAssignments, moreAssignments...)
	}

	var users []string
	for _, assignment := range userAssignments {
		users = append(users, assignment.Id)
	}
	_ = d.Set("users", convertStringSliceToSet(users))
	d.SetId(id)
	return nil
}
