package okta

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func resourceGroupRoleGroupTargets() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupRoleGroupTargetsCreate,
		ReadContext:   resourceGroupRoleGroupTargetsRead,
		UpdateContext: resourceGroupRoleGroupTargetsUpdate,
		DeleteContext: resourceGroupRoleGroupTargetsDelete,
		Importer:      nil,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Group ID for the group with the admin role",
				ForceNew:    true,
			},
			"role_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Role ID of the role assignment of the group",
				ForceNew:    true,
			},
			"group_target_list": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: "List of groups ids for the targets of the admin role",
			},
		},
	}
}

func resourceGroupRoleGroupTargetsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	groupID := d.Get("group_id").(string)
	roleID := d.Get("role_id").(string)
	groupTargets := convertInterfaceToStringSet(d.Get("group_target_list"))
	logger(m).Info("scoping admin role assignment to list of groups", "group_id", groupID, "role_id", roleID, "group_target_list", groupTargets)
	client := getOktaClientFromMetadata(m)
	currentTargets, _, err := client.Group.ListGroupTargetsForGroupRole(ctx, groupID, roleID, nil)
	if err != nil {
		return diag.Errorf("unable to get admin assignment %s for group %s: %v", roleID, groupID, err)
	}
	if len(currentTargets) > 0 {
		return diag.Errorf("group targets already found attached to role, you should not use this resource unless it is the sole manager of target groups")
	}
	err = addGroupTargetsToRole(ctx, client, groupID, roleID, groupTargets)
	if err != nil {
		return diag.Errorf("unable to add group target to role assignment %s for group %s: %v", roleID, groupID, err)
	}
	d.SetId(fmt.Sprintf("%s.%s.grouptargets", groupID, roleID))
	return resourceGroupRoleGroupTargetsRead(ctx, d, m)
}

func resourceGroupRoleGroupTargetsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	groupID := d.Get("group_id").(string)
	roleID := d.Get("role_id").(string)
	logger(m).Info("reading group targets for admin role assignment", "group_id", groupID, "role_id", roleID)
	client := getOktaClientFromMetadata(m)
	currentTargets, _, err := client.Group.ListGroupTargetsForGroupRole(ctx, groupID, roleID, nil)
	if err != nil {
		return diag.Errorf("unable to get admin assignment %s for group %s: %v", roleID, groupID, err)
	}
	if len(currentTargets) == 0 {
		logger(m).Info("no group targets found assigned to group admin assignment", "group_id", groupID, "role_id", roleID)
		d.SetId("")
		return nil
	}
	groupTargets := getGroupIds(currentTargets)
	_ = d.Set("group_target_list", groupTargets)
	return nil
}

func resourceGroupRoleGroupTargetsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	groupID := d.Get("group_id").(string)
	roleID := d.Get("role_id").(string)
	logger(m).Info("updating group targets for admin role assignment", "group_id", groupID, "role_id", roleID)
	client := getOktaClientFromMetadata(m)
	if d.HasChange("group_target_list") {
		currentTargets, _, err := client.Group.ListGroupTargetsForGroupRole(ctx, groupID, roleID, nil)
		if err != nil {
			return diag.Errorf("unable to get admin assignment %s for group %s: %v", roleID, groupID, err)
		}
		currentTargetIds := getGroupIds(currentTargets)
		err = removeGroupTargetsFromRole(ctx, client, groupID, roleID, currentTargetIds)
		if err != nil {
			return diag.Errorf("unable to remove group target from admin role assignment %s of group %s: %v", roleID, groupID, err)
		}
		newTargetIds := convertInterfaceToStringSet(d.Get("group_target_list"))
		err = addGroupTargetsToRole(ctx, client, groupID, roleID, newTargetIds)
		if err != nil {
			return diag.Errorf("unable to add group target to role assignment %s for group %s: %v", roleID, groupID, err)
		}
		_ = d.Set("group_target_list", newTargetIds)
	}
	return resourceGroupRoleGroupTargetsRead(ctx, d, m)
}

func resourceGroupRoleGroupTargetsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	groupID := d.Get("group_id").(string)
	roleID := d.Get("role_id").(string)
	logger(m).Info("removing group targets for admin role assignment", "group_id", groupID, "role_id", roleID)
	client := getOktaClientFromMetadata(m)
	targetIds := convertInterfaceToStringSet(d.Get("group_target_list"))
	err := removeGroupTargetsFromRole(ctx, client, groupID, roleID, targetIds)
	if err != nil {
		return diag.Errorf("unable to add group target to role assignment %s for group %s: %v", roleID, groupID, err)
	}
	return nil
}

func getGroupIds(groups []*okta.Group) []string {
	var groupIds []string
	for _, group := range groups {
		groupIds = append(groupIds, group.Id)
	}
	return groupIds
}

func addGroupTargetsToRole(ctx context.Context, client *okta.Client, groupID string, roleID string, groupTargets []string) error {
	for _, target := range groupTargets {
		_, err := client.Group.AddGroupTargetToGroupAdministratorRoleForGroup(ctx, groupID, roleID, target)
		if err != nil {
			return err
		}
	}
	return nil
}

func removeGroupTargetsFromRole(ctx context.Context, client *okta.Client, groupID string, roleID string, groupTargets []string) error {
	for _, target := range groupTargets {
		_, err := client.Group.RemoveGroupTargetFromGroupAdministratorRoleGivenToGroup(ctx, groupID, roleID, target)
		if err != nil {
			return err
		}
	}
	return nil
}
