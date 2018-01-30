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

	_, _, err := client.Users.GetByID(d.Get("email").(string))

	switch {

	// only create the user in okta if they do not already exist
	case client.OktaErrorCode == "E0000007":

		newUserTemplate := client.Users.NewUser()
		newUserTemplate.Profile.FirstName = d.Get("firstname").(string)
		newUserTemplate.Profile.LastName = d.Get("lastname").(string)
		newUserTemplate.Profile.Login = d.Get("email").(string)
		newUserTemplate.Profile.Email = newUserTemplate.Profile.Login

		_, err = json.Marshal(newUserTemplate)
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
		log.Println("[INFO] Okta User Created: %v", *newUser)

	case err != nil:
		return err

	default:
		log.Println("[INFO] User %v already exists in Okta. Adding to Terraform.",
			d.Get("email").(string))
	}

	// add the user resource to terraform
	d.SetId(d.Get("email").(string))

	return nil
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[INFO] List User " + d.Get("email").(string))
	client := m.(*Config).oktaClient

	userList, _, err := client.Users.GetByID(d.Get("email").(string))
	if err != nil {
		// if the user does not exist in okta, delete from terraform
		if client.OktaErrorCode == "E0000007" {
			d.SetId("")
		} else {
			log.Println("[ERROR] Error listing user: %v", err)
			return err
		}
	}
	log.Println("[INFO] User List: %v", userList)

	return nil
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("[INFO] Update User " + d.Get("email").(string))

	if d.HasChange("firstname") || d.HasChange("lastname") {
		client := m.(*Config).oktaClient

		userList, _, err := client.Users.GetByID(d.Get("email").(string))
		if err != nil {
			return err
		}
		log.Println("[INFO] User List: %v", userList)

		updateUserTemplate := client.Users.NewUser()
		updateUserTemplate.Profile.FirstName = d.Get("firstname").(string)
		updateUserTemplate.Profile.LastName = d.Get("lastname").(string)
		updateUserTemplate.Profile.Login = d.Get("email").(string)
		updateUserTemplate.Profile.Email = updateUserTemplate.Profile.Login

		_, err = json.Marshal(updateUserTemplate)
		if err != nil {
			log.Println("[ERROR] Error json formatting update user template: %v", err)
			return err
		}

		// update the user in okta
		updateUser, _, err := client.Users.Update(updateUserTemplate, d.Get("email").(string))
		if err != nil {
			log.Println("[ERROR] Error Updating User: %v", err)
			return err
		}
		log.Println("[INFO] Okta User Updated: %v", *updateUser)

		// update the user resource in terraform with the new value(s)
		d.Set("firstname", d.Get("firstname").(string))
		d.Set("lastname", d.Get("lastname").(string))
	}

	return nil
}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("[INFO] Delete User Terraform Resource" + d.Get("email").(string))

	// how do we want to handle user deletion?

	return nil
}
