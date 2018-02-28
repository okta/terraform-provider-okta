package okta

import (
	"encoding/json"
	"fmt"
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
			"role": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User okta role",
			},
		},
	}
}

func resourceUserCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Creating User %v", d.Get("email").(string))
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
			return fmt.Errorf("[ERROR] Error json formatting new user template: %v", err)
		}

		// activate user but send an email to set their password
		// okta user status will be "Password Reset" until they complete
		// the okta signup process
		createNewUserAsActive := true

		newUser, _, err := client.Users.Create(newUserTemplate, createNewUserAsActive)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Creating User: %v", err)
		}
		log.Printf("[INFO] Okta User Created: %+v", newUser)

		// assign the user a role, if specified
		if d.Get("role").(string) != "" {
			log.Printf("[INFO] Assigning role: " + d.Get("role").(string))
			_, err := client.Users.AssignRole(newUser.ID, d.Get("role").(string))
			if err != nil {
				return fmt.Errorf("[ERROR] Error assigning role to user: %v", err)
			}
		}

	case err != nil:
		return fmt.Errorf("[ERROR] Error GetByID: %v", err)

	default:
		log.Printf("[INFO] User %v already exists in Okta. Adding to Terraform.",
			d.Get("email").(string))
	}

	// add the user resource to terraform
	d.SetId(d.Get("email").(string))

	return nil
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] List User %v", d.Get("email").(string))
	client := m.(*Config).oktaClient

	userList, _, err := client.Users.GetByID(d.Get("email").(string))
	if err != nil {
		// if the user does not exist in okta, delete from terraform state
		if client.OktaErrorCode == "E0000007" {
			d.SetId("")
			return nil
		} else {
			return fmt.Errorf("[ERROR] Error GetByID: %v", err)
		}
	} else {
		userRoles, _, err := client.Users.ListRoles(userList.ID)
		if err != nil {
			return fmt.Errorf("[ERROR] Error listing user role: %v", err)
		}
		if userRoles != nil {
			if len(userRoles.Role) > 1 {
				return fmt.Errorf("[ERROR] User has more than one role. This terraform provider presently only supports a single role per user. Please review the role assignments in Okta for this user.")
			}
			log.Printf("[INFO] List Role: %+v", userRoles.Role[0])
			// if they differ, change terraform user state role to okta user role
			if d.Get("role").(string) != userRoles.Role[0].Type {
				d.Set("role", userRoles.Role[0].Type)
			}
		}
	}
	log.Printf("[INFO] List User: %+v", userList)

	return nil
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Update User %v", d.Get("email").(string))
	client := m.(*Config).oktaClient

	userList, _, err := client.Users.GetByID(d.Get("email").(string))
	if err != nil {
		return fmt.Errorf("[ERROR] Error GetByID: %v", err)
	}
	log.Printf("[INFO] User List: %+v", userList)

	if d.HasChange("firstname") || d.HasChange("lastname") {
		updateUserTemplate := client.Users.NewUser()
		updateUserTemplate.Profile.FirstName = d.Get("firstname").(string)
		updateUserTemplate.Profile.LastName = d.Get("lastname").(string)
		updateUserTemplate.Profile.Login = d.Get("email").(string)
		updateUserTemplate.Profile.Email = updateUserTemplate.Profile.Login

		_, err = json.Marshal(updateUserTemplate)
		if err != nil {
			return fmt.Errorf("[ERROR] Error json formatting update user template: %v", err)
		}

		// update the user in okta
		updateUser, _, err := client.Users.Update(updateUserTemplate, d.Get("email").(string))
		if err != nil {
			return fmt.Errorf("[ERROR] Error Updating User: %v", err)
		}
		log.Printf("[INFO] Okta User Updated: %+v", updateUser)

		// update the user resource in terraform
		d.Set("firstname", d.Get("firstname").(string))
		d.Set("lastname", d.Get("lastname").(string))
	}

	if d.HasChange("role") {
		userRoles, _, err := client.Users.ListRoles(userList.ID)
		if err != nil {
			return fmt.Errorf("[ERROR] Error listing user role: %v", err)
		}
		switch {

		case userRoles == nil && d.Get("role").(string) != "":
			log.Printf("[INFO] Assigning role: " + d.Get("role").(string))
			_, err = client.Users.AssignRole(userList.ID, d.Get("role").(string))
			if err != nil {
				return fmt.Errorf("[ERROR] Error assigning role to user: %v", err)
			}

		case userRoles != nil && d.Get("role").(string) != "":
			log.Printf("[INFO] Changing role: " + d.Get("role").(string) + " to " + userRoles.Role[0].Type)
			_, err = client.Users.UnAssignRole(userList.ID, userRoles.Role[0].ID)
			if err != nil {
				return fmt.Errorf("[ERROR] Error removing role from user: %v", err)
			}
			_, err = client.Users.AssignRole(userList.ID, d.Get("role").(string))
			if err != nil {
				return fmt.Errorf("[ERROR] Error assigning role to user: %v", err)
			}

		case userRoles != nil && d.Get("role").(string) == "":
			log.Printf("[INFO] Removing role: " + d.Get("role").(string))
			_, err = client.Users.UnAssignRole(userList.ID, userRoles.Role[0].ID)
			if err != nil {
				return fmt.Errorf("[ERROR] Error removing role from user: %v", err)
			}

		default:
			return fmt.Errorf("User role changed but Terraform was unable to apply. Please investigate.")
		}
	}

	return nil
}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Delete User %v", d.Get("email").(string))
	client := m.(*Config).oktaClient

	userList, _, err := client.Users.GetByID(d.Get("email").(string))
	if err != nil {
		return fmt.Errorf("[ERROR] Error GetByID: %v", err)
	}

	// must deactivate the user before deletion
	if userList.Status != "DEPROVISIONED" {
		_, err := client.Users.Deactivate(d.Get("email").(string))
		if err != nil {
			return fmt.Errorf("[ERROR] Error Deactivating user: %v", err)
		}
	}

	// now! delete the user
	deleteUser, err := client.Users.Delete(d.Get("email").(string))
	if err != nil {
		return fmt.Errorf("[ERROR] Error Deleting user: %v", err)
	}
	log.Printf("[INFO] Okta User Deleted: %+v", deleteUser)

	// delete the user resource from terraform
	d.SetId("")

	return nil
}
