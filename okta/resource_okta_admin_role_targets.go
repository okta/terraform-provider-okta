package okta

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

var rolesWithTargets = []string{"APP_ADMIN", "GROUP_MEMBERSHIP_ADMIN", "HELP_DESK_ADMIN", "USER_ADMIN"}

func resourceAdminRoleTargets() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAdminRoleCreate,
		ReadContext:   resourceAdminRoleRead,
		UpdateContext: resourceAdminRoleUpdate,
		DeleteContext: resourceAdminRoleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				parts := strings.Split(d.Id(), "/")
				if len(parts) != 2 {
					return nil, errors.New("invalid resource import specifier. Use: terraform import <user_id>/<role_type>")
				}
				if !contains(rolesWithTargets, parts[1]) {
					return nil, fmt.Errorf("invalid role type, use one of %s", strings.Join(rolesWithTargets, ","))
				}
				err := checkRoleAssignment(ctx, d, m, parts[0], parts[1])
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
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Type of the role that is assigned to the user and supports optional targets",
				ForceNew:         true,
				ValidateDiagFunc: elemInSlice(rolesWithTargets),
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
				Description:   "List of group IDs",
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"apps"},
			},
		},
	}
}

func resourceAdminRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	logger(m).Info("creating admin role targets", "role", d.Get("role_type").(string), "user", d.Get("user_id").(string))
	err := checkRoleAssignment(ctx, d, m, d.Get("user_id").(string), d.Get("role_type").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	if d.Get("role_type").(string) == "APP_ADMIN" {
		err = addUserAppTargets(ctx, d, m, convertInterfaceToStringSet(d.Get("apps")))
	} else {
		err = addUserGroupTargets(ctx, d, m, convertInterfaceToStringSet(d.Get("groups")))
	}
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%s/%s", d.Get("user_id").(string), d.Get("role_type").(string)))
	return resourceAdminRoleRead(ctx, d, m)
}

func resourceAdminRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	logger(m).Info("reading admin role targets", "role", d.Get("role_type").(string), "user", d.Get("user_id").(string))
	role, resp, err := getOktaClientFromMetadata(m).User.GetUserRole(ctx, d.Get("user_id").(string), d.Get("role_id").(string))
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get role assigned to a user: %v", err)
	}
	if role == nil {
		d.SetId("")
		return nil
	}
	if d.Get("role_type").(string) == "APP_ADMIN" {
		apps, err := listUserApplicationTargets(ctx, d, m)
		if err != nil {
			return diag.Errorf("failed to read app targets: %v", err)
		}
		_ = d.Set("apps", convertStringSliceToSet(apps))
	} else {
		groups, err := listUserGroupTargets(ctx, d, m)
		if err != nil {
			return diag.Errorf("failed to read group targets: %v", err)
		}
		_ = d.Set("groups", convertStringSliceToSet(groups))
	}
	return nil
}

func resourceAdminRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	logger(m).Info("updating admin role targets", "role", d.Get("role_type").(string), "user", d.Get("user_id").(string))
	if d.Get("role_type").(string) == "APP_ADMIN" {
		expectedApps := convertInterfaceToStringSet(d.Get("apps"))
		if len(expectedApps) == 0 {
			err := clearAndRefreshRole(ctx, d, m)
			if err != nil {
				return diag.FromErr(err)
			}
			return nil
		}
		existingApps, err := listUserApplicationTargets(ctx, d, m)
		if err != nil {
			return diag.Errorf("failed to update app targets: %v", err)
		}
		appsToAdd, appsToRemove := splitTargets(expectedApps, existingApps)
		err = addUserAppTargets(ctx, d, m, appsToAdd)
		if err != nil {
			return diag.Errorf("failed to update app targets: %v", err)
		}
		err = removeUserAppTargets(ctx, d, m, appsToRemove)
		if err != nil {
			return diag.Errorf("failed to update app targets: %v", err)
		}
	} else {
		expectedGroups := convertInterfaceToStringSet(d.Get("groups"))
		if len(expectedGroups) == 0 {
			err := clearAndRefreshRole(ctx, d, m)
			if err != nil {
				return diag.FromErr(err)
			}
			return nil
		}
		existingGroups, err := listUserGroupTargets(ctx, d, m)
		if err != nil {
			return diag.Errorf("failed to update group targets: %v", err)
		}
		groupsToAdd, groupsToRemove := splitTargets(expectedGroups, existingGroups)
		err = addUserGroupTargets(ctx, d, m, groupsToAdd)
		if err != nil {
			return diag.Errorf("failed to update group targets: %v", err)
		}
		err = removeUserGroupTargets(ctx, d, m, groupsToRemove)
		if err != nil {
			return diag.Errorf("failed to update group targets: %v", err)
		}
	}
	return resourceAdminRoleRead(ctx, d, m)
}

func resourceAdminRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	logger(m).Info("removing admin role targets", "role", d.Get("role_type").(string), "user", d.Get("user_id").(string))
	_, err := removeAllTargets(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func splitTargets(expectedApps, existingApps []string) (appsToAdd, appsToRemove []string) {
	for i := range expectedApps {
		if !contains(existingApps, expectedApps[i]) {
			appsToAdd = append(appsToAdd, expectedApps[i])
		}
	}
	for i := range existingApps {
		if !contains(expectedApps, existingApps[i]) {
			appsToRemove = append(appsToRemove, existingApps[i])
		}
	}
	return
}

func clearAndRefreshRole(ctx context.Context, d *schema.ResourceData, m interface{}) error {
	newRoleID, err := removeAllTargets(ctx, d, m)
	if err != nil {
		return err
	}
	_ = d.Set("role_id", newRoleID)
	return nil
}

// to avoid exception when removing last target group or app,
// the API consumer should delete the role assignment and recreate it.
func removeAllTargets(ctx context.Context, d *schema.ResourceData, m interface{}) (string, error) {
	resp, err := getOktaClientFromMetadata(m).User.RemoveRoleFromUser(ctx, d.Get("user_id").(string), d.Get("role_id").(string))
	if err := suppressErrorOn404(resp, err); err != nil {
		return "", fmt.Errorf("failed to unassign '%s' role from user: %v", d.Get("role_type").(string), err)
	}
	ctx = context.WithValue(ctx, retryOnStatusCodes, []int{http.StatusConflict, http.StatusBadRequest})
	role, _, err := getOktaClientFromMetadata(m).User.AssignRoleToUser(ctx, d.Get("user_id").(string),
		okta.AssignRoleRequest{Type: d.Get("role_type").(string)}, nil)
	if err != nil {
		d.SetId("")
		return "", fmt.Errorf("failed to assign '%s' role back to user: %v", d.Get("role_type").(string), err)
	}
	return role.Id, nil
}

func addUserAppTargets(ctx context.Context, d *schema.ResourceData, m interface{}, apps []string) error {
	for i := range apps {
		app := strings.Split(apps[i], ".")
		if len(app) == 1 {
			_, err := getOktaClientFromMetadata(m).User.AddApplicationTargetToAdminRoleForUser(ctx,
				d.Get("user_id").(string), d.Get("role_id").(string), app[0])
			if err != nil {
				return fmt.Errorf("failed to add an app target to an app administrator role given to a user: %v", err)
			}
		} else {
			_, err := getOktaClientFromMetadata(m).User.AddApplicationTargetToAppAdminRoleForUser(ctx,
				d.Get("user_id").(string), d.Get("role_id").(string), app[0], strings.Join(app[1:], ""))
			if err != nil {
				return fmt.Errorf("failed to add an app instance target to an app administrator role given to a user: %v", err)
			}
		}
	}
	return nil
}

func removeUserAppTargets(ctx context.Context, d *schema.ResourceData, m interface{}, apps []string) error {
	for i := range apps {
		app := strings.Split(apps[i], ".")
		if len(app) == 1 {
			_, err := getOktaClientFromMetadata(m).User.RemoveApplicationTargetFromApplicationAdministratorRoleForUser(ctx,
				d.Get("user_id").(string), d.Get("role_id").(string), app[0])
			if err != nil {
				return fmt.Errorf("failed to remove an app target from an app administrator role given to a user: %v", err)
			}
		} else {
			_, err := getOktaClientFromMetadata(m).User.RemoveApplicationTargetFromAdministratorRoleForUser(ctx,
				d.Get("user_id").(string), d.Get("role_id").(string), app[0], strings.Join(app[1:], ""))
			if err != nil {
				return fmt.Errorf("failed to remove an app instance target from an app administrator role given to a user: %v", err)
			}
		}
	}
	return nil
}

func addUserGroupTargets(ctx context.Context, d *schema.ResourceData, m interface{}, groups []string) error {
	for i := range groups {
		_, err := getOktaClientFromMetadata(m).User.AddGroupTargetToRole(ctx,
			d.Get("user_id").(string), d.Get("role_id").(string), groups[i])
		if err != nil {
			return fmt.Errorf("failed to add a group target to a group administrator role given to a user: %v", err)
		}
	}
	return nil
}

func removeUserGroupTargets(ctx context.Context, d *schema.ResourceData, m interface{}, groups []string) error {
	for i := range groups {
		_, err := getOktaClientFromMetadata(m).User.RemoveGroupTargetFromRole(ctx,
			d.Get("user_id").(string), d.Get("role_id").(string), groups[i])
		if err != nil {
			return fmt.Errorf("failed to add a group target to a group administrator role given to a user: %v", err)
		}
	}
	return nil
}

func listUserGroupTargets(ctx context.Context, d *schema.ResourceData, m interface{}) ([]string, error) {
	var resGroups []string
	groups, resp, err := getOktaClientFromMetadata(m).User.
		ListGroupTargetsForRole(ctx, d.Get("user_id").(string), d.Get("role_id").(string), &query.Params{Limit: defaultPaginationLimit})
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

func listUserApplicationTargets(ctx context.Context, d *schema.ResourceData, m interface{}) ([]string, error) {
	var resApps []string
	apps, resp, err := getOktaClientFromMetadata(m).User.
		ListApplicationTargetsForApplicationAdministratorRoleForUser(
			ctx, d.Get("user_id").(string), d.Get("role_id").(string), &query.Params{Limit: defaultPaginationLimit})
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

func checkRoleAssignment(ctx context.Context, d *schema.ResourceData, m interface{}, userID, roleType string) error {
	roles, _, err := getOktaClientFromMetadata(m).User.ListAssignedRolesForUser(ctx, userID, nil)
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
