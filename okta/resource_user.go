package okta

import (
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/okta/okta-sdk-golang/okta"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,

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

	userBody := okta.User{Profile: profile}

	user, _, err := client.User.CreateUser(userBody, nil)

	if err != nil {
		return fmt.Errorf("[ERROR] Error Creating User from Okta: %v", err)
	}

	d.SetId(user.Id)

	return nil
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
	return nil
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

//regex lovingly lifted from: http://www.golangprograms.com/regular-expression-to-validate-email-address.html
func matchEmailRegexp(val interface{}, key string) (warnings []string, errors []error) {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if re.MatchString(val.(string)) == false {
		errors = append(errors, fmt.Errorf("%s field not a valid email address", key))
	}
	return warnings, errors
}
