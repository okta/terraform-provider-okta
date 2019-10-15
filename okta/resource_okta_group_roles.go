package okta

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-okta/sdk"
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
				d.Set("group_id", d.Id())
				d.SetId(getGroupRoleId(d.Id()))
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"group_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of group to attach admin roles to",
			},
			"admin_roles": &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Admin roles associated with the group. This can also be done per user.",
			},
		},
	}
}

func buildGroupRole(d *schema.ResourceData, role string) *sdk.Role {
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

func getGroupRoleId(groupId string) string {
	return fmt.Sprintf("%s.roles", groupId)
}

func resourceGroupRolesCreate(d *schema.ResourceData, m interface{}) error {
	groupId := d.Get("group_id").(string)
	adminRoles := convertInterfaceToStringSet(d.Get("admin_roles"))

	for _, role := range adminRoles {
		groupRole := buildGroupRole(d, role)
		_, _, err := getSupplementFromMetadata(m).CreateAdminRole(groupId, groupRole, nil)
		if err != nil {
			return err
		}
	}

	d.SetId(getGroupRoleId(groupId))

	return resourceGroupRolesRead(d, m)
}

func resourceGroupRolesRead(d *schema.ResourceData, m interface{}) error {
	groupId := d.Get("group_id").(string)
	existingRoles, resp, err := getSupplementFromMetadata(m).ListAdminRoles(groupId, nil)

	if is404(resp.StatusCode) {
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
	d.Set("admin_roles", convertStringSetToInterface(adminRoles))

	return nil
}

func resourceGroupRolesUpdate(d *schema.ResourceData, m interface{}) error {
	client := getSupplementFromMetadata(m)
	groupId := d.Get("group_id").(string)
	existingRoles, _, err := client.ListAdminRoles(groupId, nil)
	if err != nil {
		return err
	}
	adminRoles := convertInterfaceToStringSet(d.Get("admin_roles"))

	rolesToAdd, rolesToRemove := splitRoles(existingRoles, adminRoles)

	for _, role := range rolesToAdd {
		groupRole := buildGroupRole(d, role)
		_, _, err := client.CreateAdminRole(groupId, groupRole, nil)

		if err != nil {
			return err
		}
	}

	for _, roleId := range rolesToRemove {
		_, err := client.DeleteAdminRole(groupId, roleId)

		if err != nil {
			return err
		}
	}

	return resourceGroupRolesRead(d, m)
}

func resourceGroupRolesDelete(d *schema.ResourceData, m interface{}) error {
	client := getSupplementFromMetadata(m)
	groupId := d.Get("group_id").(string)
	existingRoles, _, err := client.ListAdminRoles(groupId, nil)
	if err != nil {
		return err
	}

	for _, role := range existingRoles {
		_, err := client.DeleteAdminRole(groupId, role.Id)

		if err != nil {
			return err
		}
	}

	return nil
}

func splitRoles(existingRoles []*sdk.Role, expectedRoles []string) (rolesToAdd []string, rolesToRemove []string) {
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
