package idaas

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
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
				permission must be specified when creating custom role. Valid values: "okta.users.manage",
				"okta.users.create",
				"okta.users.read",
				"okta.users.credentials.manage",
				"okta.users.credentials.resetFactors",
				"okta.users.credentials.resetPassword",
				"okta.users.credentials.expirePassword",
				"okta.users.userprofile.manage",
				"okta.users.lifecycle.manage",
				"okta.users.lifecycle.activate",
				"okta.users.lifecycle.deactivate",
				"okta.users.lifecycle.suspend",
				"okta.users.lifecycle.unsuspend",
				"okta.users.lifecycle.delete",
				"okta.users.lifecycle.unlock",
				"okta.users.lifecycle.clearSessions",
				"okta.users.groupMembership.manage",
				"okta.users.appAssignment.manage",
				"okta.users.apitokens.manage",
				"okta.users.apitokens.read",
				"okta.groups.manage",
				"okta.groups.create",
				"okta.groups.members.manage",
				"okta.groups.read",
				"okta.groups.appAssignment.manage",
				"okta.apps.read",
				"okta.apps.manage",
				"okta.apps.assignment.manage",
				"okta.profilesources.import.run",
				"okta.authzServers.read",
				"okta.users.userprofile.manage",
				"okta.authzServers.manage",
				"okta.customizations.read",
				"okta.customizations.manage",
				"okta.identityProviders.read",
				"okta.identityProviders.manage",
				"okta.workflows.read",
				"okta.workflows.invoke",
				"okta.governance.accessCertifications.manage",
				"okta.governance.accessRequests.manage",
				"okta.apps.manageFirstPartyApps",
				"okta.agents.manage",
				"okta.agents.register",
				"okta.agents.view",
				"okta.directories.manage",
				"okta.directories.read",
				"okta.devices.manage",
				"okta.devices.lifecycle.manage",
				"okta.devices.lifecycle.activate",
				"okta.devices.lifecycle.deactivate",
				"okta.devices.lifecycle.suspend",
				"okta.devices.lifecycle.unsuspend",
				"okta.devices.lifecycle.delete",
				"okta.devices.read",
				"okta.iam.read",
				"okta.support.cases.manage",`,
			},
		},
	}
}

func resourceAdminRoleCustomCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cr, err := buildCustomAdminRole(d, true)
	if err != nil {
		return diag.Errorf("failed to create custom admin role: %v", err)
	}
	role, _, err := getAPISupplementFromMetadata(meta).CreateCustomRole(ctx, *cr)
	if err != nil {
		return diag.Errorf("failed to create custom admin role: %v", err)
	}
	d.SetId(role.Id)
	return nil
}

func resourceAdminRoleCustomRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	role, resp, err := getAPISupplementFromMetadata(meta).GetCustomRole(ctx, d.Id())
	if err := utils.SuppressErrorOn404(resp, err); err != nil {
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
	perms, _, err := getAPISupplementFromMetadata(meta).ListCustomRolePermissions(ctx, d.Id())
	if err != nil {
		return diag.Errorf("failed to list permissions for custom admin role: %v", err)
	}
	_ = d.Set("permissions", flattenPermissions(perms.Permissions))
	return nil
}

func resourceAdminRoleCustomUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := getAPISupplementFromMetadata(meta)
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

	permissionsToAdd := utils.ConvertInterfaceArrToStringArr(newSet.Difference(oldSet).List())
	permissionsToRemove := utils.ConvertInterfaceArrToStringArr(oldSet.Difference(newSet).List())

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

func resourceAdminRoleCustomDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	resp, err := getAPISupplementFromMetadata(meta).DeleteCustomRole(ctx, d.Id())
	if err := utils.SuppressErrorOn404(resp, err); err != nil {
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
		cr.Permissions = utils.ConvertInterfaceToStringSetNullable(d.Get("permissions"))
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
	// Extract permission labels and normalize them
	permissionLabels := make([]string, len(permissions))
	for i := range permissions {
		permissionLabels[i] = permissions[i].Label
	}

	// Normalize permissions to handle API expansion
	normalizedPermissions := normalizePermissions(permissionLabels)

	arr := make([]interface{}, len(normalizedPermissions))
	for i, perm := range normalizedPermissions {
		arr[i] = perm
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

// normalizePermissions handles API permission expansion by mapping expanded permissions
// back to the user's intended configuration. This prevents Terraform drift when the API
// returns additional permissions alongside the ones explicitly configured.
func normalizePermissions(apiPermissions []string) []string {
	// Create a map to track which permissions we've seen
	permissionMap := make(map[string]bool)
	for _, perm := range apiPermissions {
		permissionMap[perm] = true
	}

	// Track which permissions to include in the normalized result
	normalizedSet := make(map[string]bool)

	// Handle workflow permissions expansion
	// When user configures "okta.workflows.read", API returns both:
	// - "okta.workflows.read" (original)
	// - "okta.workflows.flows.read" (expanded)
	// We preserve the user's original intent
	if permissionMap["okta.workflows.read"] {
		normalizedSet["okta.workflows.read"] = true
		// Don't include the expanded version in normalized output
		delete(permissionMap, "okta.workflows.flows.read")
	} else if permissionMap["okta.workflows.flows.read"] {
		// If only the expanded version exists, keep it
		normalizedSet["okta.workflows.flows.read"] = true
	}

	// Handle workflow invoke permissions expansion
	// When user configures "okta.workflows.invoke", API returns both:
	// - "okta.workflows.invoke" (original)
	// - "okta.workflows.flows.invoke" (expanded)
	// We preserve the user's original intent
	if permissionMap["okta.workflows.invoke"] {
		normalizedSet["okta.workflows.invoke"] = true
		// Don't include the expanded version in normalized output
		delete(permissionMap, "okta.workflows.flows.invoke")
	} else if permissionMap["okta.workflows.flows.invoke"] {
		// If only the expanded version exists, keep it
		normalizedSet["okta.workflows.flows.invoke"] = true
	}

	// Add all other permissions that weren't handled by the workflow normalization
	for perm := range permissionMap {
		if perm != "okta.workflows.flows.read" && perm != "okta.workflows.flows.invoke" {
			normalizedSet[perm] = true
		}
	}

	// Convert back to slice
	result := make([]string, 0, len(normalizedSet))
	for perm := range normalizedSet {
		result = append(result, perm)
	}

	return result
}
