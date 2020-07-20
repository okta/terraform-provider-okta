package okta

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	//"github.com/terraform-providers/terraform-provider-okta/sdk"
)

func resourceGroupRoleGroupTarget() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupRoleGroupTargetCreate,
		//Exists: resourceGroupRoleGroupTargetExists,
		Read:   resourceGroupRoleGroupTargetRead,
		Update: resourceGroupRoleGroupTargetUpdate,
		Delete: resourceGroupRoleGroupTargetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"group_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of source group",
			},
			"role_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of source group role to attach group target to",
			},
			"target_group_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of target group",
			},
		},
	}
}

func getGroupRoleGroupTargetId(groupId, roleId, targetGroupId string) string {
	return fmt.Sprintf("%s/role/%s/target/group/%s", groupId, roleId, targetGroupId)
}

func resourceGroupRoleGroupTargetCreate(d *schema.ResourceData, m interface{}) error {
	client := getSupplementFromMetadata(m)
	groupId := d.Get("group_id").(string)
	roleId := d.Get("role_id").(string)
	targetGroupId := d.Get("target_group_id").(string)

	_, err := client.CreateAdminRoleGroupTarget(groupId, roleId, targetGroupId)
	if err != nil {
		return err
	}

	d.SetId(getGroupRoleGroupTargetId(groupId, roleId, targetGroupId))

	return resourceGroupRoleGroupTargetRead(d, m)
}

func resourceGroupRoleGroupTargetRead(d *schema.ResourceData, m interface{}) error {
	client := getSupplementFromMetadata(m)
	groupId := d.Get("group_id").(string)
	roleId := d.Get("role_id").(string)
	targetGroupId := d.Get("target_group_id").(string)

	existingGroupTargets, resp, err := client.ListAdminRoleGroupTargets(groupId, roleId, nil)

	if is404(resp.StatusCode) {
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	for _, groupTarget := range existingGroupTargets {
		if groupTarget.Id == targetGroupId {
			return nil
		}
	}

	// Target must've been deleted, set null
	d.SetId("")
	return nil
}

func resourceGroupRoleGroupTargetUpdate(d *schema.ResourceData, m interface{}) error {
	// There's nothing to update...
	return resourceGroupRoleRead(d, m)
}

func resourceGroupRoleGroupTargetDelete(d *schema.ResourceData, m interface{}) error {
	client := getSupplementFromMetadata(m)
	groupId := d.Get("group_id").(string)
	roleId := d.Get("role_id").(string)
	targetGroupId := d.Get("target_group_id").(string)

	_, err := client.DeleteAdminRoleGroupTarget(groupId, roleId, targetGroupId)

	if err != nil {
		return err
	}

	return nil
}
