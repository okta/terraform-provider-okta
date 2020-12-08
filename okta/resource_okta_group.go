package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		UpdateContext: resourceGroupUpdate,
		DeleteContext: resourceGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Group name",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Group description",
			},
			"users": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Users associated with the group. This can also be done per user.",
			},
		},
	}
}

func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	logger(m).Info("creating group", "name", d.Get("name").(string))
	group := buildGroup(d)
	responseGroup, _, err := getOktaClientFromMetadata(m).Group.CreateGroup(ctx, *group)
	if err != nil {
		return diag.Errorf("failed to create group: %v", err)
	}
	d.SetId(responseGroup.Id)
	err = updateGroupUsers(ctx, d, m)
	if err != nil {
		return diag.Errorf("failed to update group users: %v", err)
	}
	return resourceGroupRead(ctx, d, m)
}

func resourceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	logger(m).Info("reading group", "id", d.Id(), "name", d.Get("name").(string))
	g, resp, err := getOktaClientFromMetadata(m).Group.GetGroup(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get group: %v", err)
	}
	if g == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("name", g.Profile.Name)
	_ = d.Set("description", g.Profile.Description)
	err = syncGroupUsers(ctx, d, m)
	if err != nil {
		return diag.Errorf("failed to get group users: %v", err)
	}
	return nil
}

func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	logger(m).Info("updating group", "id", d.Id(), "name", d.Get("name").(string))
	group := buildGroup(d)
	_, _, err := getOktaClientFromMetadata(m).Group.UpdateGroup(ctx, d.Id(), *group)
	if err != nil {
		return diag.Errorf("failed to update group: %v", err)
	}
	err = updateGroupUsers(ctx, d, m)
	if err != nil {
		return diag.Errorf("failed to update group users: %v", err)
	}
	return resourceGroupRead(ctx, d, m)
}

func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	logger(m).Info("deleting group", "id", d.Id(), "name", d.Get("name").(string))
	_, err := getOktaClientFromMetadata(m).Group.DeleteGroup(ctx, d.Id())
	if err != nil {
		return diag.Errorf("failed to delete group: %v", err)
	}
	return nil
}

func syncGroupUsers(ctx context.Context, d *schema.ResourceData, m interface{}) error {
	// Only sync when the user opts in by outlining users in the group config
	if _, exists := d.GetOk("users"); !exists {
		return nil
	}
	userIDList, err := listGroupUserIDs(ctx, m, d.Id())
	if err != nil {
		return err
	}
	return d.Set("users", convertStringSetToInterface(userIDList))
}

func updateGroupUsers(ctx context.Context, d *schema.ResourceData, m interface{}) error {
	// Only sync when the user opts in by outlining users in the group config
	// To remove all users, define an empty set
	arr, exists := d.GetOk("users")
	if !exists {
		return nil
	}

	client := getOktaClientFromMetadata(m)
	existingUserList, _, err := client.Group.ListGroupUsers(ctx, d.Id(), nil)
	if err != nil {
		return err
	}

	rawArr := arr.(*schema.Set).List()
	userIDList := make([]string, len(rawArr))

	for i, ifaceID := range rawArr {
		userID := ifaceID.(string)
		userIDList[i] = userID

		if !containsUser(existingUserList, userID) {
			resp, err := client.Group.AddUserToGroup(ctx, d.Id(), userID)
			if err != nil {
				return responseErr(resp, err)
			}
		}
	}

	for _, user := range existingUserList {
		if !contains(userIDList, user.Id) {
			err := suppressErrorOn404(client.Group.RemoveUserFromGroup(ctx, d.Id(), user.Id))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func containsUser(users []*okta.User, id string) bool {
	for _, user := range users {
		if user.Id == id {
			return true
		}
	}
	return false
}

func buildGroup(d *schema.ResourceData) *okta.Group {
	return &okta.Group{
		Profile: &okta.GroupProfile{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
		},
	}
}
