package okta

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

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
				Optional:     true,
				Description:  "User primary email address. Default = user login",
				ValidateFunc: matchEmailRegexp,
				// supress diff if no email value is given since it defauls to login
				// and it's technically required by Okta
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool { return new == "" },
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
	if len(d.Get("admin_roles").([]interface{})) > 0 {
		err = assignAdminRolesToUser(user.Id, d.Get("admin_roles").([]interface{}), client)

		if err != nil {
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

	exists, err := resourceUserExists(d, m)

	if err != nil {
		return err
	}

	// set Id to "" if the user can't be found to remove from state
	if exists != true {
		d.SetId("")
		return nil
	}

	user, _, err := client.User.GetUser(d.Id())

	if err != nil {
		return fmt.Errorf("[ERROR] Error Getting User from Okta: %v", err)
	}

	d.Set("status", user.Status)

	// set all the attributes in state that were returned from user.Profile
	for k, v := range *user.Profile {
		if v != nil {
			attribute := camelCaseToUnderscore(k)

			log.Printf("[INFO] Setting %v to %v", attribute, v)
			if err := d.Set(attribute, v); err != nil {
				return fmt.Errorf("error setting %s for resource %s: %s", attribute, d.Id(), err)
			}
		}
	}

	// set all roles currently attached to user in state
	roles, _, err := client.User.ListAssignedRoles(d.Id(), nil)

	if err != nil {
		return err
	}

	r := make([]string, 0)
	for _, role := range roles {
		r = append(r, role.Type)
	}

	d.Set("admin_roles", r)

	return nil
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

		if d.HasChange("admin_roles") {
			err := updateAdminRolesOnUser(d.Id(), d.Get("admin_roles").([]interface{}), client)

			if err != nil {
				return err
			}
		}

		if err != nil {
			return fmt.Errorf("[ERROR] Error Updating User in Okta: %v", err)
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

func populateUserProfile(d *schema.ResourceData) *okta.UserProfile {
	profile := okta.UserProfile{}

	profile["firstName"] = d.Get("first_name").(string)
	profile["lastName"] = d.Get("last_name").(string)
	profile["login"] = d.Get("login").(string)

	if _, ok := d.GetOk("email"); ok {
		profile["email"] = d.Get("email").(string)
	} else {
		profile["email"] = d.Get("login").(string)
	}

	if _, ok := d.GetOk("city"); ok {
		profile["city"] = d.Get("city").(string)
	}

	if _, ok := d.GetOk("cost_center"); ok {
		profile["costCenter"] = d.Get("cost_center").(string)
	}

	if _, ok := d.GetOk("country_code"); ok {
		profile["countryCode"] = d.Get("country_code").(string)
	}

	if _, ok := d.GetOk("department"); ok {
		profile["department"] = d.Get("department").(string)
	}

	if _, ok := d.GetOk("display_name"); ok {
		profile["displayName"] = d.Get("display_name").(string)
	}

	if _, ok := d.GetOk("division"); ok {
		profile["division"] = d.Get("division").(string)
	}

	if _, ok := d.GetOk("employee_number"); ok {
		profile["employeeNumber"] = d.Get("employee_number").(string)
	}

	if _, ok := d.GetOk("honorific_prefix"); ok {
		profile["honorificPrefix"] = d.Get("honorific_prefix").(string)
	}

	if _, ok := d.GetOk("honorific_suffix"); ok {
		profile["honorificSuffix"] = d.Get("honorific_suffix").(string)
	}

	if _, ok := d.GetOk("locale"); ok {
		profile["locale"] = d.Get("locale").(string)
	}

	if _, ok := d.GetOk("manager"); ok {
		profile["manager"] = d.Get("manager").(string)
	}

	if _, ok := d.GetOk("manager_id"); ok {
		profile["managerId"] = d.Get("manager_id").(string)
	}

	if _, ok := d.GetOk("middle_name"); ok {
		profile["middleName"] = d.Get("middle_name").(string)
	}

	if _, ok := d.GetOk("mobile_phone"); ok {
		profile["mobilePhone"] = d.Get("mobile_phone").(string)
	}

	if _, ok := d.GetOk("nick_name"); ok {
		profile["nickName"] = d.Get("nick_name").(string)
	}

	if _, ok := d.GetOk("organization"); ok {
		profile["organization"] = d.Get("organization").(string)
	}

	if _, ok := d.GetOk("postal_address"); ok {
		profile["postalAddress"] = d.Get("postal_address").(string)
	}

	if _, ok := d.GetOk("preferred_language"); ok {
		profile["preferredLanguage"] = d.Get("preferred_language").(string)
	}

	if _, ok := d.GetOk("primary_phone"); ok {
		profile["primaryPhone"] = d.Get("primary_phone").(string)
	}

	if _, ok := d.GetOk("profile_url"); ok {
		profile["profileUrl"] = d.Get("profile_url").(string)
	}

	if _, ok := d.GetOk("second_email"); ok {
		profile["secondEmail"] = d.Get("second_email").(string)
	}

	if _, ok := d.GetOk("state"); ok {
		profile["state"] = d.Get("state").(string)
	}

	if _, ok := d.GetOk("street_address"); ok {
		profile["streetAddress"] = d.Get("street_address").(string)
	}

	if _, ok := d.GetOk("timezone"); ok {
		profile["timezone"] = d.Get("timezone").(string)
	}

	if _, ok := d.GetOk("title"); ok {
		profile["title"] = d.Get("title").(string)
	}

	if _, ok := d.GetOk("user_type"); ok {
		profile["userType"] = d.Get("user_type").(string)
	}

	if _, ok := d.GetOk("zip_code"); ok {
		profile["zipCode"] = d.Get("zip_code").(string)
	}

	return &profile
}

func assignAdminRolesToUser(u string, r []interface{}, c *okta.Client) error {
	validRoles := []string{"SUPER_ADMIN", "ORG_ADMIN", "API_ACCESS_MANAGEMENT_ADMIN", "APP_ADMIN", "USER_ADMIN", "MOBILE_ADMIN", "READ_ONLY_ADMIN", "HELP_DESK_ADMIN"}

	for _, role := range r {
		if contains(validRoles, role.(string)) {
			roleStruct := okta.Role{Type: role.(string)}
			_, _, err := c.User.AddRoleToUser(u, roleStruct)

			if err != nil {
				return fmt.Errorf("[ERROR] Error Assigning Admin Roles to User: %v", err)
			}
		} else {
			return fmt.Errorf("[ERROR] %v is not a valid Okta role", role)
		}
	}

	return nil
}

// need to remove from all current admin roles and reassign based on terraform configs when a change is detected
func updateAdminRolesOnUser(u string, r []interface{}, c *okta.Client) error {
	roles, _, err := c.User.ListAssignedRoles(u, nil)

	if err != nil {
		return fmt.Errorf("[ERROR] Error Updating Admin Roles On User: %v", err)
	}

	for _, role := range roles {
		_, err := c.User.RemoveRoleFromUser(u, role.Id)

		if err != nil {
			return fmt.Errorf("[ERROR] Error Updating Admin Roles On User: %v", err)
		}
	}

	err = assignAdminRolesToUser(u, r, c)

	if err != nil {
		return err
	}

	return nil
}

// regex lovingly lifted from: http://www.golangprograms.com/regular-expression-to-validate-email-address.html
func matchEmailRegexp(val interface{}, key string) (warnings []string, errors []error) {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if re.MatchString(val.(string)) == false {
		errors = append(errors, fmt.Errorf("%s field not a valid email address", key))
	}
	return warnings, errors
}

// handle setting of user status based on what the current status is because okta
// only allows transitions to certain statuses from other statuses - consult okta User API docs for more info
// https://developer.okta.com/docs/api/resources/users#lifecycle-operations
func updateUserStatus(u string, d string, c *okta.Client) error {
	user, _, err := c.User.GetUser(u)

	if err != nil {
		return err
	}

	var statusErr error
	switch d {
	case "SUSPENDED":
		_, statusErr = c.User.SuspendUser(u)
	case "DEPROVISIONED":
		_, statusErr = c.User.DeactivateUser(u)
	case "ACTIVE":
		if user.Status == "SUSPENDED" {
			_, statusErr = c.User.UnsuspendUser(u)
		} else {
			_, _, statusErr = c.User.ActivateUser(u, nil)
		}
	}

	if statusErr != nil {
		return statusErr
	}

	err = waitForStatusTransition(u, c)

	if err != nil {
		return err
	}

	return nil
}

// need to wait for user.TransitioningToStatus field to be empty before allowing Terraform to continue
// so the proper current status gets set in the state during the Read operation after a Status update
func waitForStatusTransition(u string, c *okta.Client) error {
	user, _, err := c.User.GetUser(u)

	if err != nil {
		return err
	}

	for {
		if user.TransitioningToStatus == "" {
			return nil
		} else {
			log.Printf("[INFO] Transitioning to status = %v; waiting for 5 more seconds...", user.TransitioningToStatus)
			time.Sleep(5 * time.Second)

			user, _, err = c.User.GetUser(u)

			if err != nil {
				return err
			}
		}
	}
}
