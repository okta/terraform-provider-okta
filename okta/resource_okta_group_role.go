package okta

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func resourceGroupRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupRoleCreate,
		ReadContext:   resourceGroupRoleRead,
		UpdateContext: resourceGroupRoleUpdate,
		DeleteContext: resourceGroupRoleDelete,
		Importer:      &schema.ResourceImporter{StateContext: resourceGroupRoleImporter},
		CustomizeDiff: customdiff.ForceNewIf("target_group_list", func(_ context.Context, d *schema.ResourceDiff, m interface{}) bool {
			if d.HasChange("target_group_list") {
				if len(convertInterfaceToStringSet(d.Get("target_group_list"))) == 0 {
					return true
				}
			}
			return false
		}),
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of group to attach admin roles to",
				ForceNew:    true,
			},
			"role_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of Role to assign",
				ForceNew:    true,
				ValidateDiagFunc: stringInSlice([]string{
					"API_ACCESS_MANAGEMENT_ADMIN",
					"APP_ADMIN",
					"GROUP_MEMBERSHIP_ADMIN",
					"HELP_DESK_ADMIN",
					"MOBILE_ADMIN",
					"ORG_ADMIN",
					"READ_ONLY_ADMIN",
					"REPORT_ADMIN",
					"SUPER_ADMIN",
					"USER_ADMIN",
				}),
			},
			"target_group_list": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of groups ids for the targets of the admin role.",
				Optional:    true,
			},
		},
	}
}

func resourceGroupRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	groupID := d.Get("group_id").(string)
	roleType := d.Get("role_type").(string)
	client := getOktaClientFromMetadata(m)
	logger(m).Info("assigning role to group", "group_id", groupID, "role_type", roleType)
	role, _, err := client.Group.AssignRoleToGroup(ctx, groupID, okta.AssignRoleRequest{
		Type: roleType,
	}, nil)
	if err != nil {
		return diag.Errorf("failed to assign role %s to group %s: %v", roleType, groupID, err)
	}

	groupTargets := convertInterfaceToStringSet(d.Get("target_group_list"))
	if len(groupTargets) > 0 && supportsGroupTargets(roleType) {
		logger(m).Info("scoping admin role assignment to list of groups", "group_id", groupID, "role_id", role.Id, "target_group_list", groupTargets)
		err = addGroupTargetsToRole(ctx, client, groupID, role.Id, groupTargets)
		if err != nil {
			return diag.Errorf("unable to add group target to role assignment %s for group %s: %v", role.Id, groupID, err)
		}
	}

	d.SetId(role.Id)
	return resourceGroupRoleRead(ctx, d, m)
}

func resourceGroupRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	groupID := d.Get("group_id").(string)
	client := getOktaClientFromMetadata(m)
	rolesAssigned, _, err := client.Group.ListGroupAssignedRoles(ctx, groupID, nil)
	if err != nil {
		return diag.Errorf("failed to list roles assigned to group %s: %v", groupID, err)
	}
	for _, role := range rolesAssigned {
		if role.Id == d.Id() {

			if supportsGroupTargets(role.Type) {
				currentTargets, _, err := client.Group.ListGroupTargetsForGroupRole(ctx, groupID, role.Id, nil)
				if err != nil {
					return diag.Errorf("unable to get admin assignment %s for group %s: %v", role.Id, groupID, err)
				}
				groupIDs := getGroupIds(currentTargets)
				_ = d.Set("target_group_list", groupIDs)
			}

			_ = d.Set("role_type", role.Type)
			return nil
		}
	}
	logger(m).Info("no roles found assigned to group", "group_id", groupID)
	d.SetId("")
	return nil
}

func resourceGroupRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	groupID := d.Get("group_id").(string)
	roleID := d.Id()
	roleType := d.Get("role_type").(string)
	client := getOktaClientFromMetadata(m)
	if d.HasChange("target_group_list") && supportsGroupTargets(roleType) {
		currentTargets, _, err := client.Group.ListGroupTargetsForGroupRole(ctx, groupID, roleID, nil)
		if err != nil {
			return diag.Errorf("unable to get admin assignment %s for group %s: %v", roleID, groupID, err)
		}
		currentTargetIds := getGroupIds(currentTargets)
		err = removeGroupTargetsFromRole(ctx, client, groupID, roleID, currentTargetIds)
		if err != nil {
			return diag.Errorf("unable to remove group target from admin role assignment %s of group %s: %v", roleID, groupID, err)
		}
		newTargetIds := convertInterfaceToStringSet(d.Get("target_group_list"))
		err = addGroupTargetsToRole(ctx, client, groupID, roleID, newTargetIds)
		if err != nil {
			return diag.Errorf("unable to add group target to role assignment %s for group %s: %v", roleID, groupID, err)
		}
		_ = d.Set("target_group_list", newTargetIds)
	}
	return resourceGroupRoleRead(ctx, d, m)
}

func resourceGroupRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	groupID := d.Get("group_id").(string)
	roleType := d.Get("role_type").(string)
	logger(m).Info("deleting assigned role from group", "group_id", groupID, "role_type", roleType)
	_, err := getOktaClientFromMetadata(m).Group.RemoveRoleFromGroup(ctx, groupID, d.Id())
	if err != nil {
		return diag.Errorf("failed to remove role %s assigned to group %s: %v", roleType, groupID, err)
	}
	return nil
}

func resourceGroupRoleImporter(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	importID := strings.Split(d.Id(), "/")
	if len(importID) != 2 {
		err := fmt.Errorf("invalid format used for import ID, format must be group_id/role_assignment_id")
		return nil, err
	}
	groupID := importID[0]
	roleID := importID[1]
	client := getOktaClientFromMetadata(m)
	rolesAssigned, _, err := client.Group.ListGroupAssignedRoles(ctx, groupID, nil)
	if err != nil {
		return nil, err
	}
	for _, role := range rolesAssigned {
		if role.Id == roleID {
			d.SetId(roleID)
			_ = d.Set("group_id", groupID)
			_ = d.Set("role_type", role.Type)

			if supportsGroupTargets(role.Type) {
				currentTargets, _, err := client.Group.ListGroupTargetsForGroupRole(ctx, groupID, role.Id, nil)
				if err != nil {
					return nil, fmt.Errorf("unable to get admin assignment %s for group %s: %v", role.Id, groupID, err)
				}
				groupIDs := getGroupIds(currentTargets)
				_ = d.Set("target_group_list", groupIDs)
			}
			return []*schema.ResourceData{d}, nil
		}
	}
	err = fmt.Errorf("unable to find the role ID %s assigned to the group %s", roleID, groupID)
	return nil, err
}

// Group Role Target Helpers

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

func supportsGroupTargets(roleType string) bool {
	supportedRoles := []string{"GROUP_MEMBERSHIP_ADMIN", "HELP_DESK_ADMIN", "USER_ADMIN"}
	for _, role := range supportedRoles {
		if roleType == role {
			return true
		}
	}
	return false
}
