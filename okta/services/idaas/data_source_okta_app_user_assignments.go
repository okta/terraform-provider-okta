package idaas

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/okta/terraform-provider-okta/sdk/query"
)

// dataSourceAppUserAssignments returns a Terraform data source for retrieving
// detailed information about users assigned to an Okta application.
// This data source exposes the full API response including user status, scope,
// profile data, credentials, and timestamps.
func dataSourceAppUserAssignments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppUserAssignmentsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the Okta App being queried for users",
				ForceNew:    true,
			},
			"users": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "List of users assigned to the app with detailed information",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique identifier for the Okta User",
						},
						"external_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the user in the target app that's linked to the Okta Application User object",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of an Application User (ACTIVE, INACTIVE, PROVISIONED, etc.)",
						},
						"scope": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates if the assignment is direct (USER) or by group membership (GROUP)",
						},
						"sync_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The synchronization state for the Application User (DISABLED, ENABLED, ERROR, etc.)",
						},
						"created": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Timestamp when the Application User was created",
						},
						"last_updated": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Timestamp when the Application User was last updated",
						},
						"last_sync": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Timestamp of the last synchronization operation",
						},
						"password_changed": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Timestamp when the Application User password was last changed",
						},
						"status_changed": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Timestamp when the Application User status was last changed",
						},
						"profile": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Profile properties for the user in this application",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"credentials": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Credentials for the Application User",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"user_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Username for the Application User",
									},
									"password": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Password for the Application User",
										Sensitive:   true,
									},
								},
							},
						},
					},
				},
			},
		},
		Description: "Get a set of users assigned to an Okta application with detailed information.",
	}
}

// dataSourceAppUserAssignmentsRead retrieves all users assigned to an Okta application
// and returns detailed information about each user including their status, scope,
// profile data, credentials, and various timestamps.
func dataSourceAppUserAssignmentsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(meta)
	id := d.Get("id").(string)

	userAssignments, resp, err := client.Application.ListApplicationUsers(ctx, id, &query.Params{})
	if err != nil {
		return diag.Errorf("unable to query for users from app (%s): %s", id, err)
	}

	for resp.HasNextPage() {
		var moreAssignments []*sdk.AppUser
		resp, err = resp.Next(ctx, &moreAssignments)
		if err != nil {
			return diag.Errorf("unable to query for users from app (%s): %s", id, err)
		}
		userAssignments = append(userAssignments, moreAssignments...)
	}

	var users []map[string]interface{}
	for _, assignment := range userAssignments {
		user := processAppUserAssignment(assignment)
		users = append(users, user)
	}

	_ = d.Set("users", users)
	d.SetId(id)
	return nil
}

// processAppUserAssignment converts an AppUser to a Terraform-compatible map
func processAppUserAssignment(assignment *sdk.AppUser) map[string]interface{} {
	user := map[string]interface{}{
		"id": assignment.Id,
	}

	if assignment.ExternalId != "" {
		user["external_id"] = assignment.ExternalId
	}
	if assignment.Status != "" {
		user["status"] = assignment.Status
	}
	if assignment.Scope != "" {
		user["scope"] = assignment.Scope
	}
	if assignment.SyncState != "" {
		user["sync_state"] = assignment.SyncState
	}
	if assignment.Created != nil {
		user["created"] = assignment.Created.Format(time.RFC3339)
	}
	if assignment.LastUpdated != nil {
		user["last_updated"] = assignment.LastUpdated.Format(time.RFC3339)
	}
	if assignment.LastSync != nil {
		user["last_sync"] = assignment.LastSync.Format(time.RFC3339)
	}
	if assignment.PasswordChanged != nil {
		user["password_changed"] = assignment.PasswordChanged.Format(time.RFC3339)
	}
	if assignment.StatusChanged != nil {
		user["status_changed"] = assignment.StatusChanged.Format(time.RFC3339)
	}

	// Handle profile data
	if assignment.Profile != nil {
		if profileMap, ok := assignment.Profile.(map[string]interface{}); ok {
			// Convert profile values to strings for Terraform compatibility
			stringProfile := make(map[string]string)
			for k, v := range profileMap {
				if v == nil {
					continue // Skip null values
				}
				if str, ok := v.(string); ok {
					stringProfile[k] = str
				} else {
					// Convert other types to string representation
					stringProfile[k] = fmt.Sprintf("%v", v)
				}
			}
			if len(stringProfile) > 0 {
				user["profile"] = stringProfile
			}
		}
	}

	// Handle credentials
	if assignment.Credentials != nil {
		credentials := map[string]interface{}{}
		if assignment.Credentials.UserName != "" {
			credentials["user_name"] = assignment.Credentials.UserName
		}
		if assignment.Credentials.Password != nil && assignment.Credentials.Password.Value != "" {
			credentials["password"] = assignment.Credentials.Password.Value
		}
		if len(credentials) > 0 {
			user["credentials"] = []map[string]interface{}{credentials}
		}
	}

	return user
}
