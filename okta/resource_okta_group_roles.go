package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func resourceGroupRoles() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This resource is deprecated and will be removed in favor of using \"okta_group_role\", please migrate as soon as possible.",
		CreateContext:      resourceGroupRolesCreate,
		ReadContext:        resourceGroupRolesRead,
		UpdateContext:      resourceGroupRolesUpdate,
		DeleteContext:      resourceGroupRolesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				_ = d.Set("group_id", d.Id())
				d.SetId(getGroupRoleID(d.Id()))
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of group to attach admin roles to",
			},
			"admin_roles": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: elemInSlice(validAdminRoles),
				},
				Description: "Admin roles associated with the group. This can also be done per user.",
			},
		},
	}
}

func resourceGroupRolesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	groupID := d.Get("group_id").(string)
	adminRoles := convertInterfaceToStringSet(d.Get("admin_roles"))
	for _, role := range adminRoles {
		_, _, err := getOktaClientFromMetadata(m).Group.AssignRoleToGroup(ctx, groupID, okta.AssignRoleRequest{Type: role}, nil)
		if err != nil {
			return diag.Errorf("failed to assign role %s to group %s: %v", role, groupID, err)
		}
	}
	d.SetId(getGroupRoleID(groupID))
	return resourceGroupRolesRead(ctx, d, m)
}

func resourceGroupRolesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	groupID := d.Get("group_id").(string)
	existingRoles, resp, err := getOktaClientFromMetadata(m).Group.ListGroupAssignedRoles(ctx, groupID, nil)
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get list of group assigned roles: %v", err)
	}
	adminRoles := make([]string, len(existingRoles))
	for i, role := range existingRoles {
		adminRoles[i] = role.Type
	}
	_ = d.Set("admin_roles", convertStringSliceToSet(adminRoles))
	return nil
}

func resourceGroupRolesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(m)
	groupID := d.Get("group_id").(string)
	existingRoles, resp, err := getOktaClientFromMetadata(m).Group.ListGroupAssignedRoles(ctx, groupID, nil)
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get list of group assigned roles: %v", err)
	}
	adminRoles := convertInterfaceToStringSet(d.Get("admin_roles"))
	rolesToAdd, rolesToRemove := splitRoles(existingRoles, adminRoles)
	for _, role := range rolesToAdd {
		_, _, err := client.Group.AssignRoleToGroup(ctx, groupID, okta.AssignRoleRequest{Type: role}, nil)
		if err != nil {
			return diag.Errorf("failed to assign role %s to group %s: %v", role, groupID, err)
		}
	}
	for _, roleID := range rolesToRemove {
		_, err := client.Group.RemoveRoleFromGroup(ctx, groupID, roleID)
		if err != nil {
			return diag.Errorf("failed to remove role %s from group %s: %v", roleID, groupID, err)
		}
	}
	return resourceGroupRolesRead(ctx, d, m)
}

func resourceGroupRolesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(m)
	groupID := d.Get("group_id").(string)
	existingRoles, resp, err := client.Group.ListGroupAssignedRoles(ctx, groupID, nil)
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get list of group assigned roles: %v", err)
	}
	for _, role := range existingRoles {
		_, err := client.Group.RemoveRoleFromGroup(ctx, groupID, role.Id)
		if err != nil {
			return diag.Errorf("failed to remove role %s from group %s: %v", role.Id, groupID, err)
		}
	}
	return nil
}

func splitRoles(existingRoles []*okta.Role, expectedRoles []string) (rolesToAdd, rolesToRemove []string) {
	for _, roleName := range expectedRoles {
		if !containsRole(existingRoles, roleName) {
			rolesToAdd = append(rolesToAdd, roleName)
		}
	}

	for _, role := range existingRoles {
		if !contains(expectedRoles, role.Type) {
			rolesToRemove = append(rolesToRemove, role.Id)
		}
	}

	return
}

func containsRole(roles []*okta.Role, roleName string) bool {
	for _, role := range roles {
		if role.Type == roleName {
			return true
		}
	}
	return false
}

func getGroupRoleID(groupID string) string {
	return fmt.Sprintf("%s.roles", groupID)
}
