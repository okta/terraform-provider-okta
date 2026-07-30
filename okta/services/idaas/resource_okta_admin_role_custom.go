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
				"okta.workflows.invoke".
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
				"okta.support.cases.manage".,`,
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
	// The Okta API can return both the legacy workflow permission labels
	// (okta.workflows.invoke / okta.workflows.read) and their newer aliases
	// (okta.workflows.flows.invoke / okta.workflows.flows.read). Keep whichever
	// form is already recorded in state so we don't introduce perpetual drift.
	statePermissions := utils.ConvertInterfaceToStringSetNullable(d.Get("permissions"))
	apiPermissions := reconcileWorkflowPermissions(statePermissions, perms.Permissions)
	_ = d.Set("permissions", flattenPermissions(apiPermissions))
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
	arr := make([]interface{}, len(permissions))
	for i := range permissions {
		arr[i] = permissions[i].Label
	}
	return schema.NewSet(schema.HashString, arr)
}

// reconcileWorkflowPermissions removes a redundant workflow permission alias
// from the API response when state already tracks the other form. The Okta API
// exposes the same workflow permission under a legacy label and a newer alias:
//
//	okta.workflows.invoke <-> okta.workflows.flows.invoke
//	okta.workflows.read   <-> okta.workflows.flows.read
//
// Whichever form the user already has in state wins, so the opposite form is
// discarded from the API response before it is written back to state.
func reconcileWorkflowPermissions(statePermissions []string, apiPermissions []*sdk.Permission) []*sdk.Permission {
	inState := make(map[string]bool, len(statePermissions))
	for _, p := range statePermissions {
		inState[p] = true
	}

	existingToNew := map[string]string{
		"okta.workflows.invoke": "okta.workflows.flows.invoke",
		"okta.workflows.read":   "okta.workflows.flows.read",
	}

	discard := make(map[string]bool)
	for legacy, modern := range existingToNew {
		switch {
		case inState[legacy]:
			// state uses the legacy label, discard the newer alias from the API response
			discard[modern] = true
		case inState[modern]:
			// state uses the newer alias, discard the legacy label from the API response
			discard[legacy] = true
		}
	}
	if len(discard) == 0 {
		return apiPermissions
	}

	filtered := make([]*sdk.Permission, 0, len(apiPermissions))
	for _, p := range apiPermissions {
		if p != nil && discard[p.Label] {
			continue
		}
		filtered = append(filtered, p)
	}
	return filtered
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
