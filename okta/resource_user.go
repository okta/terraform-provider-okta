package okta

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/query"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,
		Exists: resourceUserExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"admin_roles": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "User Okta admin roles - ie. ['APP_ADMIN', 'USER_ADMIN']",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"city": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User city",
			},
			"cost_center": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User cost center",
			},
			"country_code": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User country code",
			},
			"custom_profile_attributes": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			"department": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User department",
			},
			"display_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User display name, suitable to show end users",
			},
			"division": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User division",
			},
			"email": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "User primary email address",
				ValidateFunc: matchEmailRegexp,
			},
			"employee_number": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User employee number",
			},
			"first_name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "User first name",
			},
			"group_memberships": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The groups that you want this user to be a part of",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"honorific_prefix": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User honorific prefix",
			},
			"honorific_suffix": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User honorific suffix",
			},
			"last_name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "User last name",
			},
			"locale": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User default location",
			},
			"login": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "User Okta login (must be an email address)",
				ForceNew:     true,
				ValidateFunc: matchEmailRegexp,
			},
			"manager": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Manager of User",
			},
			"manager_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Manager ID of User",
			},
			"middle_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User middle name",
			},
			"mobile_phone": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User mobile phone number",
			},
			"nick_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User nickname",
			},
			"organization": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User organization",
			},
			"postal_address": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User mailing address",
			},
			"preferred_language": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User preferred language",
			},
			"primary_phone": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User primary phone number",
			},
			"profile_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User online profile (web page)",
			},
			"second_email": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User secondary email address, used for account recovery",
			},
			"state": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User state or region",
			},
			"status": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The status of the User in Okta - remove to set user back to active/provisioned",
				Default:      "ACTIVE",
				ValidateFunc: validation.StringInSlice([]string{"ACTIVE", "STAGED", "DEPROVISIONED", "SUSPENDED"}, false),
				// ignore diff changing to ACTIVE if state is set to PROVISIONED
				// since this is a similar status in Okta terms
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return old == "PROVISIONED" && new == "ACTIVE"
				},
			},
			"street_address": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User street address",
			},
			"timezone": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User default timezone",
			},
			"title": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User title",
			},
			"user_type": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User employee type",
			},
			"zip_code": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User zipcode or postal code",
			},
		},
	}
}

func resourceUserCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Create User for %v", d.Get("login").(string))

	client := m.(*Config).oktaClient
	profile := populateUserProfile(d)

	qp := query.NewQueryParams()

	// setting activate to false on user creation will leave the user with a status of STAGED
	if d.Get("status").(string) == "STAGED" {
		qp = query.NewQueryParams(query.WithActivate(false))
	}

	userBody := okta.User{Profile: profile}
	user, _, err := client.User.CreateUser(userBody, qp)

	if err != nil {
		return fmt.Errorf("[ERROR] Error Creating User from Okta: %v", err)
	}

	// set the user id into state before setting roles and status in case they fail
	d.SetId(user.Id)

	// role assigning can only happen after the user is created so order matters here
	roles := convertInterfaceToStringArrNullable(d.Get("admin_roles"))
	if roles != nil {
		if err = assignAdminRolesToUser(user.Id, roles, client); err != nil {
			return err
		}
	}

	// group assigning can only happen after the user is created as well
	groups := convertInterfaceToStringArrNullable(d.Get("group_memberships"))
	if groups != nil {
		if err = assignGroupsToUser(user.Id, groups, client); err != nil {
			return err
		}
	}

	// status changing can only happen after user is created as well
	if d.Get("status").(string) == "SUSPENDED" || d.Get("status").(string) == "DEPROVISIONED" {
		err := updateUserStatus(user.Id, d.Get("status").(string), client)

		if err != nil {
			return fmt.Errorf("[ERROR] Error Updating Status for User: %v", err)
		}
	}

	return resourceUserRead(d, m)
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] List User %v", d.Get("login").(string))
	client := m.(*Config).oktaClient

	user, _, err := client.User.GetUser(d.Id())

	if err != nil {
		return fmt.Errorf("[ERROR] Error Getting User from Okta: %v", err)
	}

	d.Set("status", user.Status)

	if err = setUserProfileAttributes(d, user); err != nil {
		return err
	}

	if err = setAdminRoles(d, client); err != nil {
		return err
	}

	return setGroups(d, client)
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Update User %v", d.Get("login").(string))

	if d.Get("status").(string) == "STAGED" {
		return fmt.Errorf("[ERROR] Okta will not allow a user to be updated to STAGED. Can set to STAGED on user creation only.")
	}

	client := m.(*Config).oktaClient

	// run the update status func first so a user that was previously deprovisioned
	// can be updated further if it's status changed in it's terraform configs
	if d.HasChange("status") {
		err := updateUserStatus(d.Id(), d.Get("status").(string), client)

		if err != nil {
			return fmt.Errorf("[ERROR] Error Updating Status for User: %v", err)
		}
	}

	if d.Get("status") == "DEPROVISIONED" {
		return fmt.Errorf("[ERROR] Cannot update a DEPROVISIONED user")
	} else {
		profile := populateUserProfile(d)
		userBody := okta.User{Profile: profile}

		_, _, err := client.User.UpdateUser(d.Id(), userBody)

		if err != nil {
			return fmt.Errorf("[ERROR] Error Updating User in Okta: %v", err)
		}

		if d.HasChange("admin_roles") {
			roles := convertInterfaceToStringArr(d.Get("admin_roles"))
			if err := updateAdminRolesOnUser(d.Id(), roles, client); err != nil {
				return err
			}
		}

		if d.HasChange("group_memberships") {
			groups := convertInterfaceToStringArr(d.Get("group_memberships"))
			if err := updateGroupsOnUser(d.Id(), groups, client); err != nil {
				return err
			}
		}
	}

	return resourceUserRead(d, m)
}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Delete User %v", d.Get("login").(string))
	client := m.(*Config).oktaClient

	// only deprovisioned users can be deleted fully from okta
	// make two passes on the user if they aren't deprovisioned already to deprovision them first
	passes := 2

	if d.Get("status") == "DEPROVISIONED" {
		passes = 1
	}

	for i := 0; i < passes; i++ {
		_, err := client.User.DeactivateOrDeleteUser(d.Id())

		if err != nil {
			return fmt.Errorf("[ERROR] Error Deleting User in Okta: %v", err)
		}
	}

	return nil
}

func resourceUserExists(d *schema.ResourceData, m interface{}) (bool, error) {
	log.Printf("[INFO] Checking Exists for User %v", d.Get("login").(string))
	client := m.(*Config).oktaClient

	_, resp, err := client.User.GetUser(d.Id())

	if err != nil {
		return false, fmt.Errorf("[ERROR] Error Getting User from Okta: %v", err)
	}

	if strings.Contains(resp.Response.Status, "404") {
		return false, nil
	}

	return true, nil
}
