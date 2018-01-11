package okta

import (
	"encoding/json"
	"log"
	"time"

	"github.com/articulate/oktasdk-go/okta"
	"github.com/hashicorp/terraform/helper/schema"
)

func Users() *schema.Resource {
	return &schema.Resource{
		Create: UserCreate,
		Read:   UserRead,
		Update: UserUpdate,
		Delete: UserDelete,

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

func UserCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("[INFO] Creating User")
	client := okta.NewClient(nil, orgName, apiToken, isProductionOKTAORG)
	log.Println("[INFO] Client Base URL: %v\n\n", client.BaseURL)

	newUserTemplate := client.Users.NewUser()
	newUserTemplate.Profile.FirstName = d.Get("firstname").(string)
	newUserTemplate.Profile.LastName = d.Get("lastname").(string) + time.Now().Format("2006-01-02")
	newUserTemplate.Profile.Login = d.Get("email").(string)
	newUserTemplate.Profile.Email = newUserTemplate.Profile.Login

	jsonTest, _ := json.Marshal(newUserTemplate)

	log.Println("[INFO] User Json\n\t%v\n\n", string(jsonTest))
	createNewUserAsActive := false

	newUser, _, err := client.Users.Create(newUserTemplate, createNewUserAsActive)

	if err != nil {

		log.Println("[ERROR] Error Creating User:\n \t%v", err)
		return err
	}
	return nil
}

func UserRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func UserUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func UserDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
