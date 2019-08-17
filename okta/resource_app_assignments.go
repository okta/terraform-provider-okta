package okta

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAppAssignments() *schema.Resource {
	return &schema.Resource{
		// No point in having an exist function, since only the group has to exist
		Create: resourceGroupRolesCreate,
		Read:   resourceGroupRolesRead,
		Update: resourceGroupRolesUpdate,
		Delete: resourceGroupRolesDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				d.Set("app_id", d.Id())
				d.SetId(getGroupRoleId(d.Id()))
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"app_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of group to attach admin roles to",
			},
			"groups": &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Groups associated with the application",
			},
			"users": &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        appUserResource,
				Description: "Users associated with the application",
			},
		},
	}
}

func resourceGroupAppAssignmentCreate(d *schema.ResourceData, m interface{}) error {
	appId := d.Get("app_id").(string)
	adminRoles := convertInterfaceToStringSet(d.Get("admin_roles"))

	for _, role := range adminRoles {
		groupRole := buildGroupRole(d, role)
		_, _, err := getSupplementFromMetadata(m).CreateAdminRole(appId, groupRole, nil)
		if err != nil {
			return err
		}
	}

	d.SetId(getGroupRoleId(appId))

	return resourceGroupAppAssignmentRead(d, m)
}

func resourceGroupAppAssignmentRead(d *schema.ResourceData, m interface{}) error {
	groupId := d.Get("group_id").(string)
	existingRoles, _, err := getSupplementFromMetadata(m).ListAdminRoles(groupId, nil)
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

func resourceGroupAppAssignmentUpdate(d *schema.ResourceData, m interface{}) error {
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

	return resourceGroupAppAssignmentRead(d, m)
}

func resourceGroupAppAssignmentDelete(d *schema.ResourceData, m interface{}) error {
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
