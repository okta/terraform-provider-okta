package okta

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

func resourceGroupRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupRoleCreate,
		ReadContext:   resourceGroupRoleRead,
		UpdateContext: resourceGroupRoleUpdate,
		DeleteContext: resourceGroupRoleDelete,
		Importer:      createNestedResourceImporter([]string{"group_id", "id"}),
		CustomizeDiff: customdiff.All(
			customdiff.ForceNewIf("target_group_list", func(_ context.Context, d *schema.ResourceDiff, m interface{}) bool {
				if d.HasChange("target_group_list") {
					// to avoid exception when removing last group target from a role assignment,
					// the API consumer should delete the role assignment and recreate it.
					if len(convertInterfaceToStringSet(d.Get("target_group_list"))) == 0 {
						return true
					}
				}
				return false
			}),
			customdiff.ForceNewIf("target_app_list", func(_ context.Context, d *schema.ResourceDiff, m interface{}) bool {
				if d.HasChange("target_app_list") {
					// to avoid exception when removing last app target from a role assignment,
					// the API consumer should delete the role assignment and recreate it.
					if len(convertInterfaceToStringSet(d.Get("target_app_list"))) == 0 {
						return true
					}
					// deleting the app from an org which is the only one left assigned to this group role
					// will cause the role to be unassigned from the group thus resource should be recreated
					oldValue, newValue := d.GetChange("target_app_list")
					oldApps := convertInterfaceToStringSet(oldValue)
					newApps := convertInterfaceToStringSet(newValue)
					if len(oldApps) > 0 && !containsOne(oldApps, newApps...) {
						return true
					}
				}
				return false
			}),
		),
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of group to attach admin roles to",
				ForceNew:    true,
			},
			"role_type": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Type of Role to assign",
				ForceNew:         true,
				ValidateDiagFunc: elemInSlice(validAdminRoles),
			},
			"target_group_list": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "List of groups ids for the targets of the admin role.",
			},
			"target_app_list": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "List of apps ids for the targets of the admin role.",
			},
			"disable_notifications": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "When this setting is enabled, the admins won't receive any of the default Okta administrator emails",
				Default:     false,
			},
		},
	}
}

func resourceGroupRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	groupID := d.Get("group_id").(string)
	roleType := d.Get("role_type").(string)
	client := getOktaClientFromMetadata(m)
	logger(m).Info("assigning role to group", "group_id", groupID, "role_type", roleType)
	role, _, err := client.Group.AssignRoleToGroup(ctx, groupID, okta.AssignRoleRequest{Type: roleType},
		&query.Params{DisableNotifications: boolPtr(d.Get("disable_notifications").(bool))})
	if err != nil {
		return diag.Errorf("failed to assign role %s to group %s: %v", roleType, groupID, err)
	}
	groupTargets := convertInterfaceToStringSet(d.Get("target_group_list"))
	if len(groupTargets) > 0 && supportsGroupTargets(roleType) {
		logger(m).Info("scoping admin role assignment to list of groups", "group_id", groupID, "role_id", role.Id, "target_group_list", groupTargets)
		err = addGroupTargetsToRole(ctx, client, groupID, role.Id, groupTargets)
		if err != nil {
			return diag.Errorf("unable to add group target to role assignment %s for group %s: %v", role.Id, groupID, err)
		}
	}
	appTargets := convertInterfaceToStringSet(d.Get("target_app_list"))
	if len(appTargets) > 0 && roleType == "APP_ADMIN" {
		logger(m).Info("scoping admin role assignment to list of apps", "group_id", groupID, "role_id", role.Id, "target_app_list", appTargets)
		err = addGroupAppTargetsToRole(ctx, client, groupID, role.Id, appTargets)
		if err != nil {
			return diag.Errorf("unable to add app targets to role assignment %s for group %s: %v", role.Id, groupID, err)
		}
	}
	d.SetId(role.Id)
	bOff := backoff.NewExponentialBackOff()
	bOff.MaxElapsedTime = time.Second * 10
	bOff.InitialInterval = time.Second
	err = backoff.Retry(func() error {
		err := resourceGroupRoleRead(ctx, d, m)
		if err != nil {
			return backoff.Permanent(fmt.Errorf("%s", err[0].Summary))
		}
		if d.Id() != "" {
			return nil
		}
		return fmt.Errorf("role %s was not assigned to a group %s", roleType, groupID)
	}, bOff)
	return diag.FromErr(err)
}

func resourceGroupRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	groupID := d.Get("group_id").(string)
	role, resp, err := getOktaClientFromMetadata(m).Group.GetRole(ctx, groupID, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get role '%s' assigned to group '%s': %v", d.Id(), groupID, err)
	}
	if role == nil {
		role, err = findRole(ctx, d, m)
		if err != nil {
			return diag.Errorf("failed to get role '%s' assigned to group '%s': %v", d.Id(), groupID, err)
		}
		if role == nil {
			d.SetId("")
			return nil
		}
	}
	if supportsGroupTargets(role.Type) {
		groupIDs, err := listGroupTargetsIDs(ctx, m, groupID, role.Id)
		if err != nil {
			return diag.Errorf("unable to list group targets for role %s and group %s: %v", role.Id, groupID, err)
		}
		if len(groupIDs) != 0 {
			_ = d.Set("target_group_list", groupIDs)
		}
	} else if role.Type == "APP_ADMIN" {
		apps, err := listGroupAppsTargets(ctx, d, m)
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

func resourceGroupRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	groupID := d.Get("group_id").(string)
	roleID := d.Id()
	roleType := d.Get("role_type").(string)
	client := getOktaClientFromMetadata(m)
	if d.HasChange("disable_notifications") {
		_, _, err := client.Group.AssignRoleToGroup(ctx, groupID, okta.AssignRoleRequest{},
			&query.Params{DisableNotifications: boolPtr(d.Get("disable_notifications").(bool))})
		if err != nil {
			return diag.Errorf("failed to update group's '%s' notification settings: %v", groupID, err)
		}
	}
	if d.HasChange("target_group_list") && supportsGroupTargets(roleType) {
		expectedGroupIDs := convertInterfaceToStringSet(d.Get("target_group_list"))
		existingGroupIDs, err := listGroupTargetsIDs(ctx, m, groupID, roleID)
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
		expectedApps := convertInterfaceToStringSet(d.Get("target_app_list"))
		existingApps, err := listGroupAppsTargets(ctx, d, m)
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
	return resourceGroupRoleRead(ctx, d, m)
}

func resourceGroupRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	groupID := d.Get("group_id").(string)
	roleType := d.Get("role_type").(string)
	logger(m).Info("deleting assigned role from group", "group_id", groupID, "role_type", roleType)
	resp, err := getOktaClientFromMetadata(m).Group.RemoveRoleFromGroup(ctx, groupID, d.Id())
	err = suppressErrorOn404(resp, err)
	if err != nil {
		return diag.Errorf("failed to remove role %s assigned to group %s: %v", roleType, groupID, err)
	}
	return nil
}

func listGroupTargetsIDs(ctx context.Context, m interface{}, groupID, roleID string) ([]string, error) {
	var resIDs []string
	targets, resp, err := getOktaClientFromMetadata(m).Group.ListGroupTargetsForGroupRole(ctx, groupID, roleID, &query.Params{Limit: defaultPaginationLimit})
	if err != nil {
		return nil, fmt.Errorf("unable to get admin assignment %s for group %s: %v", roleID, groupID, err)
	}
	for _, target := range targets {
		resIDs = append(resIDs, target.Id)
	}
	for resp.HasNextPage() {
		var additionalTargets []*okta.Group
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

func listGroupAppsTargets(ctx context.Context, d *schema.ResourceData, m interface{}) ([]string, error) {
	var resApps []string
	apps, resp, err := getOktaClientFromMetadata(m).Group.
		ListApplicationTargetsForApplicationAdministratorRoleForGroup(
			ctx, d.Get("group_id").(string), d.Id(), &query.Params{Limit: defaultPaginationLimit, Status: "ACTIVE"})
	if err != nil {
		return nil, err
	}
	for {
		for _, app := range apps {
			if app.Id != "" {
				a := okta.NewApplication()
				_, resp, err := getOktaClientFromMetadata(m).Application.GetApplication(ctx, app.Id, a, nil)
				if err := suppressErrorOn404(resp, err); err != nil {
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

func addGroupTargetsToRole(ctx context.Context, client *okta.Client, groupID, roleID string, groupTargets []string) error {
	for i := range groupTargets {
		_, err := client.Group.AddGroupTargetToGroupAdministratorRoleForGroup(ctx, groupID, roleID, groupTargets[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func removeGroupTargetsFromRole(ctx context.Context, client *okta.Client, groupID, roleID string, groupTargets []string) error {
	for i := range groupTargets {
		resp, err := client.Group.RemoveGroupTargetFromGroupAdministratorRoleGivenToGroup(ctx, groupID, roleID, groupTargets[i])
		err = suppressErrorOn404(resp, err)
		if err != nil {
			return err
		}
	}
	return nil
}

func addGroupAppTargetsToRole(ctx context.Context, client *okta.Client, groupID, roleID string, apps []string) error {
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

func removeGroupAppTargets(ctx context.Context, client *okta.Client, groupID, roleID string, apps []string) error {
	for i := range apps {
		app := strings.Split(apps[i], ".")
		if len(app) == 1 {
			resp, err := client.Group.RemoveApplicationTargetFromApplicationAdministratorRoleGivenToGroup(ctx,
				groupID, roleID, app[0])
			err = suppressErrorOn404(resp, err)
			if err != nil {
				return fmt.Errorf("failed to remove an app target from an app administrator role given to a group: %v", err)
			}
		} else {
			resp, err := client.Group.RemoveApplicationTargetFromAdministratorRoleGivenToGroup(ctx,
				groupID, roleID, app[0], strings.Join(app[1:], ""))
			err = suppressErrorOn404(resp, err)
			if err != nil {
				return fmt.Errorf("failed to remove an app instance target from an app administrator role given to a group: %v", err)
			}
		}
	}
	return nil
}

func findRole(ctx context.Context, d *schema.ResourceData, m interface{}) (*okta.Role, error) {
	rt := d.Get("role_type").(string)
	if rt == "" {
		return nil, nil
	}
	groupID := d.Get("group_id").(string)
	roles, resp, err := getOktaClientFromMetadata(m).Group.ListGroupAssignedRoles(ctx, groupID, nil)
	if err := suppressErrorOn404(resp, err); err != nil {
		return nil, fmt.Errorf("failed to list roles assigned to group %s: %v", groupID, err)
	}
	if len(roles) == 0 {
		logger(m).Error("either group (%s) which had these roles assigned no longer exists or no roles were assigned", groupID)
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
	return contains([]string{"GROUP_MEMBERSHIP_ADMIN", "HELP_DESK_ADMIN", "USER_ADMIN"}, roleType)
}
