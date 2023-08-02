package okta

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/okta/terraform-provider-okta/sdk/query"
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
			"track_all_users": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "The resource concerns itself with all users added/deleted to the group; even those managed outside of the resource.",
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
	boc := newExponentialBackOffWithContext(ctx, 10*time.Second)
	// During create the Okta service can have eventual consistency issues when
	// adding users to a group. Use a backoff to wait for at list one user to be
	// associated with the group.
	err = backoff.Retry(func() error {
		// TODO, should we wait for all users to be added to the group?
		ok, err := checkIfGroupHasUsers(ctx, client, groupId, users)
		if doNotRetry(m, err) {
			return backoff.Permanent(err)
		}
		if err != nil {
			return backoff.Permanent(err)
		}
		if ok {
			return nil
		}
		return fmt.Errorf("group (%s) did not have expected user memberships after multiple checks", groupId)
	}, boc)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(groupId)
	return nil
}

func resourceGroupMembershipsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(m)
	groupId := d.Get("group_id").(string)
	oldUsers := convertInterfaceToStringSetNullable(d.Get("users"))
	trackAllUsers := d.Get("track_all_users").(bool)

	// New behavior, tracking all users.
	if trackAllUsers {
		changed, newUserIDs, err := checkIfUsersHaveChanged(ctx, client, groupId, &oldUsers)
		if err != nil {
			return diag.Errorf("An error occured checking user ids for group %q, error: %+v", groupId, err)
		}
		if changed {
			// set the new user ids if users have changed
			d.Set("users", convertStringSliceToSet(*newUserIDs))
		}

		return nil
	}

	// Legacy behavior is just to check if any users have left the group.
	changed, newUserIDs, err := checkIfUsersHaveBeenRemoved(ctx, client, groupId, &oldUsers)
	if err != nil {
		return diag.Errorf("An error occured checking user ids for group %q, error: %+v", groupId, err)
	}
	if changed {
		// The user list has changed, set the new user ids to the users value.
		d.Set("users", convertStringSliceToSet(*newUserIDs))
	}
	return nil
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

// checkIfUsersHaveChanged If the function returns true then users have been
// changed and the returned user ids should be considered the new set of users.
// Returns error for API errors. Returns false if no users have changed and the
// slice of returned strings will be empty.
func checkIfUsersHaveChanged(ctx context.Context, client *sdk.Client, groupId string, users *[]string) (bool, *[]string, error) {
	noop := []string{}
	if users == nil || len(*users) == 0 {
		return false, &noop, nil
	}

	// We are using the old users map as a ledger to find users that have been
	// removed from the user list.
	oldUsers := toStrIndexedMap(users)
	changed := false
	// Collect all user ids that are returned from the API
	usersFromAPI := []string{}

	groupUsers, resp, err := client.Group.ListGroupUsers(ctx, groupId, &query.Params{Limit: defaultPaginationLimit})
	if err := suppressErrorOn404(resp, err); err != nil {
		return false, &noop, fmt.Errorf("unable to list users for group (%s) from API, error: %+v", groupId, err)
	}

	for _, user := range groupUsers {
		// if the new user id is not in the old users map then the list of users has changed
		if _, found := (*oldUsers)[user.Id]; !found {
			changed = true
		}
		usersFromAPI = append(usersFromAPI, user.Id)
	}

	for resp.HasNextPage() {
		groupUsers = nil
		resp, err = resp.Next(context.Background(), &groupUsers)
		if err != nil {
			return false, &noop, fmt.Errorf("unable to list users for group (%s) from API, error: %+v", groupId, err)
		}

		for _, user := range groupUsers {
			// if the new user id is not in the old users map then the list of users has changed
			if _, found := (*oldUsers)[user.Id]; !found {
				changed = true
			}
			usersFromAPI = append(usersFromAPI, user.Id)
		}
	}
	if len(*oldUsers) != len(usersFromAPI) {
		changed = true
	}

	var result *[]string = &noop
	if changed {
		result = &usersFromAPI
	}

	return changed, result, nil
}

// checkIfUsersHaveBeenRemoved If the function returns true then users have been
// removed and the subset of returned user ids should be considered the new set
// of users. Returns error for API errors. Returns false if no users have been
// removed and the slice of returned strings will be empty.
func checkIfUsersHaveBeenRemoved(ctx context.Context, client *sdk.Client, groupId string, users *[]string) (bool, *[]string, error) {
	noop := []string{}
	if users == nil || len(*users) == 0 {
		return false, &noop, nil
	}

	// We are using the old users map as a ledger to find users that have been
	// removed from the user list. If it ever becomes sized 0 then we've found
	// all of our user ids and no longer have to make API calls.
	oldUsers := toStrIndexedMap(users)

	groupUsers, resp, err := client.Group.ListGroupUsers(ctx, groupId, &query.Params{Limit: defaultPaginationLimit})
	if err := suppressErrorOn404(resp, err); err != nil {
		return false, &noop, fmt.Errorf("unable to list users for group (%s) from API, error: %+v", groupId, err)
	}

	for _, user := range groupUsers {
		// Deleting user from API from the old users map
		delete(*oldUsers, user.Id)
		if len(*oldUsers) == 0 {
			// All old users have been accounted for.
			return false, &noop, nil
		}
	}

	for resp.HasNextPage() {
		groupUsers = nil
		resp, err = resp.Next(context.Background(), &groupUsers)
		if err != nil {
			return false, &noop, fmt.Errorf("unable to list users for group (%s) from API, error: %+v", groupId, err)
		}
		for _, user := range groupUsers {
			// Deleting user from API from the old users map
			delete(*oldUsers, user.Id)
			if len(*oldUsers) == 0 {
				// All old users have been accounted for.
				return false, &noop, nil
			}
		}
	}

	// Any old users left are the IDs that have been removed from group.

	// This loop keeps the returned new user list in the same order as the users
	// list passed in minus the now missing user IDs.
	newUsers := []string{}
	for _, userId := range *users {
		// Any single user id found in old users should not be appended to the
		// new users result.
		if _, found := (*oldUsers)[userId]; !found {
			newUsers = append(newUsers, userId)
		}
	}

	return true, &newUsers, nil
}

func checkIfGroupHasUsers(ctx context.Context, client *sdk.Client, groupId string, users []string) (bool, error) {
	groupUsers, resp, err := client.Group.ListGroupUsers(ctx, groupId, &query.Params{Limit: defaultPaginationLimit})
	if err := suppressErrorOn404(resp, err); err != nil {
		return false, fmt.Errorf("unable to return membership for group (%s) from API", groupId)
	}
	return (len(groupUsers) > 0), nil
}

func toStrIndexedMap(strs *[]string) *map[string]int {
	result := map[string]int{}
	if strs == nil {
		return &result
	}

	length := len(*strs)
	for i := 0; i < length; i++ {
		result[(*strs)[i]] = i
	}

	return &result
}
