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
	log.Println("[INFO] Creating User")
	log.Println("[INFO] " + d.Get("email").(string))
	client := m.(*Config).oktaClient
	log.Printf("USERS RESOURCE %+v\n", client)

	newUserTemplate := client.Users.NewUser()
	newUserTemplate.Profile.FirstName = d.Get("firstname").(string)
	newUserTemplate.Profile.LastName = d.Get("lastname").(string)
	newUserTemplate.Profile.Login = d.Get("email").(string)
	newUserTemplate.Profile.Email = newUserTemplate.Profile.Login

	jsonTest, _ := json.Marshal(newUserTemplate)

	log.Println("[INFO] User Json\n\t%v\n\n", string(jsonTest))
	createNewUserAsActive := false

	newUser, _, err := client.Users.Create(newUserTemplate, createNewUserAsActive)
	log.Println("[INFO] newUser \n\t%v\n\n", newUser)

	if err != nil {

		log.Println("[ERROR] Error Creating User:\n \t%v", err)
		return err
	}
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
