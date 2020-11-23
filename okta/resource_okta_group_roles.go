package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
)

func resourceGroupRoles() *schema.Resource {
	return &schema.Resource{
		// No point in having an exist function, since only the group has to exist
		Create: resourceGroupRolesCreate,
		Read:   resourceGroupRolesRead,
		Update: resourceGroupRolesUpdate,
		Delete: resourceGroupRolesDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Admin roles associated with the group. This can also be done per user.",
			},
		},
	}
}

func buildGroupRole(role string) *sdk.Role {
	return &sdk.Role{
		AssignmentType: "GROUP",
		Type:           role,
	}
}

func containsRole(roles []*sdk.Role, roleName string) bool {
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

func resourceGroupRolesCreate(d *schema.ResourceData, m interface{}) error {
	groupID := d.Get("group_id").(string)
	adminRoles := convertInterfaceToStringSet(d.Get("admin_roles"))

	for _, role := range adminRoles {
		_, _, err := getOktaClientFromMetadata(m).Group.AssignRoleToGroup(context.Background(), groupID, okta.AssignRoleRequest{
			Type: role,
		}, nil)
		if err != nil {
			return err
		}
	}

	d.SetId(getGroupRoleID(groupID))

	return resourceGroupRolesRead(d, m)
}

func resourceGroupRolesRead(d *schema.ResourceData, m interface{}) error {
	groupID := d.Get("group_id").(string)
	existingRoles, resp, err := getSupplementFromMetadata(m).ListAdminRoles(groupID, nil)

	if resp != nil && is404(resp.StatusCode) {
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	adminRoles := make([]string, len(existingRoles))

	for i, role := range existingRoles {
		adminRoles[i] = role.Type
	}
	_ = d.Set("admin_roles", convertStringSetToInterface(adminRoles))

	return nil
}

func resourceGroupRolesUpdate(d *schema.ResourceData, m interface{}) error {
	client := getSupplementFromMetadata(m)
	groupID := d.Get("group_id").(string)
	existingRoles, _, err := client.ListAdminRoles(groupID, nil)
	if err != nil {
		return err
	}
	adminRoles := convertInterfaceToStringSet(d.Get("admin_roles"))

	rolesToAdd, rolesToRemove := splitRoles(existingRoles, adminRoles)

	for _, role := range rolesToAdd {
		groupRole := buildGroupRole(role)
		_, _, err := client.CreateAdminRole(groupID, groupRole, nil)

		if err != nil {
			return err
		}
	}

	for _, roleID := range rolesToRemove {
		_, err := client.DeleteAdminRole(groupID, roleID)

		if err != nil {
			return err
		}
	}

	return resourceGroupRolesRead(d, m)
}

func resourceGroupRolesDelete(d *schema.ResourceData, m interface{}) error {
	client := getSupplementFromMetadata(m)
	groupID := d.Get("group_id").(string)
	existingRoles, _, err := client.ListAdminRoles(groupID, nil)
	if err != nil {
		return err
	}

	for _, role := range existingRoles {
		_, err := client.DeleteAdminRole(groupID, role.Id)

		if err != nil {
			return err
		}
	}

	return nil
}

func splitRoles(existingRoles []*sdk.Role, expectedRoles []string) (rolesToAdd, rolesToRemove []string) {
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
