package okta

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

func resourceOktaAppUserAssignments() *schema.Resource {
	return &schema.Resource{
		CreateContext: nil,
		ReadContext:   nil,
		UpdateContext: nil,
		DeleteContext: nil,
		Importer:      nil,
		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "App to associate user with",
			},
			"users": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Set of users to associate with the app",
				MinItems:    1,
				Set: func(v interface{}) int {
					buf := bytes.NewBuffer(nil)
					user := v.(map[string]interface{})

					buf.WriteString(fmt.Sprintf("%s-", user["id"].(string)))
					buf.WriteString(fmt.Sprintf("%s-", user["username"].(string)))
					buf.WriteString(fmt.Sprintf("%s-", normalizeDataJSON(user["profile"])))

					return schema.HashString(buf.String())
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "User associated with the application",
						},
						"username": {
							Type:     schema.TypeString,
							Required: true,
						},
						"password": {
							Type:      schema.TypeString,
							Sensitive: true,
							Optional:  true,
						},
						"profile": {
							Type:             schema.TypeString,
							ValidateDiagFunc: stringIsJSON,
							StateFunc:        normalizeDataJSON,
							Optional:         true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return new == ""
							},
						},
						"retain_assignment": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Retain the user assignment on destroy. If set to true, the resource will be removed from state but not from the Okta app.",
						},
					},
				},
			},
		},
	}
}

func resourceOktaAppUserAssignmentsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	users := d.Get("users").(*schema.Set).List()
	client := getOktaClientFromMetadata(m)
	appID := d.Get("app_id").(string)

	assignments := tfUsersToUserAssignments(users...)

	err := addUserAssignments(ctx, client, appID, assignments)
	if err != nil {
		return diag.FromErr(err)
	}

	//okta_app_user_assignments completely controls all assignments for an application
	d.SetId(appID)
	return nil //TODO: Use read function here
}

func tfUsersToUserAssignments(users ...interface{}) map[string]okta.AppUser {
	assignments := map[string]okta.AppUser{}

	for _, rawUser := range users {
		user := rawUser.(map[string]interface{})

		id := user["id"].(string)
		if id == "" {
			continue
		}

		username := user["username"].(string)

		password := user["password"].(string)

		rawProfile := user["profile"]
		var profile interface{}
		_ = json.Unmarshal([]byte(rawProfile.(string)), &profile)

		assignments[id] = okta.AppUser{
			Id:      id,
			Profile: profile,
			Credentials: &okta.AppUserCredentials{
				UserName: username,
				Password: &okta.AppUserPasswordCredential{
					Value: password,
				},
			},
		}
	}
	return assignments
}

func addUserAssignments(ctx context.Context, client *okta.Client, appID string, assignments map[string]okta.AppUser) error {
	for userID, assignment := range assignments {
		_, _, err := client.Application.AssignUserToApplication(ctx, appID, assignment)
		if err != nil {
			return fmt.Errorf("failed to assign user (%s) to app (%s): %s", userID, appID, err)
		}
	}
	return nil
}

func removeUserAssignments(ctx context.Context, client *okta.Client, appID string, assignments map[string]okta.AppUser) error {
	for userID, _ := range assignments {
		_, err := client.Application.DeleteApplicationUser(ctx, appID, userID, &query.Params{})
		if err != nil {
			return fmt.Errorf("failed to unassign user (%s) from app (%s): %s", userID, appID, err)
		}
	}
	return nil
}
