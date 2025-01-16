package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUserAdminRoles() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserAdminRolesCreate,
		ReadContext:   resourceUserAdminRolesRead,
		UpdateContext: resourceUserAdminRolesUpdate,
		DeleteContext: resourceUserAdminRolesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				_ = d.Set("user_id", d.Id())
				d.SetId(d.Id())
				return []*schema.ResourceData{d}, nil
			},
		},
		Description: "Resource to manage a set of administrator roles for a specific user. This resource allows you to manage admin roles for a single user, independent of the user schema itself.",
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of a Okta User",
				ForceNew:    true,
			},
			"admin_roles": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "The list of Okta user admin roles, e.g. `['APP_ADMIN', 'USER_ADMIN']` See [API Docs](https://developer.okta.com/docs/api/openapi/okta-management/guides/roles/#standard-roles).",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"disable_notifications": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "When this setting is enabled, the admins won't receive any of the default Okta administrator emails. These admins also won't have access to contact Okta Support and open support cases on behalf of your org.",
				Default:     false,
			},
		},
	}
}

func resourceUserAdminRolesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	userId := d.Get("user_id").(string)
	roles := convertInterfaceToStringSetNullable(d.Get("admin_roles"))
	client := getOktaClientFromMetadata(m)
	err := assignAdminRolesToUser(ctx, userId, roles, d.Get("disable_notifications").(bool), client)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(userId)
	return nil
}

func resourceUserAdminRolesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := setAdminRoles(ctx, d, m)
	if err != nil {
		return diag.Errorf("failed to set read user's roles: %v", err)
	}
	return nil
}

func resourceUserAdminRolesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	userId := d.Get("user_id").(string)
	client := getOktaClientFromMetadata(m)
	if !d.HasChange("admin_roles") && d.HasChange("disable_notifications") {
		roles := convertInterfaceToStringSet(d.Get("admin_roles"))
		// we just need to update an existing role assignment status by passing just the query parameter.
		if err := assignAdminRolesToUser(ctx, userId, roles, d.Get("disable_notifications").(bool), client); err != nil {
			return diag.Errorf("failed to update user admin roles: %v", err)
		}
		return nil
	}
	oldRoles, newRoles := d.GetChange("admin_roles")
	oldSet := oldRoles.(*schema.Set)
	newSet := newRoles.(*schema.Set)
	rolesToAdd := convertInterfaceArrToStringArr(newSet.Difference(oldSet).List())
	rolesToRemove := convertInterfaceArrToStringArr(oldSet.Difference(newSet).List())
	roles, _, err := listUserOnlyRoles(ctx, client, d.Id())
	if err != nil {
		return diag.Errorf("failed to list user's roles: %v", err)
	}
	for _, role := range roles {
		if contains(rolesToRemove, role.Type) {
			resp, err := client.User.RemoveRoleFromUser(ctx, d.Id(), role.Id)
			if err := suppressErrorOn404(resp, err); err != nil {
				return diag.Errorf("failed to remove user's role: %v", err)
			}
		}
	}
	err = assignAdminRolesToUser(ctx, d.Id(), rolesToAdd, d.Get("disable_notifications").(bool), client)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceUserAdminRolesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	userID := d.Get("user_id").(string)
	client := getOktaClientFromMetadata(m)
	roles, _, err := listUserOnlyRoles(ctx, client, userID)
	if err != nil {
		return diag.Errorf("failed to list user's roles: %v", err)
	}
	for _, role := range roles {
		resp, err := client.User.RemoveRoleFromUser(ctx, userID, role.Id)
		if err := suppressErrorOn404(resp, err); err != nil {
			return diag.Errorf("failed to remove user's role: %v", err)
		}
	}
	return nil
}
