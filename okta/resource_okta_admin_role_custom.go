package okta

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourceAdminRoleCustom() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAdminRoleCustomCreate,
		ReadContext:   resourceAdminRoleCustomRead,
		UpdateContext: resourceAdminRoleCustomUpdate,
		DeleteContext: resourceAdminRoleCustomDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: `Resource to manage administrative Role assignments for a User

These operations allow the creation and manipulation of custom roles as custom collections of permissions.`,
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
					Type: schema.TypeString,
				},
				Description: `The permissions that the new Role grants. At least one
				permission must be specified when creating custom role. Valid values: "okta.authzServers.manage",
			  "okta.authzServers.read",
			  "okta.apps.assignment.manage",
			  "okta.apps.manage",
			  "okta.apps.read",
			  "okta.customizations.manage",
			  "okta.customizations.read",
			  "okta.groups.appAssignment.manage",
			  "okta.groups.create",
			  "okta.groups.manage",
			  "okta.groups.members.manage",
			  "okta.groups.read",
			  "okta.profilesources.import.run",
			  "okta.users.appAssignment.manage",
			  "okta.users.create",
			  "okta.users.credentials.expirePassword",
			  "okta.users.credentials.manage",
			  "okta.users.credentials.resetFactors",
			  "okta.users.credentials.resetPassword",
			  "okta.users.groupMembership.manage",
			  "okta.users.lifecycle.activate",
			  "okta.users.lifecycle.clearSessions",
			  "okta.users.lifecycle.deactivate",
			  "okta.users.lifecycle.delete",
			  "okta.users.lifecycle.manage",
			  "okta.users.lifecycle.suspend",
			  "okta.users.lifecycle.unlock",
			  "okta.users.lifecycle.unsuspend",
			  "okta.users.manage",
			  "okta.users.read",
			  "okta.users.userprofile.manage",
			  "okta.workflows.invoke".,`,
			},
		},
	}
}

func resourceAdminRoleCustomCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cr, err := buildCustomAdminRole(d, true)
	if err != nil {
		return diag.Errorf("failed to create custom admin role: %v", err)
	}
	role, _, err := getAPISupplementFromMetadata(m).CreateCustomRole(ctx, *cr)
	if err != nil {
		return diag.Errorf("failed to create custom admin role: %v", err)
	}
	d.SetId(role.Id)
	return nil
}

func resourceAdminRoleCustomRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	role, resp, err := getAPISupplementFromMetadata(m).GetCustomRole(ctx, d.Id())
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
	perms, _, err := getAPISupplementFromMetadata(m).ListCustomRolePermissions(ctx, d.Id())
	if err != nil {
		return diag.Errorf("failed to list permissions for custom admin role: %v", err)
	}
	_ = d.Set("permissions", flattenPermissions(perms.Permissions))
	return nil
}

func resourceAdminRoleCustomUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getAPISupplementFromMetadata(m)
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
	resp, err := getAPISupplementFromMetadata(m).DeleteCustomRole(ctx, d.Id())
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
