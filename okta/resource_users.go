package okta

import (
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceUsers() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,

		CustomizeDiff: func(d *schema.ResourceDiff, v interface{}) error {
			//for an existing user, the login field cannot change
			prev, _ := d.GetChange("login")
			if prev.(string) != "" && d.HasChange("login") {
				return fmt.Errorf("You cannot change the login field for an existing User")
			}
			//regex lovingly lifted from: http://www.golangprograms.com/regular-expression-to-validate-email-address.html
			re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
			if re.MatchString(d.Get("login").(string)) == false {
				return fmt.Errorf("Login field not a valid email address")
			}
			if _, ok := d.GetOk("email"); ok {
				if re.MatchString(d.Get("email").(string)) == false {
					return fmt.Errorf("Email field not a valid email address")
				}
			}
			return nil
		},

		Schema: map[string]*schema.Schema{
			"login": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "User Okta login (must be an email address)",
			},
			"email": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User primary email address. Default = user login",
			},
			"firstname": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "User first name",
			},
			"lastname": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "User last name",
			},
			"middlename": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User middle name",
			},
			"role": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User Okta role",
			},
			"secondemail": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User secondary email address, used for account recovery",
			},
			"honprefix": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User honorific prefix",
			},
			"honsuffix": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User honorific suffix",
			},
			"title": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User title",
			},
			"displayname": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User display name, suitable to show end users",
			},
			"nickname": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User nickname",
			},
			"profileurl": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User online profile (web page)",
			},
			"primaryphone": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User primary phone number",
			},
			"mobilephone": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User mobile phone number",
			},
			"streetaddress": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User street address",
			},
			"city": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User city",
			},
			"state": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User state or region",
			},
			"zipcode": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User zipcode or postal code",
			},
			"countrycode": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User country code",
			},
			"postaladdress": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User mailing address",
			},
			"language": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User preferred language",
			},
			"locale": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User default location",
			},
			"timezone": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User default timezone",
			},
			"usertype": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User employee type",
			},
			"empnumber": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User employee number",
			},
			"costcenter": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User cost center",
			},
			"organization": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User organization",
			},
			"division": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User division",
			},
			"department": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User department",
			},
			"managerid": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Manager ID of User",
			},
			"manager": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Manager of User",
			},
		},
	}
}

func resourceUserCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Creating User %v", d.Get("login").(string))
	client := m.(*Config).oktaClient

	// check if our user exists in Okta, search by login
	filter := client.Users.UserListFilterOptions()
	filter.LoginEqualTo = d.Get("login").(string)
	newUser, _, err := client.Users.ListWithFilter(&filter)
	if len(newUser) == 0 {
		userTemplate("create", d, m)
		if err != nil {
			return err
		}
		return nil
	}
	if len(newUser) > 1 {
		return fmt.Errorf("[ERROR] Retrieved more than one Okta user for the login %v", d.Get("login").(string))
	}
	log.Printf("[INFO] User already exists in Okta. Adding to Terraform")
	// add the user resource to terraform
	d.SetId(newUser[0].ID)

	return nil
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] List User %v", d.Get("login").(string))
	client := m.(*Config).oktaClient

	_, _, err := client.Users.GetByID(d.Id())
	if err != nil {
		// if the user does not exist in okta, delete from terraform state
		if client.OktaErrorCode == "E0000007" {
			d.SetId("")
			return nil
		} else {
			return fmt.Errorf("[ERROR] Error GetByID: %v", err)
		}
	} else {
		userRoles, _, err := client.Users.ListRoles(d.Id())
		if err != nil {
			return fmt.Errorf("[ERROR] Error listing user role: %v", err)
		}
		if userRoles != nil {
			if len(userRoles.Role) > 1 {
				return fmt.Errorf("[ERROR] User has more than one role. This terraform provider presently only supports a single role per user. Please review the user's role assignments in Okta.")
			}
		}
	}

	return nil
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Update User %v", d.Get("login").(string))
	client := m.(*Config).oktaClient
	d.Partial(true)

	_, _, err := client.Users.GetByID(d.Id())
	if err != nil {
		return fmt.Errorf("[ERROR] Error GetByID: %v", err)
	}
	userTemplate("update", d, m)
	if err != nil {
		return err
	}
	d.Partial(false)

	return nil
}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Delete User %v", d.Get("login").(string))
	client := m.(*Config).oktaClient

	userList, _, err := client.Users.GetByID(d.Id())
	if err != nil {
		return fmt.Errorf("[ERROR] Error GetByID: %v", err)
	}
	// must deactivate the user before deletion
	if userList.Status != "DEPROVISIONED" {
		_, err := client.Users.Deactivate(d.Id())
		if err != nil {
			return fmt.Errorf("[ERROR] Error Deactivating user: %v", err)
		}
	}
	// delete the user
	_, err = client.Users.Delete(d.Id())
	if err != nil {
		return fmt.Errorf("[ERROR] Error Deleting user: %v", err)
	}
	// delete the user resource from terraform
	d.SetId("")

	return nil
}

func userTemplate(action string, d *schema.ResourceData, m interface{}) error {
	client := m.(*Config).oktaClient

	template := client.Users.NewUser()
	template.Profile.Login = d.Get("login").(string)
	template.Profile.Email = d.Get("login").(string)
	template.Profile.FirstName = d.Get("firstname").(string)
	template.Profile.LastName = d.Get("lastname").(string)
	if _, ok := d.GetOk("email"); ok {
		template.Profile.Email = d.Get("email").(string)
	}
	if _, ok := d.GetOk("middlename"); ok {
		template.Profile.MiddleName = d.Get("middlename").(string)
	}
	if _, ok := d.GetOk("secondemail"); ok {
		template.Profile.SecondEmail = d.Get("secondemail").(string)
	}
	if _, ok := d.GetOk("honprefix"); ok {
		template.Profile.HonPrefix = d.Get("honprefix").(string)
	}
	if _, ok := d.GetOk("honsuffix"); ok {
		template.Profile.HonSuffix = d.Get("honsuffix").(string)
	}
	if _, ok := d.GetOk("title"); ok {
		template.Profile.Title = d.Get("title").(string)
	}
	if _, ok := d.GetOk("displayname"); ok {
		template.Profile.DisplayName = d.Get("displayname").(string)
	}
	if _, ok := d.GetOk("nickname"); ok {
		template.Profile.NickName = d.Get("nickname").(string)
	}
	if _, ok := d.GetOk("profileurl"); ok {
		template.Profile.ProfileURL = d.Get("profileurl").(string)
	}
	if _, ok := d.GetOk("primaryphone"); ok {
		template.Profile.PrimaryPhone = d.Get("primaryphone").(string)
	}
	if _, ok := d.GetOk("mobilephone"); ok {
		template.Profile.MobilePhone = d.Get("mobilephone").(string)
	}
	if _, ok := d.GetOk("streetaddress"); ok {
		template.Profile.StreetAddress = d.Get("streetaddress").(string)
	}
	if _, ok := d.GetOk("city"); ok {
		template.Profile.City = d.Get("city").(string)
	}
	if _, ok := d.GetOk("state"); ok {
		template.Profile.State = d.Get("state").(string)
	}
	if _, ok := d.GetOk("zipcode"); ok {
		template.Profile.ZipCode = d.Get("zipcode").(string)
	}
	if _, ok := d.GetOk("countrycode"); ok {
		template.Profile.CountryCode = d.Get("countrycode").(string)
	}
	if _, ok := d.GetOk("postaladdress"); ok {
		template.Profile.PostalAddress = d.Get("postaladdress").(string)
	}
	if _, ok := d.GetOk("language"); ok {
		template.Profile.PreferredLanguage = d.Get("language").(string)
	}
	if _, ok := d.GetOk("locale"); ok {
		template.Profile.Locale = d.Get("locale").(string)
	}
	if _, ok := d.GetOk("timezone"); ok {
		template.Profile.Timezone = d.Get("timezone").(string)
	}
	if _, ok := d.GetOk("usertype"); ok {
		template.Profile.UserType = d.Get("usertype").(string)
	}
	if _, ok := d.GetOk("empnumber"); ok {
		template.Profile.EmployeeNumber = d.Get("empnumber").(string)
	}
	if _, ok := d.GetOk("costcenter"); ok {
		template.Profile.CostCenter = d.Get("costcenter").(string)
	}
	if _, ok := d.GetOk("organization"); ok {
		template.Profile.Organization = d.Get("organization").(string)
	}
	if _, ok := d.GetOk("division"); ok {
		template.Profile.Division = d.Get("division").(string)
	}
	if _, ok := d.GetOk("department"); ok {
		template.Profile.Department = d.Get("department").(string)
	}
	if _, ok := d.GetOk("managerid"); ok {
		template.Profile.ManagerID = d.Get("managerid").(string)
	}
	if _, ok := d.GetOk("manager"); ok {
		template.Profile.Manager = d.Get("manager").(string)
	}

	switch action {
	case "create":
		// activate user but send an email to set their password
		// okta user status will be "Password Reset" until they complete
		// the okta signup process
		createNewUserAsActive := true

		newUser, _, err := client.Users.Create(template, createNewUserAsActive)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Creating User: %v", err)
		}
		log.Printf("[INFO] Okta User Created: %+v", newUser)

		// assign the user a role, if specified
		if _, ok := d.GetOk("role"); ok {
			log.Printf("[INFO] Assigning role: " + d.Get("role").(string))
			_, err := client.Users.AssignRole(newUser.ID, d.Get("role").(string))
			if err != nil {
				return fmt.Errorf("[ERROR] Error assigning role to user: %v", err)
			}
		}
		// add the user resource to terraform
		d.SetId(newUser.ID)

	case "update":
		updateUser, _, err := client.Users.Update(template, d.Id())
		if err != nil {
			return fmt.Errorf("[ERROR] Error Updating User: %v", err)
		}
		log.Printf("[INFO] Okta User Updated: %+v", updateUser)

		userRoles, _, err := client.Users.ListRoles(d.Id())
		if err != nil {
			return fmt.Errorf("[ERROR] Error listing user role: %v", err)
		}

		if d.HasChange("role") {
			if userRoles != nil {
				log.Printf("[INFO] Removing role: " + userRoles.Role[0].Type)
				_, err = client.Users.UnAssignRole(d.Id(), userRoles.Role[0].ID)
				if err != nil {
					return fmt.Errorf("[ERROR] Error removing role from user: %v", err)
				}
			}
			if _, ok := d.GetOk("role"); ok {
				log.Printf("[INFO] Assigning role: " + d.Get("role").(string))
				_, err = client.Users.AssignRole(d.Id(), d.Get("role").(string))
				if err != nil {
					return fmt.Errorf("[ERROR] Error assigning role to user: %v", err)
				}
			}
		}

	default:
		return fmt.Errorf("[ERROR] userTemplate action only supports \"create\" and \"update\"")
	}

	return nil
}
