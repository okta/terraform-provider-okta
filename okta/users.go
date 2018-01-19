package okta

import (
	"encoding/json"
	"log"

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
	// okta user status will be "Password reset" until they complete
	// the okta signup process
	createNewUserAsActive := true

	newUser, _, err := client.Users.Create(newUserTemplate, createNewUserAsActive)
	if err != nil {
		log.Println("[ERROR] Error Creating User: %v", err)
		return err
	}

	// set the resource ID in terraform
	d.SetId(d.Get("email").(string))

	log.Println("[INFO] User Created: %v", *newUser)
	return nil
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[INFO] List User " + d.Get("email").(string))
	client := m.(*Config).oktaClient

	listFilter := client.Users.UserListFilterOptions()
	listFilter.EmailEqualTo = d.Get("email").(string)

	_, err := json.Marshal(listFilter)
	if err != nil {
		log.Println("[ERROR] Error json formatting user filter template: %v", err)
		return err
	}

	userList, _, err := client.Users.ListWithFilter(&listFilter)
	if err != nil {
		log.Println("[ERROR] Error listing user: %v", err)
		return err
	}
	log.Println("[INFO] User List: %v", userList)

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
