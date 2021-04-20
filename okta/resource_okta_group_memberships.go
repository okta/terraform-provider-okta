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
		Importer:      nil,
		Description:   "Resource to manage a set of group memberships for a specific user.",
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
	err := addGroupMembers(ctx, client, groupId, users)
	if err != nil {
		return diag.FromErr(err)
	}
	bOff := backoff.NewExponentialBackOff()
	bOff.MaxElapsedTime = time.Second * 10
	bOff.InitialInterval = time.Second
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
	users := convertInterfaceToStringSetNullable(d.Get("users"))
	client := getOktaClientFromMetadata(m)
	ok, err := checkIfGroupHasUsers(ctx, client, groupId, users)
	if err != nil {
		return diag.Errorf("unable to complete group check for user: %v", err)
	}
	if ok {
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

	old, new := d.GetChange("users")

	oldSet := old.(*schema.Set)
	newSet := new.(*schema.Set)

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

func checkIfGroupHasUsers(ctx context.Context, client *okta.Client, groupId string, users []string) (bool, error) {
	groupUsers, resp, err := client.Group.ListGroupUsers(ctx, groupId, &query.Params{Limit: defaultPaginationLimit})
	exists, err := doesResourceExist(resp, err)
	if err != nil {
		return false, fmt.Errorf("unable to return membership for group (%s) from API", groupId)
	}

	if !exists {
		return false, nil
	}
	if resp.HasNextPage() {
		for {
			var additionalUsers []*okta.User
			resp, err := resp.Next(ctx, additionalUsers)
			if err != nil {
				return false, fmt.Errorf("unable to return membership for group (%s) from API", groupId)
			}
			groupUsers = append(groupUsers, additionalUsers...)
			if !resp.HasNextPage() {
				break
			}
		}
	}

	// Create set of users
	expectedUserSet := make(map[string]bool)

	for _, user := range users {
		expectedUserSet[user] = false
	}

	// Use users pulled from user and mark set if found
	for _, user := range groupUsers {
		if _, ok := expectedUserSet[user.Id]; ok {
			expectedUserSet[user.Id] = true
		}
	}

	// Check set for any missing values
	for _, state := range expectedUserSet {
		if !state {
			return false, nil
		}
	}

	return true, nil
}
