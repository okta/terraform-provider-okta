package okta

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"strings"
)

func resourceGroupRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupRoleCreate,
		ReadContext:   resourceGroupRoleRead,
		UpdateContext: nil,
		DeleteContext: resourceGroupRoleDelete,
		Importer:      &schema.ResourceImporter{StateContext: resourceGroupRoleImporter},
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
		},
	}
}

func resourceGroupRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	groupID := d.Get("group_id").(string)
	roleType := d.Get("role_type").(string)
	logger(m).Info("assigning role to group", "group_id", groupID, "role_type", roleType)
	role, _, err := getOktaClientFromMetadata(m).Group.AssignRoleToGroup(ctx, groupID, okta.AssignRoleRequest{
		Type: roleType,
	}, nil)
	if err != nil {
		return diag.Errorf("failed to assign role %s to group %s: %v", roleType, groupID, err)
	}
	d.SetId(role.Id)
	return resourceGroupRoleRead(ctx, d, m)
}

func resourceGroupRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	groupID := d.Get("group_id").(string)
	rolesAssigned, _, err := getOktaClientFromMetadata(m).Group.ListGroupAssignedRoles(ctx, groupID, nil)
	if err != nil {
		return diag.Errorf("failed to list roles assigned to group %s: %v", groupID, err)
	}
	for _, role := range rolesAssigned {
		if role.Id == d.Id() {
			_ = d.Set("role_type", role.Type)
			return nil
		}
	}
	logger(m).Info("no roles found assigned to group", "group_id", groupID)
	d.SetId("")
	return nil
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
	rolesAssigned, _, err := getOktaClientFromMetadata(m).Group.ListGroupAssignedRoles(ctx, groupID, nil)
	if err != nil {
		return nil, err
	}
	for _, role := range rolesAssigned {
		if role.Id == roleID {
			d.SetId(roleID)
			_ = d.Set("group_id", groupID)
			_ = d.Set("role_type", role.Type)
			return []*schema.ResourceData{d}, nil
		}
	}
	err = fmt.Errorf("unable to find the role ID %s assigned to the group %s", roleID, groupID)
	return nil, err
}
