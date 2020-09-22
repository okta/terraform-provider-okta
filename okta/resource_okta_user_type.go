package okta

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-okta/sdk"
)

func resourceUserType() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserTypeCreate,
		Read:   resourceUserTypeRead,
		Update: resourceUserTypeUpdate,
		Delete: resourceUserTypeDelete,
		Exists: resourceUserTypeExists,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The display name for the type",
			},
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: "The display name for the type	",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A human-readable description of the type",
			},
		},
	}
}

func buildUserType(d *schema.ResourceData) *sdk.UserType {
	return &sdk.UserType{
		Name:        d.Get("name").(string),
		DisplayName: d.Get("display_name").(string),
		Description: d.Get("description").(string),
	}
}

func resourceUserTypeCreate(d *schema.ResourceData, m interface{}) error {
	client := getSupplementFromMetadata(m)
	userType := buildUserType(d)
	newUserType, _, err := client.CreateUserType(*userType, nil)
	if err != nil {
		return err
	}

	d.SetId(newUserType.Id)
	if err != nil {
		return err
	}

	return resourceUserTypeRead(d, m)
}

func resourceUserTypeUpdate(d *schema.ResourceData, m interface{}) error {
	client := getSupplementFromMetadata(m)
	userType := buildUserType(d)
	_, _, err := client.UpdateUserType(d.Id(), *userType, nil)

	if err != nil {
		return err
	}

	return resourceUserTypeRead(d, m)
}

func resourceUserTypeRead(d *schema.ResourceData, m interface{}) error {
	userType, resp, err := getSupplementFromMetadata(m).GetUserType(d.Id())

	if is404(resp.StatusCode) {
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	d.Set("name", userType.Name)
	d.Set("display_name", userType.DisplayName)
	d.Set("description", userType.Description)

	return nil
}

func resourceUserTypeDelete(d *schema.ResourceData, m interface{}) error {
	client := getSupplementFromMetadata(m)
	_, err := client.DeleteUserType(d.Id())

	return err
}

func resourceUserTypeExists(d *schema.ResourceData, m interface{}) (bool, error) {
	_, _, err := getSupplementFromMetadata(m).GetUserType(d.Id())
	if err != nil {
		return false, fmt.Errorf("[ERROR] Error Getting User Type: %v", err)
	}
	return true, nil
}
