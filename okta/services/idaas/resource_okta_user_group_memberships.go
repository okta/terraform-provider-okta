package idaas

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	v6okta "github.com/okta/okta-sdk-golang/v6/okta"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func resourceUserGroupMemberships() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserGroupMembershipsCreate,
		ReadContext:   resourceUserGroupMembershipsRead,
		UpdateContext: resourceUserGroupMembershipsUpdate,
		DeleteContext: resourceUserGroupMembershipsDelete,
		Importer:      nil,
		Description: `Resource to manage a set of group memberships for a specific user.
This resource allows you to bulk manage groups for a single user, independent of
the user schema itself. This allows you to manage group membership in terraform
without overriding other automatic membership operations performed by group rules
and other non-managed actions.
**Important**: The default behavior of the resource is to only maintain the
state of group ids that are assigned to it. This behavior will signal drift only if
those groups stop being part of the user's memberships. If the desired behavior is
to track all groups that are added/removed from the user make use of the
'track_all_groups' argument with this resource.`,
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of a Okta User",
				ForceNew:    true,
			},
			"groups": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "The list of Okta group IDs which the user should have membership managed for.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceUserGroupMembershipsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	userId := d.Get("user_id").(string)
	groups := utils.ConvertInterfaceToStringSetNullable(d.Get("groups"))
	client := getOktaV6ClientFromMetadata(meta)
	err := addUserToGroups(ctx, client, userId, groups)
	if err != nil {
		return diag.FromErr(err)
	}
	boc := utils.NewExponentialBackOffWithContext(ctx, 10*time.Second)
	err = backoff.Retry(func() error {
		ok, err := checkIfUserHasGroups(ctx, client, userId, groups)
		if doNotRetry(meta, err) {
			return backoff.Permanent(err)
		}
		if err != nil {
			return backoff.Permanent(err)
		}
		if ok {
			return nil
		}
		return fmt.Errorf("user (%s) did not have expected group memberships after multiple checks", userId)
	}, boc)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(userId)
	return nil
}

func resourceUserGroupMembershipsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	userId := d.Get("user_id").(string)
	groups := utils.ConvertInterfaceToStringSetNullable(d.Get("groups"))
	client := getOktaV6ClientFromMetadata(meta)
	trackAllGroups := d.Get("track_all_groups").(bool)

	if trackAllGroups {
		changed, newGroupIDs, err := checkIfGroupsHaveChanged(ctx, client, userId, &groups)
		if err != nil {
			return diag.Errorf("an error occurred checking group ids for user %q, error: %+v", userId, err)
		}
		if changed {
			d.Set("groups", utils.ConvertStringSliceToSet(*newGroupIDs))
		}
		return nil
	}

	// Legacy behavior: if any managed groups were removed externally, shrink the
	// state to the surviving subset so Terraform plans to add them back rather
	// than recreating the entire resource.
	changed, newGroupIDs, err := checkIfGroupsHaveBeenRemoved(ctx, client, userId, &groups)
	if err != nil {
		return diag.Errorf("an error occurred checking group ids for user %q, error: %+v", userId, err)
	}
	if changed {
		d.Set("groups", utils.ConvertStringSliceToSet(*newGroupIDs))
	}
	return nil
}

func resourceUserGroupMembershipsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	userId := d.Get("user_id").(string)
	groups := utils.ConvertInterfaceToStringSetNullable(d.Get("groups"))
	client := getOktaV6ClientFromMetadata(meta)
	err := removeUserFromGroups(ctx, client, userId, groups)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceUserGroupMembershipsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	userId := d.Get("user_id").(string)
	client := getOktaV6ClientFromMetadata(meta)

	old, new := d.GetChange("groups")
	oldSet := old.(*schema.Set)
	newSet := new.(*schema.Set)

	groupsToAdd := utils.ConvertInterfaceArrToStringArr(newSet.Difference(oldSet).List())
	groupsToRemove := utils.ConvertInterfaceArrToStringArr(oldSet.Difference(newSet).List())

	err := addUserToGroups(ctx, client, userId, groupsToAdd)
	if err != nil {
		return diag.FromErr(err)
	}
	err = removeUserFromGroups(ctx, client, userId, groupsToRemove)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

// checkIfGroupsHaveBeenRemoved returns true when any of the tracked groups are
// no longer in the user's actual memberships. The returned slice contains only
// the groups that are still present, so the caller can narrow state and let
// Terraform plan to re-add the missing ones.
func checkIfGroupsHaveBeenRemoved(ctx context.Context, client *v6okta.APIClient, userId string, groups *[]string) (bool, *[]string, error) {
	noop := []string{}
	if groups == nil || len(*groups) == 0 {
		return false, &noop, nil
	}

	// Ledger of managed groups. When all have been found in the API response
	// we can stop paginating early.
	remaining := toStrIndexedMap(groups)

	userGroups, resp, err := client.UserResourcesAPI.ListUserGroups(ctx, userId).Execute()
	if err := utils.SuppressErrorOn404_V6(resp, err); err != nil {
		return false, &noop, fmt.Errorf("unable to list groups for user (%s) from API, error: %+v", userId, err)
	}

	for _, group := range userGroups {
		delete(*remaining, group.GetId())
		if len(*remaining) == 0 {
			return false, &noop, nil
		}
	}

	for resp.HasNextPage() {
		userGroups = nil
		resp, err = resp.Next(&userGroups)
		if err != nil {
			return false, &noop, fmt.Errorf("unable to list groups for user (%s) from API, error: %+v", userId, err)
		}
		for _, group := range userGroups {
			delete(*remaining, group.GetId())
			if len(*remaining) == 0 {
				return false, &noop, nil
			}
		}
	}

	// Any IDs still in remaining were not found — they have been removed.
	// Build the surviving subset in the original order.
	newGroups := []string{}
	for _, groupId := range *groups {
		if _, missing := (*remaining)[groupId]; !missing {
			newGroups = append(newGroups, groupId)
		}
	}

	return true, &newGroups, nil
}

// checkIfGroupsHaveChanged returns true when the user's actual group memberships
// differ from the tracked set in state. When changed, it also returns the full
// current list so the caller can update the state attribute.
func checkIfGroupsHaveChanged(ctx context.Context, client *v6okta.APIClient, userId string, groups *[]string) (bool, *[]string, error) {
	noop := []string{}
	if groups == nil {
		return false, &noop, nil
	}

	oldGroups := toStrIndexedMap(groups)
	changed := false
	groupsFromAPI := []string{}

	userGroups, resp, err := client.UserResourcesAPI.ListUserGroups(ctx, userId).Execute()
	if err := utils.SuppressErrorOn404_V6(resp, err); err != nil {
		return false, &noop, fmt.Errorf("unable to list groups for user (%s) from API, error: %+v", userId, err)
	}

	for _, group := range userGroups {
		if group.GetType() == "BUILT_IN" {
			continue
		}
		if _, found := (*oldGroups)[group.GetId()]; !found {
			changed = true
		}
		groupsFromAPI = append(groupsFromAPI, group.GetId())
	}

	for resp.HasNextPage() {
		userGroups = nil
		resp, err = resp.Next(&userGroups)
		if err != nil {
			return false, &noop, fmt.Errorf("unable to list groups for user (%s) from API, error: %+v", userId, err)
		}
		for _, group := range userGroups {
			if group.GetType() == "BUILT_IN" {
				continue
			}
			if _, found := (*oldGroups)[group.GetId()]; !found {
				changed = true
			}
			groupsFromAPI = append(groupsFromAPI, group.GetId())
		}
	}

	if len(*oldGroups) != len(groupsFromAPI) {
		changed = true
	}

	var result *[]string = &noop
	if changed {
		result = &groupsFromAPI
	}

	return changed, result, nil
}

func checkIfUserHasGroups(ctx context.Context, client *v6okta.APIClient, userId string, groups []string) (bool, error) {
	userGroups, resp, err := client.UserResourcesAPI.ListUserGroups(ctx, userId).Execute()
	if err := utils.SuppressErrorOn404_V6(resp, err); err != nil {
		return false, fmt.Errorf("unable to return groups for user (%s) from API", userId)
	}
	var nextUserGroups []v6okta.Group

	for resp.HasNextPage() {
		resp, err = resp.Next(&nextUserGroups)

		if err := utils.SuppressErrorOn404_V6(resp, err); err != nil {
			return false, fmt.Errorf("unable to get next page of groups for user (%s) from API", userId)
		}

		userGroups = append(userGroups, nextUserGroups...)
	}

	if len(userGroups) == 0 {
		return false, nil
	}

	// Create set of groups
	expectedGroupSet := make(map[string]bool)

	for _, group := range groups {
		expectedGroupSet[group] = false
	}

	// Use groups pulled from user and mark set if found
	for _, group := range userGroups {
		if _, ok := expectedGroupSet[group.GetId()]; ok {
			expectedGroupSet[group.GetId()] = true
		}
	}

	// Check set for any missing values
	for _, state := range expectedGroupSet {
		if !state {
			return false, nil
		}
	}

	return true, nil
}
