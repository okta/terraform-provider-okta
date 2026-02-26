package idaas

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
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
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				importID := strings.Split(d.Id(), "/")
				if len(importID) > 2 {
					return nil, errors.New("invalid format used for import ID, format must be 'group_id' or 'group_id/true'")
				}
				if len(importID) == 2 {
					d.Set("track_all_users", importID[1] == "true")
				}
				d.SetId(importID[0])
				d.Set("group_id", importID[0])
				return []*schema.ResourceData{d}, nil
			},
		},
		Description: `Resource to manage a set of memberships for a specific group.
This resource will allow you to bulk manage group membership in Okta for a given
group. This offers an interface to pass multiple users into a single resource
call, for better API resource usage. If you need a relationship of a single 
user to many groups, please use the 'okta_user_group_memberships' resource.
**Important**: The default behavior of the resource is to only maintain the
state of user ids that are assigned it. This behavior will signal drift only if
those users stop being part of the group. If the desired behavior is track all
users that are added/removed from the group make use of the 'track_all_users'
argument with this resource.`,
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

func resourceGroupMembershipsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	groupId := d.Get("group_id").(string)
	users := utils.ConvertInterfaceToStringSetNullable(d.Get("users"))
	// if read is being called via import "id" will not be blank and "group_id"
	// will be blank, so set group_id accordingly
	if d.Id() != "" && groupId == "" {
		groupId = d.Id()
		d.Set("group_id", groupId)
	}

	client := getOktaClientFromMetadata(meta)

	if len(users) == 0 {
		d.SetId(groupId)
		return nil
	}
	err := addGroupMembers(ctx, client, groupId, users)
	if err != nil {
		return diag.FromErr(err)
	}
	boc := utils.NewExponentialBackOffWithContext(ctx, 10*time.Second)
	// During create the Okta service can have eventual consistency issues when
	// adding users to a group. Use a backoff to wait for at list one user to be
	// associated with the group.
	err = backoff.Retry(func() error {
		// TODO, should we wait for all users to be added to the group?
		ok, err := checkIfGroupHasUsers(ctx, client, groupId, users)
		if doNotRetry(meta, err) {
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

func resourceGroupMembershipsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(meta)
	groupId := d.Get("group_id").(string)
	oldUsers := utils.ConvertInterfaceToStringSetNullable(d.Get("users"))
	trackAllUsers := d.Get("track_all_users").(bool)

	// New behavior, tracking all users.
	if trackAllUsers {
		changed, newUserIDs, err := checkIfUsersHaveChanged(ctx, client, groupId, &oldUsers)
		if err != nil {
			return diag.Errorf("An error occurred checking user ids for group %q, error: %+v", groupId, err)
		}
		if changed {
			// set the new user ids if users have changed
			d.Set("users", utils.ConvertStringSliceToSet(*newUserIDs))
		}

		return nil
	}

	// Legacy behavior is just to check if any users have left the group.
	changed, newUserIDs, err := checkIfUsersHaveBeenRemoved(ctx, client, groupId, &oldUsers)
	if err != nil {
		return diag.Errorf("An error occurred checking user ids for group %q, error: %+v", groupId, err)
	}
	if changed {
		// The user list has changed, set the new user ids to the users value.
		d.Set("users", utils.ConvertStringSliceToSet(*newUserIDs))
	}
	return nil
}

func resourceGroupMembershipsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	groupId := d.Get("group_id").(string)
	users := utils.ConvertInterfaceToStringSetNullable(d.Get("users"))
	client := getOktaClientFromMetadata(meta)
	err := removeGroupMembers(ctx, client, groupId, users)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceGroupMembershipsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	groupId := d.Get("group_id").(string)
	client := getOktaClientFromMetadata(meta)

	oldUsers, newUsers := d.GetChange("users")

	oldSet := oldUsers.(*schema.Set)
	newSet := newUsers.(*schema.Set)

	usersToAdd := utils.ConvertInterfaceArrToStringArr(newSet.Difference(oldSet).List())
	usersToRemove := utils.ConvertInterfaceArrToStringArr(oldSet.Difference(newSet).List())

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
	// users slice can be sized 0 if this is a read from import
	if users == nil {
		return false, &noop, nil
	}

	// We are using the old users map as a ledger to find users that have been
	// removed from the user list.
	oldUsers := toStrIndexedMap(users)
	changed := false
	// Collect all user ids that are returned from the API
	usersFromAPI := []string{}

	groupUsers, resp, err := client.Group.ListGroupUsers(ctx, groupId, &query.Params{Limit: utils.DefaultPaginationLimit})
	if err := utils.SuppressErrorOn404(resp, err); err != nil {
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

	groupUsers, resp, err := client.Group.ListGroupUsers(ctx, groupId, &query.Params{Limit: utils.DefaultPaginationLimit})
	if err := utils.SuppressErrorOn404(resp, err); err != nil {
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

func checkIfGroupHasUsers(ctx context.Context, client *sdk.Client, groupId string, _ []string) (bool, error) {
	groupUsers, resp, err := client.Group.ListGroupUsers(ctx, groupId, &query.Params{Limit: utils.DefaultPaginationLimit})
	if err := utils.SuppressErrorOn404(resp, err); err != nil {
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
