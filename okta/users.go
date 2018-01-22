package okta

import (
	"encoding/json"
	"log"
	"strings"
	"unsafe"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceUsers() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,

		Schema: map[string]*schema.Schema{
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

			"email": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "User email address",
			},
		},
	}
}

func resourceUserCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("[INFO] Creating User" + d.Get("email").(string))
	client := m.(*Config).oktaClient

	newUserTemplate := client.Users.NewUser()
	newUserTemplate.Profile.FirstName = d.Get("firstname").(string)
	newUserTemplate.Profile.LastName = d.Get("lastname").(string)
	newUserTemplate.Profile.Login = d.Get("email").(string)
	newUserTemplate.Profile.Email = newUserTemplate.Profile.Login

	_, err := json.Marshal(newUserTemplate)
	if err != nil {
		log.Println("[ERROR] Error json formatting new user template: %v", err)
		return err
	}

	// activate user but send an email to set their password
	// okta user status will be "Password Reset" until they complete
	// the okta signup process
	createNewUserAsActive := true

	newUser, _, err := client.Users.Create(newUserTemplate, createNewUserAsActive)
	if err != nil {
		log.Println("[ERROR] Error Creating User: %v", err)
		return err
	}

	// add the user resource to terraform
	d.SetId(d.Get("email").(string))

	log.Println("[INFO] User Created: %v", *newUser)
	return nil
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[INFO] List User " + d.Get("email").(string))
	client := m.(*Config).oktaClient

	userList, _, err := client.Users.GetByID(d.Get("email").(string))
	if err != nil {
		// this is a placeholder
		// need to add to the sdk spitting out okta error codes
		type error interface {
			Error() string
		}
		arr := strings.Split(err.Error(), ",")
		brr := strings.Split(arr[1], ":")
		log.Println("[ERROR] BRR: %v", brr[1])
		if strings.TrimSpace(brr[1]) == "E0000007" {
			d.SetId("")
			return nil
		}
		log.Println("[ERROR] Error listing user: %v", err)
		return err
	}
	log.Println("[INFO] User List: %v", userList)

	// remove our user resource from terraform if our okta query
	// results in an empty array
	if unsafe.Sizeof(userList) == 0 {
		d.SetId("")
		return nil
	}
	return nil
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
//	log.Println("[INFO] Destroy Managed Users")
//	client := m.(*Config).oktaClient

//	destroy user
//	d.SetId(")
	return nil
}
