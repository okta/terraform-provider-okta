package okta

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func resourceUserGroupMemberships() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserGroupMembershipsCreate,
		ReadContext:   resourceUserGroupMembershipsRead,
		UpdateContext: resourceUserGroupMembershipsUpdate,
		DeleteContext: resourceUserGroupMembershipsDelete,
		Importer:      nil,
		Description:   "Resource to manage a set of group memberships for a specific user.",
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

func resourceUserGroupMembershipsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	userId := d.Get("user_id").(string)
	groups := convertInterfaceToStringSetNullable(d.Get("groups"))
	client := getOktaClientFromMetadata(m)
	err := addUserToGroups(ctx, client, userId, groups)
	if err != nil {
		return diag.FromErr(err)
	}
	bOff := backoff.NewExponentialBackOff()
	bOff.MaxElapsedTime = time.Second * 10
	bOff.InitialInterval = time.Second
	err = backoff.Retry(func() error {
		ok, err := checkIfUserHasGroups(ctx, client, userId, groups)
		if err != nil {
			return backoff.Permanent(err)
		}
		if ok {
			return nil
		}
		return fmt.Errorf("user (%s) did not have expected group memberships after multiple checks", userId)
	}, bOff)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(userId)
	return nil
}

func resourceUserGroupMembershipsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	userId := d.Get("user_id").(string)
	groups := convertInterfaceToStringSetNullable(d.Get("groups"))
	client := getOktaClientFromMetadata(m)
	ok, err := checkIfUserHasGroups(ctx, client, userId, groups)
	if err != nil {
		return diag.Errorf("unable to complete group check for user: %v", err)
	}
	if ok {
		return nil
	} else {
		d.SetId("")
		logger(m).Info("user (%s) did not have expected group memberships or did not exist", userId)
		return nil
	}
}

func resourceUserGroupMembershipsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	userId := d.Get("user_id").(string)
	groups := convertInterfaceToStringSetNullable(d.Get("groups"))
	client := getOktaClientFromMetadata(m)
	err := removeUserFromGroups(ctx, client, userId, groups)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceUserGroupMembershipsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	userId := d.Get("user_id").(string)
	client := getOktaClientFromMetadata(m)

	old, new := d.GetChange("groups")

	oldSet := old.(*schema.Set)
	newSet := new.(*schema.Set)

	groupsToAdd := convertInterfaceArrToStringArr(newSet.Difference(oldSet).List())
	groupsToRemove := convertInterfaceArrToStringArr(oldSet.Difference(newSet).List())

	err := addUserToGroups(ctx, client, userId, groupsToAdd)
	if err != nil {
		diag.FromErr(err)
	}

	err = removeUserFromGroups(ctx, client, userId, groupsToRemove)
	if err != nil {
		diag.FromErr(err)
	}

	return nil
}

func checkIfUserHasGroups(ctx context.Context, client *okta.Client, userId string, groups []string) (bool, error) {
	userGroups, resp, err := client.User.ListUserGroups(ctx, userId)
	exists, err := doesResourceExist(resp, err)
	if err != nil {
		return false, fmt.Errorf("unable to return groups for user (%s) from API", userId)
	}

	if !exists {
		return false, nil
	}

	// Create set of groups
	expectedGroupSet := make(map[string]bool)

	for _, group := range groups {
		expectedGroupSet[group] = false
	}

	// Use groups pulled from user and mark set if found
	for _, group := range userGroups {
		if _, ok := expectedGroupSet[group.Id]; ok {
			expectedGroupSet[group.Id] = true
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

func addUserToGroups(ctx context.Context, client *okta.Client, userId string, groups []string) error {
	for _, group := range groups {
		resp, err := client.Group.AddUserToGroup(ctx, group, userId)
		exists, err := doesResourceExist(resp, err)
		if err != nil {
			return fmt.Errorf("failed to add user (%s) to group (%s): %v", userId, group, err)
		}
		if !exists {
			return fmt.Errorf("targeted object does not exist: %s", err)
		}
	}
	return nil
}

func removeUserFromGroups(ctx context.Context, client *okta.Client, userId string, groups []string) error {
	for _, group := range groups {
		resp, err := client.Group.RemoveUserFromGroup(ctx, group, userId)
		err = suppressErrorOn404(resp, err)
		if err != nil {
			return fmt.Errorf("failed to remove user (%s) from group (%s): %v", userId, group, err)
		}
	}
	return nil
}
