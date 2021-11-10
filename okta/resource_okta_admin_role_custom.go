package okta

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

var validCustomRolePermissions = []string{
	"okta.users.manage", "okta.users.create", "okta.users.read", "okta.users.credentials.manage",
	"okta.users.userprofile.manage", "okta.users.lifecycle.manage", "okta.users.groupMembership.manage",
	"okta.users.appAssignment.manage", "okta.groups.manage", "okta.groups.create", "okta.groups.members.manage",
	"okta.groups.read", "okta.groups.appAssignment.manage", "okta.apps.read", "okta.apps.manage",
	"okta.apps.assignment.manage",
}

func resourceAdminRoleCustom() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAdminRoleCustomCreate,
		ReadContext:   resourceAdminRoleCustomRead,
		UpdateContext: resourceAdminRoleCustomUpdate,
		DeleteContext: resourceAdminRoleCustomDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "Resource to manage administrative Role assignments for a User",
		Schema: map[string]*schema.Schema{
			"label": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name given to the new Role",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A human-readable description of the new Role",
			},
			"permissions": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: elemInSlice(validCustomRolePermissions),
				},
				Description: "The permissions that the new Role grants.",
			},
		},
	}
}

func resourceAdminRoleCustomCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cr, err := buildCustomAdminRole(d, true)
	if err != nil {
		return diag.Errorf("failed to create custom admin role: %v", err)
	}
	role, _, err := getSupplementFromMetadata(m).CreateCustomRole(ctx, *cr)
	if err != nil {
		return diag.Errorf("failed to create custom admin role: %v", err)
	}
	d.SetId(role.Id)
	return nil
}

func resourceAdminRoleCustomRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	role, resp, err := getSupplementFromMetadata(m).GetCustomRole(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to find custom admin role: %v", err)
	}
	if role == nil {
		d.SetId("")
		return nil
	}
	// is case role label was used instead of ID for the import
	if role.Id != d.Id() {
		d.SetId(role.Id)
	}
	_ = d.Set("label", role.Label)
	_ = d.Set("description", role.Description)
	perms, _, err := getSupplementFromMetadata(m).ListCustomRolePermissions(ctx, d.Id())
	if err != nil {
		return diag.Errorf("failed to list permissions for custom admin role: %v", err)
	}
	_ = d.Set("permissions", flattenPermissions(perms.Permissions))
	return nil
}

func resourceAdminRoleCustomUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getSupplementFromMetadata(m)
	if d.HasChanges("label", "description") {
		cr, _ := buildCustomAdminRole(d, false)
		_, _, err := client.UpdateCustomRole(ctx, d.Id(), *cr)
		if err != nil {
			return diag.Errorf("failed to update custom admin role: %v", err)
		}
	}
	if !d.HasChange("permissions") {
		return nil
	}
	oldPermissions, newPermissions := d.GetChange("permissions")
	oldSet := oldPermissions.(*schema.Set)
	newSet := newPermissions.(*schema.Set)

	permissionsToAdd := convertInterfaceArrToStringArr(newSet.Difference(oldSet).List())
	permissionsToRemove := convertInterfaceArrToStringArr(oldSet.Difference(newSet).List())

	err := addCustomRolePermissions(ctx, client, d.Id(), permissionsToAdd)
	if err != nil {
		return diag.FromErr(err)
	}
	err = removeCustomRolePermissions(ctx, client, d.Id(), permissionsToRemove)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceAdminRoleCustomDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resp, err := getSupplementFromMetadata(m).DeleteCustomRole(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to delete admin custom role: %v", err)
	}
	return nil
}

func buildCustomAdminRole(d *schema.ResourceData, isNew bool) (*sdk.CustomRole, error) {
	cr := &sdk.CustomRole{
		Label:       d.Get("label").(string),
		Description: d.Get("description").(string),
	}
	if isNew {
		cr.Permissions = convertInterfaceToStringSetNullable(d.Get("permissions"))
		if len(cr.Permissions) == 0 {
			return nil, errors.New("at least one permission must be specified when creating custom role")
		}
	}
	return cr, nil
}

func flattenPermissions(permissions []*sdk.Permission) interface{} {
	if len(permissions) == 0 {
		return nil
	}
	arr := make([]interface{}, len(permissions))
	for i := range permissions {
		arr[i] = permissions[i].Label
	}
	return schema.NewSet(schema.HashString, arr)
}

func addCustomRolePermissions(ctx context.Context, client *sdk.APISupplement, roleIdOrLabel string, permissions []string) error {
	for _, permission := range permissions {
		_, _, err := client.AddCustomRolePermission(ctx, roleIdOrLabel, permission)
		if err != nil {
			return fmt.Errorf("failed to add %s permission to the custom role %s: %v", permission, roleIdOrLabel, err)
		}
	}
	return nil
}

func removeCustomRolePermissions(ctx context.Context, client *sdk.APISupplement, roleIdOrLabel string, permissions []string) error {
	for _, permission := range permissions {
		_, err := client.DeleteCustomRolePermission(ctx, roleIdOrLabel, permission)
		if err != nil {
			return fmt.Errorf("failed to remove %s permission from the custom role %s: %v", permission, roleIdOrLabel, err)
		}
	}
	return nil
}
