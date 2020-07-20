package okta

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-okta/sdk"
)

func resourceGroupRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupRoleCreate,
		//Exists: resourceGroupRoleGroupTargetExists,
		Read:   resourceGroupRoleRead,
		Update: resourceGroupRoleUpdate,
		Delete: resourceGroupRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"group_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of group to attach admin roles to",
			},
			"role_type": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Admin role associated with the group.",
			},
		},
	}
}

func buildGroupRole(d *schema.ResourceData, roleType string) *sdk.Role {
	return &sdk.Role{
		AssignmentType: "GROUP",
		Type:           roleType,
	}
}

func resourceGroupRoleCreate(d *schema.ResourceData, m interface{}) error {
	client := getSupplementFromMetadata(m)
	groupId := d.Get("group_id").(string)
	roleType := d.Get("role_type").(string)

	groupRole := buildGroupRole(d, roleType)
	role, _, err := client.CreateAdminRole(groupId, groupRole, nil)
	if err != nil {
		return err
	}

	d.SetId(role.Id)

	return resourceGroupRoleRead(d, m)
}

func resourceGroupRoleRead(d *schema.ResourceData, m interface{}) error {
	client := getSupplementFromMetadata(m)
	groupId := d.Get("group_id").(string)
	roleId := d.Id()
	existingRoles, resp, err := client.ListAdminRoles(groupId, nil)

	if is404(resp.StatusCode) {
		// Group has been deleted, Role must've been deleted, set null
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	for _, role := range existingRoles {
		if role.Id == roleId {
			return nil
		}
	}

	// Role must've been deleted, set null
	d.SetId("")
	return nil
}

func resourceGroupRoleUpdate(d *schema.ResourceData, m interface{}) error {
	// There's nothing to update...
	return resourceGroupRoleRead(d, m)
}

func resourceGroupRoleDelete(d *schema.ResourceData, m interface{}) error {
	client := getSupplementFromMetadata(m)
	groupId := d.Get("group_id").(string)
	roleId := d.Id()

	_, err := client.DeleteAdminRole(groupId, roleId)

	if err != nil {
		return err
	}

	return nil
}
