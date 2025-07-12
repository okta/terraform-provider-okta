package idaas

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourceUserGroupMemberships() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserGroupMembershipsCreate,
		ReadContext:   resourceUserGroupMembershipsRead,
		UpdateContext: resourceUserGroupMembershipsUpdate,
		DeleteContext: resourceUserGroupMembershipsDelete,
		Importer:      nil,
		Description:   "Resource to manage a set of group memberships for a specific user. This resource allows you to bulk manage groups for a single user, independent of the user schema itself. This allows you to manage group membership in terraform without overriding other automatic membership operations performed by group rules and other non-managed actions.",
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
	client := getOktaClientFromMetadata(meta)
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
	client := getOktaClientFromMetadata(meta)
	ok, err := checkIfUserHasGroups(ctx, client, userId, groups)
	if err != nil {
		return diag.Errorf("unable to complete group check for user: %v", err)
	}
	if ok {
		return nil
	} else {
		d.SetId("")
		logger(meta).Info("user (%s) did not have expected group memberships or did not exist", userId)
		return nil
	}
}

func resourceUserGroupMembershipsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	userId := d.Get("user_id").(string)
	groups := utils.ConvertInterfaceToStringSetNullable(d.Get("groups"))
	client := getOktaClientFromMetadata(meta)
	err := removeUserFromGroups(ctx, client, userId, groups)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceUserGroupMembershipsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	userId := d.Get("user_id").(string)
	client := getOktaClientFromMetadata(meta)

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

func checkIfUserHasGroups(ctx context.Context, client *sdk.Client, userId string, groups []string) (bool, error) {
	userGroups, resp, err := client.User.ListUserGroups(ctx, userId)
	if err := utils.SuppressErrorOn404(resp, err); err != nil {
		return false, fmt.Errorf("unable to return groups for user (%s) from API", userId)
	}
	var nextUserGroups []*sdk.Group

	for resp.HasNextPage() {
		resp, err = resp.Next(ctx, &nextUserGroups)

		if err := utils.SuppressErrorOn404(resp, err); err != nil {
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
