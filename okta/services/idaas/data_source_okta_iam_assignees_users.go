package idaas

import (
	"context"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/config"
)

func dataSourceIamAssigneesUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIamAssigneesUsersRead,
		Schema: map[string]*schema.Schema{
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     200,
				Description: "Maximum number of users to return per page (default: 200).",
			},
			"users": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of users with role assignments.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for the user.",
						},
						"orn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Okta Resource Name (ORN) for the user.",
						},
					},
				},
			},
		},
		Description: "Get a list of users with role assignments. Note: This datasource may take some time to complete for organizations with many users due to pagination.",
	}
}

func dataSourceIamAssigneesUsersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client := config.OktaIDaaSClient.OktaSDKClientV5()

	limit := d.Get("limit").(int)
	var allUsers []map[string]interface{}
	after := ""
	for {
		resp, _, err := client.RoleAssignmentAPI.ListUsersWithRoleAssignments(ctx).After(after).Limit(int32(limit)).Execute()
		if err != nil {
			return diag.Errorf("failed to list users with role assignments: %v", err)
		}
		if resp == nil || resp.Value == nil {
			break
		}
		for _, user := range resp.Value {
			userMap := map[string]interface{}{
				"id":  user.GetId(),
				"orn": user.GetOrn(),
			}
			allUsers = append(allUsers, userMap)
		}
		// Pagination: check for next link
		if resp.Links == nil || resp.Links.Next == nil || resp.Links.Next.Href == "" {
			break
		}
		// Extract 'after' parameter from next link
		nextURL := resp.Links.Next.Href
		if nextURL != "" {
			parsedURL, err := url.Parse(nextURL)
			if err == nil {
				after = parsedURL.Query().Get("after")
				if after == "" {
					break
				}
			} else {
				break
			}
		} else {
			break
		}
	}
	if err := d.Set("users", allUsers); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("iam_assignees_users")
	return nil
}
