package idaas

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/api"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/okta/terraform-provider-okta/sdk/query"
)

var rolesWithTargets = []string{"APP_ADMIN", "GROUP_MEMBERSHIP_ADMIN", "HELP_DESK_ADMIN", "USER_ADMIN"}

func resourceAdminRoleTargets() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAdminRoleCreate,
		ReadContext:   resourceAdminRoleRead,
		UpdateContext: resourceAdminRoleUpdate,
		DeleteContext: resourceAdminRoleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				parts := strings.Split(d.Id(), "/")
				if len(parts) != 2 {
					return nil, errors.New("invalid resource import specifier. Use: terraform import <user_id>/<role_type>")
				}
				if !utils.Contains(rolesWithTargets, parts[1]) {
					return nil, fmt.Errorf("invalid role type, use one of %s", strings.Join(rolesWithTargets, ","))
				}
				err := checkRoleAssignment(ctx, d, meta, parts[0], parts[1])
				if err != nil {
					return nil, err
				}
				_ = d.Set("user_id", parts[0])
				_ = d.Set("role_type", parts[1])
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "User associated with the role",
				ForceNew:    true,
			},
			"role_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of the role that is assigned to the user and supports optional targets. See [API Docs](https://developer.okta.com/docs/api/openapi/okta-management/guides/roles/#standard-roles)",
				ForceNew:    true,
			},
			"role_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of a role",
			},
			"apps": {
				Type:          schema.TypeSet,
				Optional:      true,
				Description:   "List of app names (name represents set of app instances) or a combination of app name and app instance ID (like 'salesforce' or 'facebook.0oapsqQ6dv19pqyEo0g3')",
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"groups"},
			},
			"groups": {
				Type:          schema.TypeSet,
				Optional:      true,
				Description:   "List of group IDs. Conflicts with apps",
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"apps"},
			},
		},
	}
}

func resourceAdminRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logger(meta).Info("creating admin role targets", "role", d.Get("role_type").(string), "user", d.Get("user_id").(string))
	err := checkRoleAssignment(ctx, d, meta, d.Get("user_id").(string), d.Get("role_type").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	if d.Get("role_type").(string) == "APP_ADMIN" {
		err = addUserAppTargets(ctx, d, meta, utils.ConvertInterfaceToStringSet(d.Get("apps")))
	} else {
		err = addUserGroupTargets(ctx, d, meta, utils.ConvertInterfaceToStringSet(d.Get("groups")))
	}
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%s/%s", d.Get("user_id").(string), d.Get("role_type").(string)))
	return resourceAdminRoleRead(ctx, d, meta)
}

func resourceAdminRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logger(meta).Info("reading admin role targets", "role", d.Get("role_type").(string), "user", d.Get("user_id").(string))
	role, resp, err := getOktaClientFromMetadata(meta).User.GetUserRole(ctx, d.Get("user_id").(string), d.Get("role_id").(string))
	if err := utils.SuppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get role assigned to a user: %v", err)
	}
	if role == nil {
		d.SetId("")
		return nil
	}
	if d.Get("role_type").(string) == "APP_ADMIN" {
		apps, err := listUserApplicationTargets(ctx, d, meta)
		if err != nil {
			return diag.Errorf("failed to read app targets: %v", err)
		}
		_ = d.Set("apps", utils.ConvertStringSliceToSet(apps))
	} else {
		groups, err := listUserGroupTargets(ctx, d, meta)
		if err != nil {
			return diag.Errorf("failed to read group targets: %v", err)
		}
		_ = d.Set("groups", utils.ConvertStringSliceToSet(groups))
	}
	return nil
}

func resourceAdminRoleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logger(meta).Info("updating admin role targets", "role", d.Get("role_type").(string), "user", d.Get("user_id").(string))
	if d.Get("role_type").(string) == "APP_ADMIN" {
		expectedApps := utils.ConvertInterfaceToStringSet(d.Get("apps"))
		if len(expectedApps) == 0 {
			err := clearAndRefreshRole(ctx, d, meta)
			if err != nil {
				return diag.FromErr(err)
			}
			return nil
		}
		existingApps, err := listUserApplicationTargets(ctx, d, meta)
		if err != nil {
			return diag.Errorf("failed to update app targets: %v", err)
		}
		appsToAdd, appsToRemove := splitTargets(expectedApps, existingApps)
		err = addUserAppTargets(ctx, d, meta, appsToAdd)
		if err != nil {
			return diag.Errorf("failed to update app targets: %v", err)
		}
		err = removeUserAppTargets(ctx, d, meta, appsToRemove)
		if err != nil {
			return diag.Errorf("failed to update app targets: %v", err)
		}
	} else {
		expectedGroups := utils.ConvertInterfaceToStringSet(d.Get("groups"))
		if len(expectedGroups) == 0 {
			err := clearAndRefreshRole(ctx, d, meta)
			if err != nil {
				return diag.FromErr(err)
			}
			return nil
		}
		existingGroups, err := listUserGroupTargets(ctx, d, meta)
		if err != nil {
			return diag.Errorf("failed to update group targets: %v", err)
		}
		groupsToAdd, groupsToRemove := splitTargets(expectedGroups, existingGroups)
		err = addUserGroupTargets(ctx, d, meta, groupsToAdd)
		if err != nil {
			return diag.Errorf("failed to update group targets: %v", err)
		}
		err = removeUserGroupTargets(ctx, d, meta, groupsToRemove)
		if err != nil {
			return diag.Errorf("failed to update group targets: %v", err)
		}
	}
	return resourceAdminRoleRead(ctx, d, meta)
}

func resourceAdminRoleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logger(meta).Info("removing admin role targets", "role", d.Get("role_type").(string), "user", d.Get("user_id").(string))
	_, err := removeAllTargets(ctx, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func splitTargets(expectedApps, existingApps []string) (appsToAdd, appsToRemove []string) {
	for i := range expectedApps {
		if !utils.Contains(existingApps, expectedApps[i]) {
			appsToAdd = append(appsToAdd, expectedApps[i])
		}
	}
	for i := range existingApps {
		if !utils.Contains(expectedApps, existingApps[i]) {
			appsToRemove = append(appsToRemove, existingApps[i])
		}
	}
	return
}

func clearAndRefreshRole(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	newRoleID, err := removeAllTargets(ctx, d, meta)
	if err != nil {
		return err
	}
	_ = d.Set("role_id", newRoleID)
	return nil
}

// to avoid exception when removing last target group or app,
// the API consumer should delete the role assignment and recreate it.
func removeAllTargets(ctx context.Context, d *schema.ResourceData, meta interface{}) (string, error) {
	resp, err := getOktaClientFromMetadata(meta).User.RemoveRoleFromUser(ctx, d.Get("user_id").(string), d.Get("role_id").(string))
	if err := utils.SuppressErrorOn404(resp, err); err != nil {
		return "", fmt.Errorf("failed to unassign '%s' role from user: %v", d.Get("role_type").(string), err)
	}
	ctx = context.WithValue(ctx, api.RetryOnStatusCodes, []int{http.StatusConflict, http.StatusBadRequest})
	role, _, err := getOktaClientFromMetadata(meta).User.AssignRoleToUser(ctx, d.Get("user_id").(string),
		sdk.AssignRoleRequest{Type: d.Get("role_type").(string)}, nil)
	if err != nil {
		d.SetId("")
		return "", fmt.Errorf("failed to assign '%s' role back to user: %v", d.Get("role_type").(string), err)
	}
	if role == nil {
		return "", errors.New("role was nil")
	}
	return role.Id, nil
}

func addUserAppTargets(ctx context.Context, d *schema.ResourceData, meta interface{}, apps []string) error {
	for i := range apps {
		app := strings.Split(apps[i], ".")
		if len(app) == 1 {
			_, err := getOktaClientFromMetadata(meta).User.AddApplicationTargetToAdminRoleForUser(ctx,
				d.Get("user_id").(string), d.Get("role_id").(string), app[0])
			if err != nil {
				return fmt.Errorf("failed to add an app target to an app administrator role given to a user: %v", err)
			}
		} else {
			_, err := getOktaClientFromMetadata(meta).User.AddApplicationTargetToAppAdminRoleForUser(ctx,
				d.Get("user_id").(string), d.Get("role_id").(string), app[0], strings.Join(app[1:], ""))
			if err != nil {
				return fmt.Errorf("failed to add an app instance target to an app administrator role given to a user: %v", err)
			}
		}
	}
	return nil
}

func removeUserAppTargets(ctx context.Context, d *schema.ResourceData, meta interface{}, apps []string) error {
	for i := range apps {
		app := strings.Split(apps[i], ".")
		if len(app) == 1 {
			_, err := getOktaClientFromMetadata(meta).User.RemoveApplicationTargetFromApplicationAdministratorRoleForUser(ctx,
				d.Get("user_id").(string), d.Get("role_id").(string), app[0])
			if err != nil {
				return fmt.Errorf("failed to remove an app target from an app administrator role given to a user: %v", err)
			}
		} else {
			_, err := getOktaClientFromMetadata(meta).User.RemoveApplicationTargetFromAdministratorRoleForUser(ctx,
				d.Get("user_id").(string), d.Get("role_id").(string), app[0], strings.Join(app[1:], ""))
			if err != nil {
				return fmt.Errorf("failed to remove an app instance target from an app administrator role given to a user: %v", err)
			}
		}
	}
	return nil
}

func addUserGroupTargets(ctx context.Context, d *schema.ResourceData, meta interface{}, groups []string) error {
	for i := range groups {
		_, err := getOktaClientFromMetadata(meta).User.AddGroupTargetToRole(ctx,
			d.Get("user_id").(string), d.Get("role_id").(string), groups[i])
		if err != nil {
			return fmt.Errorf("failed to add a group target to a group administrator role given to a user: %v", err)
		}
	}
	return nil
}

func removeUserGroupTargets(ctx context.Context, d *schema.ResourceData, meta interface{}, groups []string) error {
	for i := range groups {
		_, err := getOktaClientFromMetadata(meta).User.RemoveGroupTargetFromRole(ctx,
			d.Get("user_id").(string), d.Get("role_id").(string), groups[i])
		if err != nil {
			return fmt.Errorf("failed to add a group target to a group administrator role given to a user: %v", err)
		}
	}
	return nil
}

func listUserGroupTargets(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]string, error) {
	var resGroups []string
	groups, resp, err := getOktaClientFromMetadata(meta).User.
		ListGroupTargetsForRole(ctx, d.Get("user_id").(string), d.Get("role_id").(string), &query.Params{Limit: utils.DefaultPaginationLimit})
	if err != nil {
		return nil, err
	}
	for {
		for _, group := range groups {
			resGroups = append(resGroups, group.Id)
		}
		if resp.HasNextPage() {
			resp, err = resp.Next(ctx, &groups)
			if err != nil {
				return nil, err
			}
			continue
		} else {
			break
		}
	}
	return resGroups, nil
}

func listUserApplicationTargets(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]string, error) {
	var resApps []string
	apps, resp, err := getOktaClientFromMetadata(meta).User.
		ListApplicationTargetsForApplicationAdministratorRoleForUser(
			ctx, d.Get("user_id").(string), d.Get("role_id").(string), &query.Params{Limit: utils.DefaultPaginationLimit})
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
					return nil, fmt.Errorf("something is wrong here, becase application was in the target list just a second ago, and now it's gone")
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

func checkRoleAssignment(ctx context.Context, d *schema.ResourceData, meta interface{}, userID, roleType string) error {
	roles, _, err := getOktaClientFromMetadata(meta).User.ListAssignedRolesForUser(ctx, userID, nil)
	if err != nil {
		return fmt.Errorf("failed to get list of roles associated with the user: %v", err)
	}
	for i := range roles {
		if roles[i].Type == roleType {
			_ = d.Set("role_id", roles[i].Id)
			break
		}
	}
	if d.Get("role_id").(string) == "" {
		return fmt.Errorf("please assign '%s' to a user before creating or importing this resource", roleType)
	}
	return nil
}
