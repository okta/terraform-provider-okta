package idaas

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/okta-sdk-golang/v5/okta"
)

var permissionsWithConditionSupport = []string{
	"okta.users.userprofile.manage",
	"okta.users.read",
}

func resourceAdminRoleCustom() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAdminRoleCustomCreate,
		ReadContext:   resourceAdminRoleCustomRead,
		UpdateContext: resourceAdminRoleCustomUpdate,
		DeleteContext: resourceAdminRoleCustomDelete,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
		Description:   `Resource to create and manage a custom okta admin role that can be assigned to a user or group to grant a collection of permissions.`,
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
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The permissions granted by the custom admin role. At least one permission must be specified when creating a custom admin role. see: https://developer.okta.com/docs/api/openapi/okta-management/guides/permissions`,
				MinItems:    1,
			},
			"permission_conditions": {
				Type:     schema.TypeSet,
				Optional: true,
				Set:      hashPermissionCondition,
				Description: `Use a Condition object to further restrict a permission in a Custom Admin Role. 
Permissions that support conditions:
- okta.users.userprofile.manage
- okta.users.read

Note: support for new permission conditions may be added in the future.
A permission condition can have either include or exclude, but not both.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"permission": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The Permission/Object to which the conditions apply",
						},
						"include": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "JSON-encoded map of references to attributes that are allowed.",
						},
						"exclude": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "JSON-encoded map of references to attributes that aren't allowed.",
						},
					},
				},
			},
		},
	}
}

func validatePermissionConditions(ctx context.Context, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	var permissions []string
	if permissionsRaw, ok := d.GetOk("permissions"); ok {
		permissions = utils.ConvertInterfaceArrToStringArr(permissionsRaw.(*schema.Set).List())
	}

	if conditions, ok := d.GetOk("permission_conditions"); ok {
		for _, condition := range conditions.(*schema.Set).List() {
			condMap := condition.(map[string]interface{})
			permission := condMap["permission"].(string)

			if !slices.Contains(permissions, permission) {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("Permission %q is used in conditions but not present in permissions", permission),
					Detail:   "Ensure that each permission referenced in a condition is also listed in the permissions block.",
				})
			}

			if !slices.Contains(permissionsWithConditionSupport, permission) {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("Permission %q might not currently support permission conditions", permission),
					Detail:   "Terraform will attempt to apply this configuration, but it may fail if the specified permission does not currently support conditions. New permissions may support conditions in future Okta versions.",
				})
			}

			// Check if both include and exclude have non-empty values
			includeStr := condMap["include"].(string)
			excludeStr := condMap["exclude"].(string)
			if includeStr != "" && excludeStr != "" {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("Permission %q has both include and exclude conditions", permission),
					Detail:   "A permission condition can only have either include or exclude, but not both. This is a limitation of the Okta API.",
				})
			}

			// Validate JSON format for include and exclude
			if includeStr != "" {
				var includeMap map[string]interface{}
				if err := json.Unmarshal([]byte(includeStr), &includeMap); err != nil {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  fmt.Sprintf("Invalid JSON in include condition for permission %q", permission),
						Detail:   fmt.Sprintf("The include condition must be a valid JSON string: %v", err),
					})
				}
			}

			if excludeStr != "" {
				var excludeMap map[string]interface{}
				if err := json.Unmarshal([]byte(excludeStr), &excludeMap); err != nil {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  fmt.Sprintf("Invalid JSON in exclude condition for permission %q", permission),
						Detail:   fmt.Sprintf("The exclude condition must be a valid JSON string: %v", err),
					})
				}
			}
		}
	}

	return diags
}

func resourceAdminRoleCustomCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	diags := validatePermissionConditions(ctx, d)
	if diags.HasError() {
		return diags
	}

	client := getOktaV5ClientFromMetadata(m)

	cr := &okta.CreateIamRoleRequest{
		Label:       d.Get("label").(string),
		Description: d.Get("description").(string),
	}

	// Convert permissions set to string array
	if v, ok := d.GetOk("permissions"); ok && v != nil {
		perms := make([]string, 0, v.(*schema.Set).Len())
		for _, item := range v.(*schema.Set).List() {
			perms = append(perms, item.(string))
		}
		cr.Permissions = perms
	}

	role, _, err := client.RoleAPI.CreateRole(ctx).Instance(*cr).Execute()
	if err != nil {
		return append(diags, diag.Errorf("failed to create custom admin role: %v", err)...)
	}
	d.SetId(role.GetId())

	// Handle permission conditions
	if conditions, ok := d.GetOk("permission_conditions"); ok {
		err := managePermissionConditions(ctx, client, d.Id(), conditions.(*schema.Set).List())
		if err != nil {
			return append(diags, diag.Errorf("failed to update permission conditions for role %s: %v", d.Id(), err)...)
		}
	}

	return diags
}

func resourceAdminRoleCustomRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getOktaV5ClientFromMetadata(m)

	role, resp, err := client.RoleAPI.GetRole(ctx, d.Id()).Execute()
	if err := utils.SuppressErrorOn404_V5(resp, err); err != nil {
		return diag.Errorf("failed to find custom admin role: %v", err)
	}
	if role == nil {
		d.SetId("")
		return nil
	}

	// If role label was used instead of ID during import
	if role.GetId() != d.Id() {
		d.SetId(role.GetId())
	}

	if err := d.Set("label", role.GetLabel()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", role.GetDescription()); err != nil {
		return diag.FromErr(err)
	}

	perms, _, err := client.RoleAPI.ListRolePermissions(ctx, d.Id()).Execute()
	if err != nil {
		return diag.Errorf("failed to list permissions for custom admin role: %v", err)
	}

	permList := make([]string, 0)
	conditions := make([]map[string]interface{}, 0)

	for _, perm := range perms.GetPermissions() {
		permLabel := perm.GetLabel()
		permList = append(permList, permLabel)

		if perm.HasConditions() {
			cond := map[string]interface{}{
				"permission": permLabel,
			}

			if conditionMap, ok := perm.GetConditionsOk(); ok && conditionMap != nil {
				// Handle include conditions
				if includeMap, ok := conditionMap["include"].(map[string]interface{}); ok && len(includeMap) > 0 {
					jsonInclude, err := json.Marshal(includeMap)
					if err != nil {
						return diag.Errorf("failed to marshal include conditions for permission %s: %v", permLabel, err)
					}
					cond["include"] = string(jsonInclude)
				}

				// Handle exclude conditions
				if excludeMap, ok := conditionMap["exclude"].(map[string]interface{}); ok && len(excludeMap) > 0 {
					jsonExclude, err := json.Marshal(excludeMap)
					if err != nil {
						return diag.Errorf("failed to marshal exclude conditions for permission %s: %v", permLabel, err)
					}
					cond["exclude"] = string(jsonExclude)
				}

				conditions = append(conditions, cond)
			}
		}
	}

	if err := d.Set("permissions", schema.NewSet(schema.HashString, utils.ConvertStringSliceToInterfaceSlice(permList))); err != nil {
		return diag.FromErr(err)
	}
	condInterfaces := make([]interface{}, len(conditions))
	for i, c := range conditions {
		condInterfaces[i] = c
	}
	if err := d.Set("permission_conditions", schema.NewSet(hashPermissionCondition, condInterfaces)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceAdminRoleCustomUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	diags := validatePermissionConditions(ctx, d)
	if diags.HasError() {
		return diags
	}

	client := getOktaV5ClientFromMetadata(m)

	if d.HasChange("description") || d.HasChange("label") {
		ur := &okta.UpdateIamRoleRequest{
			Description: d.Get("description").(string),
			Label:       d.Get("label").(string),
		}
		_, _, err := client.RoleAPI.ReplaceRole(ctx, d.Id()).Instance(*ur).Execute()
		if err != nil {
			return append(diags, diag.Errorf("failed to update custom admin role: %v", err)...)
		}
	}

	if d.HasChange("permissions") {
		oldPermissions, newPermissions := d.GetChange("permissions")
		oldSet := oldPermissions.(*schema.Set)
		newSet := newPermissions.(*schema.Set)

	permissionsToAdd := utils.ConvertInterfaceArrToStringArr(newSet.Difference(oldSet).List())
	permissionsToRemove := utils.ConvertInterfaceArrToStringArr(oldSet.Difference(newSet).List())

		if err := addCustomRolePermissions(ctx, client, d.Id(), permissionsToAdd); err != nil {
			return append(diags, err...)
		}
		if err := removeCustomRolePermissions(ctx, client, d.Id(), permissionsToRemove); err != nil {
			return append(diags, err...)
		}
	}

	// Handle permission conditions update
	if d.HasChange("permission_conditions") {
		conditions := d.Get("permission_conditions").(*schema.Set).List()
		if err := managePermissionConditions(ctx, client, d.Id(), conditions); err != nil {
			return append(diags, diag.Errorf("failed to manage permission conditions: %v", err)...)
		}
	}

	return resourceAdminRoleCustomRead(ctx, d, m)
}

func resourceAdminRoleCustomDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getOktaV5ClientFromMetadata(m)
	resp, err := client.RoleAPI.DeleteRole(ctx, d.Id()).Execute()
	if err := utils.SuppressErrorOn404_V5(resp, err); err != nil {
		return diag.Errorf("failed to delete admin custom role: %v", err)
	}
	return nil
}

func addCustomRolePermissions(ctx context.Context, client *okta.APIClient, roleIdOrLabel string, permissions []string) diag.Diagnostics {
	for _, permission := range permissions {
		_, err := client.RoleAPI.CreateRolePermission(ctx, roleIdOrLabel, permission).Execute()
		if err != nil {
			return diag.Errorf("failed to add %s permission to the custom role %s: %v", permission, roleIdOrLabel, err)
		}
	}
	return nil
}

func removeCustomRolePermissions(ctx context.Context, client *okta.APIClient, roleIdOrLabel string, permissions []string) diag.Diagnostics {
	for _, permission := range permissions {
		_, err := client.RoleAPI.DeleteRolePermission(ctx, roleIdOrLabel, permission).Execute()
		if err != nil {
			return diag.Errorf("failed to remove %s permission from the custom role %s: %v", permission, roleIdOrLabel, err)
		}
	}
	return nil
}

func managePermissionConditions(ctx context.Context, client *okta.APIClient, roleIdOrLabel string, conditions []interface{}) diag.Diagnostics {
	for _, condition := range conditions {
		condMap := condition.(map[string]interface{})
		permission := condMap["permission"].(string)

		// Build conditions map in the format the API expects
		conditionData := map[string]interface{}{
			"conditions": map[string]interface{}{},
		}

		// Handle include conditions
		if includeStr, ok := condMap["include"].(string); ok && includeStr != "" {
			var includeMap map[string]interface{}
			if err := json.Unmarshal([]byte(includeStr), &includeMap); err != nil {
				return diag.Errorf("failed to parse include conditions JSON for permission %s: %v", permission, err)
			}
			conditionData["conditions"].(map[string]interface{})["include"] = includeMap
		}

		// Handle exclude conditions
		if excludeStr, ok := condMap["exclude"].(string); ok && excludeStr != "" {
			var excludeMap map[string]interface{}
			if err := json.Unmarshal([]byte(excludeStr), &excludeMap); err != nil {
				return diag.Errorf("failed to parse exclude conditions JSON for permission %s: %v", permission, err)
			}
			conditionData["conditions"].(map[string]interface{})["exclude"] = excludeMap
		}

		req := okta.CreateUpdateIamRolePermissionRequest{
			Conditions: conditionData["conditions"].(map[string]interface{}),
		}

		_, _, err := client.RoleAPI.ReplaceRolePermission(ctx, roleIdOrLabel, permission).Instance(req).Execute()
		if err != nil {
			return diag.Errorf("failed to update permission condition for %s: %v", permission, err)
		}
	}
	return nil
}

func hashPermissionCondition(v interface{}) int {
	m := v.(map[string]interface{})
	perm := m["permission"].(string)
	key := fmt.Sprintf("perm:%s|include:%v|exclude:%v", perm, m["include"], m["exclude"])
	return schema.HashString(key)
}
