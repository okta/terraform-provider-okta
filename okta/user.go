package okta

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
)

const (
	statusActive   = "ACTIVE"
	statusInactive = "INACTIVE"

	userStatusPasswordExpired = "PASSWORD_EXPIRED"
	userStatusProvisioned     = "PROVISIONED"
	userStatusDeprovisioned   = "DEPROVISIONED"
	userStatusStaged          = "STAGED"
	userStatusSuspended       = "SUSPENDED"
	userStatusRecovery        = "RECOVERY"
	userStatusLockedOut       = "LOCKED_OUT"

	userScope = "USER"

	groupProfileEveryone = "Everyone"
)

var userProfileDataSchema = map[string]*schema.Schema{
	"admin_roles": {
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
		Type: schema.TypeString,
		// StateFunc: normalizeDataJSON,
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
	return buildSchema(userProfileDataSchema, target)
}

func assignAdminRolesToUser(u string, r []string, c *okta.Client) error {
	for _, role := range r {
		if contains(sdk.ValidAdminRoles, role) {
			roleStruct := okta.AssignRoleRequest{Type: role}
			_, _, err := c.User.AssignRoleToUser(context.Background(), u, roleStruct, nil)

			if err != nil {
				return fmt.Errorf("[ERROR] Error Assigning Admin Roles to User: %v", err)
			}
		} else {
			return fmt.Errorf("[ERROR] %v is not a valid Okta role", role)
		}
	}

	return nil
}

func assignGroupsToUser(u string, g []string, c *okta.Client) error {
	for _, group := range g {
		_, err := c.Group.AddUserToGroup(context.Background(), group, u)

		if err != nil {
			return fmt.Errorf("[ERROR] Error Assigning Group to User: %v", err)
		}
	}

	return nil
}

func populateUserProfile(d *schema.ResourceData) *okta.UserProfile {
	profile := okta.UserProfile{}

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

	getSetParams := []string{"city", "costCenter", "countryCode", "department", "displayName", "division",
		"employeeNumber", "honorificPrefix", "honorificSuffix", "locale", "manager", "managerId", "middleName",
		"mobilePhone", "nickName", "organization", "preferredLanguage", "primaryPhone", "profileUrl",
		"secondEmail", "state", "streetAddress", "timezone", "title", "userType", "zipCode"}

	for i := range getSetParams {
		if res, ok := d.GetOk(camelCaseToUnderscore(getSetParams[i])); ok {
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

func listUserOnlyRoles(c *okta.Client, userID string) (userOnlyRoles []*okta.Role, resp *okta.Response, err error) {
	roles, resp, err := c.User.ListAssignedRolesForUser(context.Background(), userID, nil)
	if err != nil {
		return
	}

	for _, role := range roles {
		if role.AssignmentType == userScope {
			userOnlyRoles = append(userOnlyRoles, role)
		}
	}

	return
}

func setAdminRoles(d *schema.ResourceData, c *okta.Client) error {
	roleTypes := make([]interface{}, 0)

	// set all roles currently attached to user in state
	roles, resp, err := listUserOnlyRoles(c, d.Id())

	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusForbidden {
			log.Printf("[INFO] Insufficient permissions to get Admin Roles, skipping.")
		} else {
			return err
		}
	} else {
		for _, role := range roles {
			roleTypes = append(roleTypes, role.Type)
		}
	}

	// set the custom_profile_attributes values
	return setNonPrimitives(d, map[string]interface{}{
		"admin_roles": schema.NewSet(schema.HashString, roleTypes),
	})
}

func setGroups(d *schema.ResourceData, c *okta.Client) error {
	// set all groups currently attached to user in state
	groups, _, err := c.User.ListUserGroups(context.Background(), d.Id())
	if err != nil {
		return err
	}

	groupIDs := make([]interface{}, 0)

	// ignore saving the Everyone group into state so we don't end up with perpetual diffs
	for _, group := range groups {
		if group.Profile.Name != groupProfileEveryone {
			groupIDs = append(groupIDs, group.Id)
		}
	}

	// set the custom_profile_attributes values
	return setNonPrimitives(d, map[string]interface{}{
		"group_memberships": schema.NewSet(schema.HashString, groupIDs),
	})
}

func isCustomUserAttr(key string) bool {
	return !contains(profileKeys, key)
}

func flattenUser(u *okta.User) (map[string]interface{}, error) {
	customAttributes := make(map[string]interface{})
	attrs := map[string]interface{}{}

	for k, v := range *u.Profile {
		if v != nil {
			attrKey := camelCaseToUnderscore(k)

			if isCustomUserAttr(attrKey) {
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

	data, err := json.Marshal(customAttributes)
	if err != nil {
		return nil, fmt.Errorf("failed to load custom_attributes to JSON")
	}
	attrs["custom_profile_attributes"] = string(data)

	return attrs, nil
}

// need to remove from all current admin roles and reassign based on terraform configs when a change is detected
func updateAdminRolesOnUser(u string, r []string, c *okta.Client) error {
	roles, _, err := listUserOnlyRoles(c, u)

	if err != nil {
		return fmt.Errorf("[ERROR] Error Updating Admin Roles On User: %v", err)
	}

	for _, role := range roles {
		_, err = c.User.RemoveRoleFromUser(context.Background(), u, role.Id)

		if err != nil {
			return fmt.Errorf("[ERROR] Error Updating Admin Roles On User: %v", err)
		}
	}

	return assignAdminRolesToUser(u, r, c)
}

// need to remove from all current groups and reassign based on terraform configs when a change is detected
func updateGroupsOnUser(u string, g []string, c *okta.Client) error {
	groups, _, err := c.User.ListUserGroups(context.Background(), u)

	if err != nil {
		return fmt.Errorf("[ERROR] Error Updating Groups On User: %v", err)
	}

	for _, group := range groups {
		if group.Profile.Name != groupProfileEveryone {
			_, err = c.Group.RemoveUserFromGroup(context.Background(), group.Id, u)
			if err != nil {
				return fmt.Errorf("[ERROR] Error Updating Groups On User: %v", err)
			}
		}
	}

	return assignGroupsToUser(u, g, c)
}

// handle setting of user status based on what the current status is because okta
// only allows transitions to certain statuses from other statuses - consult okta User API docs for more info
// https://developer.okta.com/docs/api/resources/users#lifecycle-operations
func updateUserStatus(uid, desiredStatus string, c *okta.Client) error {
	user, _, err := c.User.GetUser(context.Background(), uid)

	if err != nil {
		return err
	}

	var statusErr error
	switch desiredStatus {
	case userStatusSuspended:
		_, statusErr = c.User.SuspendUser(context.Background(), uid)
	case userStatusDeprovisioned:
		_, statusErr = c.User.DeactivateUser(context.Background(), uid, nil)
	case statusActive:
		switch user.Status {
		case userStatusSuspended:
			_, statusErr = c.User.UnsuspendUser(context.Background(), uid)
		case userStatusPasswordExpired:
			// Ignore password expired status. This status is already activated.
			return nil
		case userStatusLockedOut:
			_, statusErr = c.User.UnlockUser(context.Background(), uid)
		default:
			_, _, statusErr = c.User.ActivateUser(context.Background(), uid, nil)
		}
	}

	if statusErr != nil {
		return statusErr
	}

	return waitForStatusTransition(uid, c)
}

// need to wait for user.TransitioningToStatus field to be empty before allowing Terraform to continue
// so the proper current status gets set in the state during the Read operation after a Status update
func waitForStatusTransition(u string, c *okta.Client) error {
	user, _, err := c.User.GetUser(context.Background(), u)

	if err != nil {
		return err
	}

	for {
		if user.TransitioningToStatus == "" {
			return nil
		}
		log.Printf("[INFO] Transitioning to status = %v; waiting for 5 more seconds...", user.TransitioningToStatus)
		time.Sleep(5 * time.Second)
		user, _, err = c.User.GetUser(context.Background(), u)
		if err != nil {
			return err
		}
	}
}
