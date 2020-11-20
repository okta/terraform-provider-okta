package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
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
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The display name for the type",
			},
			"display_name": {
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

func buildUserType(d *schema.ResourceData) okta.UserType {
	return okta.UserType{
		Name:        d.Get("name").(string),
		DisplayName: d.Get("display_name").(string),
		Description: d.Get("description").(string),
	}
}

func resourceUserTypeCreate(d *schema.ResourceData, m interface{}) error {
	newUserType, _, err := getOktaClientFromMetadata(m).UserType.CreateUserType(context.Background(), buildUserType(d))
	if err != nil {
		return err
	}
	d.SetId(newUserType.Id)
	return resourceUserTypeRead(d, m)
}

func resourceUserTypeUpdate(d *schema.ResourceData, m interface{}) error {
	userType := buildUserType(d)
	_, _, err := getOktaClientFromMetadata(m).UserType.UpdateUserType(context.Background(), d.Id(), userType)
	if err != nil {
		return err
	}
	return resourceUserTypeRead(d, m)
}

func resourceUserTypeRead(d *schema.ResourceData, m interface{}) error {
	userType, resp, err := getOktaClientFromMetadata(m).UserType.GetUserType(context.Background(), d.Id())
	if resp != nil && is404(resp.StatusCode) {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}
	_ = d.Set("name", userType.Name)
	_ = d.Set("display_name", userType.DisplayName)
	_ = d.Set("description", userType.Description)
	return nil
}

func resourceUserTypeDelete(d *schema.ResourceData, m interface{}) error {
	_, err := getOktaClientFromMetadata(m).UserType.DeleteUserType(context.Background(), d.Id())
	return err
}

func resourceUserTypeExists(d *schema.ResourceData, m interface{}) (bool, error) {
	_, resp, err := getOktaClientFromMetadata(m).UserType.GetUserType(context.Background(), d.Id())
	if err != nil {
		return false, fmt.Errorf("failed to get user type: %v", err)
	}
	if resp != nil && is404(resp.StatusCode) {
		return false, nil
	}
	return true, nil
}
