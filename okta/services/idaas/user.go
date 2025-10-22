package idaas

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/okta/terraform-provider-okta/sdk/query"
)

const (
	StatusActive   = "ACTIVE"
	StatusInactive = "INACTIVE"

	UserStatusPasswordExpired = "PASSWORD_EXPIRED"
	UserStatusProvisioned     = "PROVISIONED"
	UserStatusDeprovisioned   = "DEPROVISIONED"
	UserStatusStaged          = "STAGED"
	UserStatusSuspended       = "SUSPENDED"
	UserStatusRecovery        = "RECOVERY"
	UserStatusLockedOut       = "LOCKED_OUT"

	UserScope = "USER"

	GroupProfileEveryone = "Everyone"
)

var userProfileDataSchema = map[string]*schema.Schema{
	"admin_roles": {
		Type:     schema.TypeSet,
		Computed: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	},
	"roles": {
		Type:     schema.TypeSet,
		Computed: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	},
	"city": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"cost_center": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"country_code": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"custom_profile_attributes": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"department": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"display_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"division": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"email": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"employee_number": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"first_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"group_memberships": {
		Type:     schema.TypeSet,
		Computed: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	},
	"honorific_prefix": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"honorific_suffix": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"last_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"locale": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"login": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"manager": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"manager_id": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"middle_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"mobile_phone": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"nick_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"organization": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"postal_address": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"preferred_language": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"primary_phone": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"profile_url": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"second_email": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"state": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"status": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"street_address": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"timezone": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"title": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"user_type": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"zip_code": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

func buildUserDataSourceSchema(target map[string]*schema.Schema) map[string]*schema.Schema {
	return utils.BuildSchema(userProfileDataSchema, target)
}

func assignAdminRolesToUser(ctx context.Context, userID string, roles []string, disableNotifications bool, client *sdk.Client) error {
	for _, role := range roles {
		if role == "CUSTOM" {
			continue
		}
		_, _, err := client.User.AssignRoleToUser(ctx, userID, sdk.AssignRoleRequest{Type: role},
			&query.Params{DisableNotifications: utils.BoolPtr(disableNotifications)})
		if err != nil {
			return fmt.Errorf("failed to assign role '%s' to user '%s': %w", role, userID, err)
		}
	}
	return nil
}

func populateUserProfile(d *schema.ResourceData) *sdk.UserProfile {
	profile := sdk.UserProfile{}

	if rawAttrs, ok := d.GetOk("custom_profile_attributes"); ok {
		var attrs map[string]interface{}
		str := rawAttrs.(string)

		// We validate the JSON, no need to check error
		_ = json.Unmarshal([]byte(str), &attrs)
		for k, v := range attrs {
			profile[k] = v
		}
	}

	profile["firstName"] = d.Get("first_name").(string)
	profile["lastName"] = d.Get("last_name").(string)
	profile["login"] = d.Get("login").(string)
	profile["email"] = d.Get("email").(string)

	getSetParams := []string{
		"city", "costCenter", "countryCode", "department", "displayName", "division",
		"employeeNumber", "honorificPrefix", "honorificSuffix", "locale", "manager", "managerId", "middleName",
		"mobilePhone", "nickName", "organization", "preferredLanguage", "primaryPhone", "profileUrl",
		"secondEmail", "state", "streetAddress", "timezone", "title", "userType", "zipCode",
	}

	for i := range getSetParams {
		if res, ok := d.GetOk(utils.CamelCaseToUnderscore(getSetParams[i])); ok {
			profile[getSetParams[i]] = res.(string)
		}
	}

	// need to set profile.postalAddress to nil explicitly if not set because of a bug with this field
	// have a support ticket open with okta about it
	if _, ok := d.GetOk("postal_address"); ok {
		profile["postalAddress"] = d.Get("postal_address").(string)
	} else {
		profile["postalAddress"] = nil
	}

	return &profile
}

func listUserRoles(ctx context.Context, c *sdk.Client, userID string) (userOnlyRoles []*sdk.Role, resp *sdk.Response, err error) {
	roles, resp, err := c.User.ListAssignedRolesForUser(ctx, userID, nil)
	if err != nil {
		return
	}
	userOnlyRoles = append(userOnlyRoles, roles...)
	return
}

func getRoles(ctx context.Context, id string, c *sdk.Client) ([]interface{}, error) {
	roleTypes := make([]interface{}, 0)
	roles, resp, err := listUserRoles(ctx, c, id)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusForbidden {
			// no-op
		} else {
			return nil, err
		}
	} else {
		for _, role := range roles {
			roleTypes = append(roleTypes, role.Type)
		}
	}
	return roleTypes, err
}

func setRoles(ctx context.Context, d *schema.ResourceData, m interface{}) error {
	roleTypes, err := getRoles(ctx, d.Id(), getOktaClientFromMetadata(m))
	if err != nil {
		return fmt.Errorf("failed to get roles: %v", err)
	}
	// set the custom_profile_attributes values
	return utils.SetNonPrimitives(d, map[string]interface{}{
		"roles": schema.NewSet(schema.HashString, roleTypes),
	})
}

func listUserOnlyRoles(ctx context.Context, c *sdk.Client, userID string) (userOnlyRoles []*sdk.Role, resp *sdk.Response, err error) {
	roles, resp, err := c.User.ListAssignedRolesForUser(ctx, userID, nil)
	if err != nil {
		return
	}
	for _, role := range roles {
		if role.AssignmentType == UserScope && role.Type != "CUSTOM" {
			userOnlyRoles = append(userOnlyRoles, role)
		}
	}
	return
}

// set all groups currently attached to the user
func setAllGroups(ctx context.Context, d *schema.ResourceData, c *sdk.Client) error {
	groupIDs, err := getGroupsForUser(ctx, d.Id(), c)
	if err != nil {
		return err
	}
	gids := utils.ConvertStringSliceToInterfaceSlice(groupIDs)
	return utils.SetNonPrimitives(d, map[string]interface{}{
		"group_memberships": schema.NewSet(schema.HashString, gids),
	})
}

func getAdminRoles(ctx context.Context, id string, c *sdk.Client) ([]interface{}, *sdk.Response, error) {
	roleTypes := make([]interface{}, 0)
	roles, resp, err := listUserOnlyRoles(ctx, c, id)

	if err != nil {
		return roleTypes, resp, err
	} else {
		for _, role := range roles {
			roleTypes = append(roleTypes, role.Type)
		}
	}

	return roleTypes, resp, err
}

func setAdminRoles(ctx context.Context, d *schema.ResourceData, m interface{}) error {
	roleTypes, resp, err := getAdminRoles(ctx, d.Id(), getOktaClientFromMetadata(m))
	if err := utils.SuppressErrorOn403("setting admin roles", m, resp, err); err != nil {
		return fmt.Errorf("failed to get admin roles: %v", err)
	}

	// set the custom_profile_attributes values
	return utils.SetNonPrimitives(d, map[string]interface{}{
		"admin_roles": schema.NewSet(schema.HashString, roleTypes),
	})
}

func getGroupsForUser(ctx context.Context, id string, c *sdk.Client) ([]string, error) {
	groups, response, err := c.User.ListUserGroups(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to list user groups: %v", err)
	}

	groupIDs := make([]string, 0)

	for {
		for _, group := range groups {
			groupIDs = append(groupIDs, group.Id)
		}

		if !response.HasNextPage() {
			break
		}

		response, err = response.Next(ctx, &groups)
		if err != nil {
			return nil, fmt.Errorf("failed to list user groups: %v", err)
		}
	}

	return groupIDs, nil
}

func isCustomUserAttr(key string) bool {
	return !utils.Contains(profileKeys, key)
}

func flattenUser(u *sdk.User, filteredCustomAttributes []string) map[string]interface{} {
	customAttributes := make(map[string]interface{})
	attrs := map[string]interface{}{}

	for k, v := range *u.Profile {
		if v != nil {
			attrKey := utils.CamelCaseToUnderscore(k)

			if isCustomUserAttr(attrKey) {

				// Exclude any custom attributes that should be filtered
				if utils.Contains(filteredCustomAttributes, attrKey) {
					continue
				}

				// Supporting any potential type
				ref := reflect.ValueOf(v)
				switch ref.Kind() {
				case reflect.String:
					customAttributes[k] = ref.String()
				case reflect.Float64:
					customAttributes[k] = ref.Float()
				case reflect.Int:
					customAttributes[k] = ref.Int()
				case reflect.Bool:
					customAttributes[k] = ref.Bool()
				case reflect.Slice:
					rawArr := v.([]interface{})
					customAttributes[k] = rawArr
				case reflect.Map:
					rawMap := v.(map[string]interface{})
					customAttributes[k] = rawMap
				}
			} else {
				attrs[attrKey] = v
			}
		}
	}

	attrs["status"] = mapStatus(u.Status)
	if u.RealmId != nil {
		attrs["realm_id"] = u.RealmId
	}

	data, _ := json.Marshal(customAttributes)
	attrs["custom_profile_attributes"] = string(data)

	return attrs
}

// handle setting of user status based on what the current status is because okta
// only allows transitions to certain statuses from other statuses - consult okta User API docs for more info
// https://developer.okta.com/docs/api/resources/users#lifecycle-operations
func updateUserStatus(m interface{}, ctx context.Context, uid, desiredStatus string, c *sdk.Client) error {
	user, _, err := c.User.GetUser(ctx, uid)
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}
	var statusErr error
	switch desiredStatus {
	case UserStatusSuspended:
		_, statusErr = c.User.SuspendUser(ctx, uid)
	case UserStatusDeprovisioned:
		_, statusErr = c.User.DeactivateUser(ctx, uid, nil)
	case StatusActive:
		switch user.Status {
		case UserStatusSuspended:
			_, statusErr = c.User.UnsuspendUser(ctx, uid)
		case UserStatusPasswordExpired:
			// Ignore password expired status. This status is already activated.
			return nil
		case UserStatusLockedOut:
			_, statusErr = c.User.UnlockUser(ctx, uid)
		default:
			_, _, statusErr = c.User.ActivateUser(ctx, uid, nil)
		}
	}
	if statusErr != nil {
		return statusErr
	}
	return waitForStatusTransition(m, ctx, uid, c)
}

// need to wait for user.TransitioningToStatus field to be empty before allowing Terraform to continue
// so the proper current status gets set in the state during the Read operation after a Status update
func waitForStatusTransition(m interface{}, ctx context.Context, u string, c *sdk.Client) error {
	user, _, err := c.User.GetUser(ctx, u)
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}
	for {
		if user.TransitioningToStatus == "" {
			return nil
		}

		log.Printf("[INFO] Transitioning to status = %v; waiting for 5 more seconds...", user.TransitioningToStatus)
		m.(*config.Config).TimeOperations.Sleep(5 * time.Second)
		user, _, err = c.User.GetUser(ctx, u)
		if err != nil {
			return fmt.Errorf("failed to get user: %v", err)
		}
	}
}
