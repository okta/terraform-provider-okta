package okta

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

func resourceGroupMemberships() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupMembershipsCreate,
		ReadContext:   resourceGroupMembershipsRead,
		UpdateContext: resourceGroupMembershipsUpdate,
		DeleteContext: resourceGroupMembershipsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "Resource to manage a set of group memberships for a specific group.",
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of a Okta group.",
				ForceNew:    true,
			},
			"users": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "The list of Okta user IDs which the group should have membership managed for.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceGroupMembershipsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	groupId := d.Get("group_id").(string)
	users := convertInterfaceToStringSetNullable(d.Get("users"))
	client := getOktaClientFromMetadata(m)

	if len(users) == 0 {
		d.SetId(groupId)
		return nil
	}
	err := addGroupMembers(ctx, client, groupId, users)
	if err != nil {
		return diag.FromErr(err)
	}
	bOff := backoff.NewExponentialBackOff()
	bOff.MaxElapsedTime = time.Second * 10
	bOff.InitialInterval = time.Second

	// During create the Okta service can have eventual consistency issues when
	// adding users to a group. Use a backoff to wait for at list one user to be
	// associated with the group.
	err = backoff.Retry(func() error {
		ok, err := checkIfGroupHasUsers(ctx, client, groupId, users)
		if err != nil {
			return backoff.Permanent(err)
		}
		if ok {
			return nil
		}
		return fmt.Errorf("group (%s) did not have expected user memberships after multiple checks", groupId)
	}, bOff)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(groupId)
	return nil
}

func resourceGroupMembershipsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	groupId := d.Get("group_id").(string)
	client := getOktaClientFromMetadata(m)
	expected_users := convertInterfaceToStringSetNullable(d.Get("users"))

	// handle import edge case
	if len(expected_users) == 0 {
		logger(m).Info("reading group membership", "id", d.Id())
		userIDList, err := listGroupUserIDs(ctx, m, d.Id())
		if err != nil {
			return diag.Errorf("unable to return membership for group (%s) from API", d.Id())
		}
		d.Set("group_id", d.Id())
		d.Set("users", convertStringSliceToSet(userIDList))
		return nil
	}

	ok, err := checkIfGroupHasUsers(ctx, client, groupId, expected_users)
	if err != nil {
		return diag.Errorf("unable to complete group check for user: %v", err)
	}
	if ok {
		d.Set("group_id", d.Id())
		userIDList, _ := listGroupUserIDs(ctx, m, d.Id())
		d.Set("users", convertStringSliceToSet(userIDList))
		return nil
	} else {
		d.SetId("")
		logger(m).Info("group (%s) did not have expected memberships or did not exist", groupId)
		return nil
	}
}

func resourceGroupMembershipsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	groupId := d.Get("group_id").(string)
	users := convertInterfaceToStringSetNullable(d.Get("users"))
	client := getOktaClientFromMetadata(m)
	err := removeGroupMembers(ctx, client, groupId, users)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceGroupMembershipsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	groupId := d.Get("group_id").(string)
	client := getOktaClientFromMetadata(m)

	oldUsers, newUsers := d.GetChange("users")

	oldSet := oldUsers.(*schema.Set)
	newSet := newUsers.(*schema.Set)

	usersToAdd := convertInterfaceArrToStringArr(newSet.Difference(oldSet).List())
	usersToRemove := convertInterfaceArrToStringArr(oldSet.Difference(newSet).List())

	err := addGroupMembers(ctx, client, groupId, usersToAdd)
	if err != nil {
		diag.FromErr(err)
	}

	err = removeGroupMembers(ctx, client, groupId, usersToRemove)
	if err != nil {
		diag.FromErr(err)
	}

	return nil
}

// checkIfGroupHasUsers firstly checks if the group has users given an immediate
// API call to list group users. It will additionally compare the current list
// of users found in the API call with a slice of known users passed in to the
// function. This second comparison is a stateful comparison.
func checkIfGroupHasUsers(ctx context.Context, client *okta.Client, groupId string, users []string) (bool, error) {
	// TODO: This method should be renamed and/or refactored. Its name implies
	// it is only checking for users but it can return false to signal that
	// users have changed.
	groupUsers, resp, err := client.Group.ListGroupUsers(ctx, groupId, &query.Params{Limit: defaultPaginationLimit})
	if err := suppressErrorOn404(resp, err); err != nil {
		return false, fmt.Errorf("unable to return membership for group (%s) from API", groupId)
	}

	for resp.HasNextPage() {
		var additionalUsers []*okta.User
		resp, err = resp.Next(context.Background(), &additionalUsers)
		if err != nil {
			return false, fmt.Errorf("unable to return membership for group (%s) from API", groupId)
		}
		groupUsers = append(groupUsers, additionalUsers...)
	}

	// We need to return false if there isn't any users present. For eventual
	// consistency issues create has a backoff/retry when we return false
	// without an error here. Read occurs in the future and so it doesn't guard
	// against eventual consistency.
	if len(groupUsers) == 0 {
		return false, nil
	}

	// Create a set to compare the  users slice passed into the check to compare
	// with what the api returns from list group users API call.
	expectedUserSet := make(map[string]bool)

	// We train the set with false values for our previously known users.
	for _, user := range users {
		expectedUserSet[user] = false
	}

	// We confirm the latest user ids from list group users API call are still
	// present in the set.
	for _, user := range groupUsers {
		if _, ok := expectedUserSet[user.Id]; ok {
			expectedUserSet[user.Id] = true
		}
	}

	// If one of the known users in the call to the check function is no longer
	// in the list group users API call we return false.
	for _, state := range expectedUserSet {
		if !state {
			return false, nil
		}
	}

	return true, nil
}
