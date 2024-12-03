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
			"permission_conditions": {
				Type:     schema.TypeList,
				Optional: true,
				Description: "Use a Condition object to further restrict a permission in a Custom Admin Role. For example, you can restrict access to specific profile attributes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"permission": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The Permission/Object to which the conditions apply",
						},
						"include": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Permission/Object Attributes to which access is allowed",
						},
						"exclude": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Permission/Object Attributes to which access isn't allowed",
						},
					},
				},
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

	// Handle permission conditions
	if conditions, ok := d.GetOk("permission_conditions"); ok {
		err := managePermissionConditions(ctx, getAPISupplementFromMetadata(m), d.Id(), conditions.([]interface{}))
		if err != nil {
			return diag.Errorf("failed to manage permission conditions: %v", err)
		}
	}

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
	if d.HasChange("permissions") {
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
	}

	// Handle permission conditions update
	if d.HasChange("permission_conditions") {
		conditions := d.Get("permission_conditions").([]interface{})
		err := managePermissionConditions(ctx, client, d.Id(), conditions)
		if err != nil {
			return diag.Errorf("failed to manage permission conditions: %v", err)
		}
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

func managePermissionConditions(ctx context.Context, client *sdk.APISupplement, roleId string, conditions []interface{}) error {
	for _, condition := range conditions {
		condMap := condition.(map[string]interface{})
		permission := condMap["permission"].(string)
		include := convertInterfaceArrToStringArr(condMap["include"].([]interface{}))
		exclude := convertInterfaceArrToStringArr(condMap["exclude"].([]interface{}))

		// Create the condition payload
		permissionCondition := sdk.PermissionCondition{
			Include: map[string][]string{
				permission: include,
			},
			Exclude: map[string][]string{
				permission: exclude,
			},
		}

		// Update the permission with conditions
		_, err := client.UpdateCustomRolePermissionCondition(ctx, roleId, permission, permissionCondition)
		if err != nil {
			return fmt.Errorf("failed to update permission condition for %s: %v", permission, err)
		}
	}
	return nil
}
