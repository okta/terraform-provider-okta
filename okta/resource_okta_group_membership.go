package okta

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

func resourceGroupMembership() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupMembershipCreate,
		ReadContext:   resourceGroupMembershipRead,
		UpdateContext: nil,
		DeleteContext: resourceGroupMembershipDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "",
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of a Okta Group",
				ForceNew:    true,
			},
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of a Okta User",
				ForceNew:    true,
			},
		},
	}
}

func resourceGroupMembershipCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	groupId := d.Get("group_id").(string)
	userId := d.Get("user_id").(string)
	logger(m).Info("adding user to group", "group", groupId, "user", userId)
	client := getOktaClientFromMetadata(m)
	_, err := client.Group.AddUserToGroup(ctx, groupId, userId)
	if err != nil {
		return diag.Errorf("failed to add user to group: %v", err)
	}
	d.SetId(fmt.Sprintf("%s+%s", groupId, userId))
	return resourceGroupMembershipRead(ctx, d, m)
}

func resourceGroupMembershipRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	groupId := d.Get("group_id").(string)
	userId := d.Get("user_id").(string)
	logger(m).Info("checking for membership in group", "group", groupId, "user", userId)
	client := getOktaClientFromMetadata(m)
	inGroup, err := checkIfUserInGroup(ctx, client, groupId, userId)
	if err != nil {
		return diag.Errorf("unable to complete group check for user: %v", err)
	}
	if inGroup {
		return nil
	} else {
		d.SetId("")
		logger(m).Info("user is not in group", "group", groupId, "user", userId)
		return nil
	}
}

func resourceGroupMembershipDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	groupId := d.Get("group_id").(string)
	userId := d.Get("user_id").(string)
	logger(m).Info("removing user to group", "group", groupId, "user", userId)
	client := getOktaClientFromMetadata(m)
	_, err := client.Group.RemoveUserFromGroup(ctx, groupId, userId)
	if err != nil {
		return diag.Errorf("failed to remove user to group: %v", err)
	}
	return nil
}

func checkIfUserInGroup(ctx context.Context, client *okta.Client, groupId string, userId string) (bool, error) {
	for {
		users, resp, err := client.Group.ListGroupUsers(ctx, groupId, &query.Params{})
		if err != nil {
			return false, err
		}
		for _, user := range users {
			if userId == user.Id {
				return true, nil
			}
		}
		if resp.HasNextPage() {
			continue
		} else {
			break
		}
	}
	return false, nil
}
