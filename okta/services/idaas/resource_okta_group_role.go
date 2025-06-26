package idaas

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/okta/terraform-provider-okta/sdk/query"
)

func resourceGroupRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupRoleCreate,
		ReadContext:   resourceGroupRoleRead,
		UpdateContext: resourceGroupRoleUpdate,
		DeleteContext: resourceGroupRoleDelete,
		Importer:      utils.CreateNestedResourceImporter([]string{"group_id", "id"}),
		CustomizeDiff: customdiff.All(
			customdiff.ForceNewIf("target_group_list", func(_ context.Context, d *schema.ResourceDiff, meta interface{}) bool {
				if d.HasChange("target_group_list") {
					// to avoid exception when removing last group target from a role assignment,
					// the API consumer should delete the role assignment and recreate it.
					if len(utils.ConvertInterfaceToStringSet(d.Get("target_group_list"))) == 0 {
						return true
					}
				}
				return false
			}),
			customdiff.ForceNewIf("target_app_list", func(_ context.Context, d *schema.ResourceDiff, meta interface{}) bool {
				if d.HasChange("target_app_list") {
					// to avoid exception when removing last app target from a role assignment,
					// the API consumer should delete the role assignment and recreate it.
					if len(utils.ConvertInterfaceToStringSet(d.Get("target_app_list"))) == 0 {
						return true
					}
					// deleting the app from an org which is the only one left assigned to this group role
					// will cause the role to be unassigned from the group thus resource should be recreated
					oldValue, newValue := d.GetChange("target_app_list")
					oldApps := utils.ConvertInterfaceToStringSet(oldValue)
					newApps := utils.ConvertInterfaceToStringSet(newValue)
					if len(oldApps) > 0 && !utils.ContainsOne(oldApps, newApps...) {
						return true
					}
				}
				return false
			}),
		),
		Description: "Assigns Admin roles to Okta Groups. This resource allows you to assign Okta administrator roles to Okta Groups. This resource provides a one-to-one interface between the Okta group and the admin role.",
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of group to attach admin roles to",
				ForceNew:    true,
			},
			"role_type": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Admin role assigned to the group. It can be any one of the following values:
	"API_ADMIN",
	"APP_ADMIN",
	"CUSTOM",
	"GROUP_MEMBERSHIP_ADMIN",
	"HELP_DESK_ADMIN",
	"MOBILE_ADMIN",
	"ORG_ADMIN",
	"READ_ONLY_ADMIN",
	"REPORT_ADMIN",
	"SUPER_ADMIN",
	"USER_ADMIN"
	. See [API Docs](https://developer.okta.com/docs/api/openapi/okta-management/guides/roles/#standard-roles).
	- "USER_ADMIN" is the Group Administrator.`,
				ForceNew: true,
			},
			"target_group_list": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "A list of group IDs you would like as the targets of the admin role. - Only supported when used with the role types: `GROUP_MEMBERSHIP_ADMIN`, `HELP_DESK_ADMIN`, or `USER_ADMIN`.",
			},
			"target_app_list": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "A list of app names (name represents set of app instances, like 'salesforce' or 'facebook'), or a combination of app name and app instance ID (like 'facebook.0oapsqQ6dv19pqyEo0g3') you would like as the targets of the admin role. - Only supported when used with the role type `APP_ADMIN`.",
			},
			"disable_notifications": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "When this setting is enabled, the admins won't receive any of the default Okta administrator emails. These admins also won't have access to contact Okta Support and open support cases on behalf of your org.",
				Default:     false,
			},
			"role_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Role ID. Required for role_type = `CUSTOM`",
			},
			"resource_set_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Resource Set ID. Required for role_type = `CUSTOM`",
			},
		},
	}
}

func resourceGroupRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	groupID := d.Get("group_id").(string)
	roleType := d.Get("role_type").(string)
	client := getOktaClientFromMetadata(meta)
	logger(meta).Info("assigning role to group", "group_id", groupID, "role_type", roleType)
	role, _, err := client.Group.AssignRoleToGroup(
		ctx,
		groupID,
		sdk.AssignRoleRequest{
			Type:        roleType,
			Role:        d.Get("role_id").(string),
			ResourceSet: d.Get("resource_set_id").(string),
		},
		&query.Params{DisableNotifications: utils.BoolPtr(d.Get("disable_notifications").(bool))})
	if err != nil {
		return diag.Errorf("failed to assign role %s to group %s: %v", roleType, groupID, err)
	}
	groupTargets := utils.ConvertInterfaceToStringSet(d.Get("target_group_list"))
	if len(groupTargets) > 0 && supportsGroupTargets(roleType) {
		logger(meta).Info("scoping admin role assignment to list of groups", "group_id", groupID, "role_id", role.Id, "target_group_list", groupTargets)
		err = addGroupTargetsToRole(ctx, client, groupID, role.Id, groupTargets)
		if err != nil {
			return diag.Errorf("unable to add group target to role assignment %s for group %s: %v", role.Id, groupID, err)
		}
	}
	appTargets := utils.ConvertInterfaceToStringSet(d.Get("target_app_list"))
	if len(appTargets) > 0 && roleType == "APP_ADMIN" {
		logger(meta).Info("scoping admin role assignment to list of apps", "group_id", groupID, "role_id", role.Id, "target_app_list", appTargets)
		err = addGroupAppTargetsToRole(ctx, client, groupID, role.Id, appTargets)
		if err != nil {
			return diag.Errorf("unable to add app targets to role assignment %s for group %s: %v", role.Id, groupID, err)
		}
	}
	d.SetId(role.Id)
	boc := utils.NewExponentialBackOffWithContext(ctx, 10*time.Second)
	err = backoff.Retry(func() error {
		err := resourceGroupRoleRead(ctx, d, meta)
		if err != nil {
			// NOTE: we don't need a doNotRetry(m, err) check because this is going
			// to backoff permanently as-is
			return backoff.Permanent(fmt.Errorf("%s", err[0].Summary))
		}
		if d.Id() != "" {
			return nil
		}
		return fmt.Errorf("role %s was not assigned to a group %s", roleType, groupID)
	}, boc)
	return diag.FromErr(err)
}

func resourceGroupRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	groupID := d.Get("group_id").(string)
	role, resp, err := getOktaClientFromMetadata(meta).Group.GetRole(ctx, groupID, d.Id())
	if err := utils.SuppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get role '%s' assigned to group '%s': %v", d.Id(), groupID, err)
	}
	if role == nil {
		role, err = findRole(ctx, d, meta)
		if err != nil {
			return diag.Errorf("failed to get role '%s' assigned to group '%s': %v", d.Id(), groupID, err)
		}
		if role == nil {
			d.SetId("")
			return nil
		}
	}
	if supportsGroupTargets(role.Type) {
		groupIDs, err := listGroupTargetsIDs(ctx, meta, groupID, role.Id)
		if err != nil {
			return diag.Errorf("unable to list group targets for role %s and group %s: %v", role.Id, groupID, err)
		}
		if len(groupIDs) != 0 {
			_ = d.Set("target_group_list", groupIDs)
		}
	} else if role.Type == "APP_ADMIN" {
		apps, err := listGroupAppsTargets(ctx, d, meta)
		if err != nil {
			return diag.Errorf("unable to list app targets for role %s and group %s: %v", role.Id, groupID, err)
		}
		if len(apps) != 0 {
			_ = d.Set("target_app_list", apps)
		}
	}
	_ = d.Set("role_type", role.Type)
	return nil
}

func resourceGroupRoleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	groupID := d.Get("group_id").(string)
	roleID := d.Id()
	roleType := d.Get("role_type").(string)
	client := getOktaClientFromMetadata(meta)
	if d.HasChange("disable_notifications") {
		_, _, err := client.Group.AssignRoleToGroup(ctx, groupID, sdk.AssignRoleRequest{},
			&query.Params{DisableNotifications: utils.BoolPtr(d.Get("disable_notifications").(bool))})
		if err != nil {
			return diag.Errorf("failed to update group's '%s' notification settings: %v", groupID, err)
		}
	}
	if d.HasChange("target_group_list") && supportsGroupTargets(roleType) {
		expectedGroupIDs := utils.ConvertInterfaceToStringSet(d.Get("target_group_list"))
		existingGroupIDs, err := listGroupTargetsIDs(ctx, meta, groupID, roleID)
		if err != nil {
			return diag.FromErr(err)
		}
		targetsToAdd, targetsToRemove := splitTargets(expectedGroupIDs, existingGroupIDs)
		err = addGroupTargetsToRole(ctx, client, groupID, roleID, targetsToAdd)
		if err != nil {
			return diag.Errorf("failed to add group target to role assignment %s for group %s: %v", roleID, groupID, err)
		}
		err = removeGroupTargetsFromRole(ctx, client, groupID, roleID, targetsToRemove)
		if err != nil {
			return diag.Errorf("failed to remove group target from admin role assignment %s of group %s: %v", roleID, groupID, err)
		}
	}
	if d.HasChange("target_app_list") && roleType == "APP_ADMIN" {
		expectedApps := utils.ConvertInterfaceToStringSet(d.Get("target_app_list"))
		existingApps, err := listGroupAppsTargets(ctx, d, meta)
		if err != nil {
			return diag.Errorf("unable to list app targets for role %s and group %s: %v", d.Id(), groupID, err)
		}
		targetsToAdd, targetsToRemove := splitTargets(expectedApps, existingApps)
		err = addGroupAppTargetsToRole(ctx, client, groupID, roleID, targetsToAdd)
		if err != nil {
			return diag.Errorf("unable to add app target to role assignment %s for group %s: %v", roleID, groupID, err)
		}
		err = removeGroupAppTargets(ctx, client, groupID, roleID, targetsToRemove)
		if err != nil {
			return diag.Errorf("failed to remove app target from admin role assignment %s of group %s: %v", roleID, groupID, err)
		}
	}
	return resourceGroupRoleRead(ctx, d, meta)
}

func resourceGroupRoleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	groupID := d.Get("group_id").(string)
	roleType := d.Get("role_type").(string)
	logger(meta).Info("deleting assigned role from group", "group_id", groupID, "role_type", roleType)
	resp, err := getOktaClientFromMetadata(meta).Group.RemoveRoleFromGroup(ctx, groupID, d.Id())
	err = utils.SuppressErrorOn404(resp, err)
	if err != nil {
		return diag.Errorf("failed to remove role %s assigned to group %s: %v", roleType, groupID, err)
	}
	return nil
}

func listGroupTargetsIDs(ctx context.Context, meta interface{}, groupID, roleID string) ([]string, error) {
	var resIDs []string
	targets, resp, err := getOktaClientFromMetadata(meta).Group.ListGroupTargetsForGroupRole(ctx, groupID, roleID, &query.Params{Limit: utils.DefaultPaginationLimit})
	if err != nil {
		return nil, fmt.Errorf("unable to get admin assignment %s for group %s: %v", roleID, groupID, err)
	}
	for _, target := range targets {
		resIDs = append(resIDs, target.Id)
	}
	for resp.HasNextPage() {
		var additionalTargets []*sdk.Group
		resp, err = resp.Next(ctx, &additionalTargets)
		if err != nil {
			return nil, err
		}
		for _, target := range additionalTargets {
			resIDs = append(resIDs, target.Id)
		}
	}
	return resIDs, nil
}

func listGroupAppsTargets(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]string, error) {
	var resApps []string
	apps, resp, err := getOktaClientFromMetadata(meta).Group.
		ListApplicationTargetsForApplicationAdministratorRoleForGroup(
			ctx, d.Get("group_id").(string), d.Id(), &query.Params{Limit: utils.DefaultPaginationLimit, Status: "ACTIVE"})
	if err != nil {
		return nil, err
	}
	for {
		for _, app := range apps {
			if app.Id != "" {
				a := sdk.NewApplication()
				_, resp, err := getOktaClientFromMetadata(meta).Application.GetApplication(ctx, app.Id, a, nil)
				if err := utils.SuppressErrorOn404(resp, err); err != nil {
					return nil, err
				}
				if a.Name == "" {
					return nil, fmt.Errorf("something is wrong here, becase application was in the group target list just a second ago, and now it does not exist")
				}
				resApps = append(resApps, fmt.Sprintf("%s.%s", a.Name, a.Id))
			} else {
				resApps = append(resApps, app.Name)
			}
		}
		if resp.HasNextPage() {
			resp, err = resp.Next(ctx, &apps)
			if err != nil {
				return nil, err
			}
			continue
		} else {
			break
		}
	}
	return resApps, nil
}

func addGroupTargetsToRole(ctx context.Context, client *sdk.Client, groupID, roleID string, groupTargets []string) error {
	for i := range groupTargets {
		_, err := client.Group.AddGroupTargetToGroupAdministratorRoleForGroup(ctx, groupID, roleID, groupTargets[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func removeGroupTargetsFromRole(ctx context.Context, client *sdk.Client, groupID, roleID string, groupTargets []string) error {
	for i := range groupTargets {
		resp, err := client.Group.RemoveGroupTargetFromGroupAdministratorRoleGivenToGroup(ctx, groupID, roleID, groupTargets[i])
		err = utils.SuppressErrorOn404(resp, err)
		if err != nil {
			return err
		}
	}
	return nil
}

func addGroupAppTargetsToRole(ctx context.Context, client *sdk.Client, groupID, roleID string, apps []string) error {
	for i := range apps {
		app := strings.Split(apps[i], ".")
		if len(app) == 1 {
			_, err := client.Group.AddApplicationTargetToAdminRoleGivenToGroup(ctx,
				groupID, roleID, app[0])
			if err != nil {
				return fmt.Errorf("failed to add an app target to an app administrator role given to a group: %v", err)
			}
		} else {
			_, err := client.Group.AddApplicationInstanceTargetToAppAdminRoleGivenToGroup(ctx,
				groupID, roleID, app[0], strings.Join(app[1:], ""))
			if err != nil {
				return fmt.Errorf("failed to add an app instance target to an app administrator role given to a group: %v", err)
			}
		}
	}
	return nil
}

func removeGroupAppTargets(ctx context.Context, client *sdk.Client, groupID, roleID string, apps []string) error {
	for i := range apps {
		app := strings.Split(apps[i], ".")
		if len(app) == 1 {
			resp, err := client.Group.RemoveApplicationTargetFromApplicationAdministratorRoleGivenToGroup(ctx,
				groupID, roleID, app[0])
			err = utils.SuppressErrorOn404(resp, err)
			if err != nil {
				return fmt.Errorf("failed to remove an app target from an app administrator role given to a group: %v", err)
			}
		} else {
			resp, err := client.Group.RemoveApplicationTargetFromAdministratorRoleGivenToGroup(ctx,
				groupID, roleID, app[0], strings.Join(app[1:], ""))
			err = utils.SuppressErrorOn404(resp, err)
			if err != nil {
				return fmt.Errorf("failed to remove an app instance target from an app administrator role given to a group: %v", err)
			}
		}
	}
	return nil
}

func findRole(ctx context.Context, d *schema.ResourceData, meta interface{}) (*sdk.Role, error) {
	rt := d.Get("role_type").(string)
	if rt == "" {
		return nil, nil
	}
	groupID := d.Get("group_id").(string)
	roles, resp, err := getOktaClientFromMetadata(meta).Group.ListGroupAssignedRoles(ctx, groupID, nil)
	if err := utils.SuppressErrorOn404(resp, err); err != nil {
		return nil, fmt.Errorf("failed to list roles assigned to group %s: %v", groupID, err)
	}
	if len(roles) == 0 {
		logger(meta).Error("either group (%s) which had these roles assigned no longer exists or no roles were assigned", groupID)
		return nil, nil
	}
	for i := range roles {
		if roles[i].Type == rt {
			d.SetId(roles[i].Id)
			return roles[i], nil
		}
	}
	return nil, nil
}

func supportsGroupTargets(roleType string) bool {
	return utils.Contains([]string{"GROUP_MEMBERSHIP_ADMIN", "HELP_DESK_ADMIN", "USER_ADMIN"}, roleType)
}
