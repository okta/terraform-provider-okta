package okta

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func resourceAppUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppUserCreate,
		Exists: resourceAppUserExists,
		Read:   resourceAppUserRead,
		Update: resourceAppUserUpdate,
		Delete: resourceAppUserDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				parts := strings.Split(d.Id(), "/")
				if len(parts) != 2 {
					return nil, errors.New("invalid resource import specifier. Use: terraform import <app_id>/<user_id>")
				}

				_ = d.Set("app_id", parts[0])
				_ = d.Set("user_id", parts[1])

				assignment, _, err := getOktaClientFromMetadata(m).Application.
					GetApplicationUser(context.Background(), parts[0], parts[1], nil)

				if err != nil {
					return nil, err
				}

				d.SetId(assignment.Id)

				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "App to associate user with",
			},
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "User associated with the application",
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
			},
			"profile": {
				Type:      schema.TypeString,
				StateFunc: normalizeDataJSON,
				Optional:  true,
				Default:   "{}",
			},
		},
	}
}

func resourceAppUserExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := getOktaClientFromMetadata(m)
	g, _, err := client.Application.GetApplicationUser(
		context.Background(),
		d.Get("app_id").(string),
		d.Get("user_id").(string),
		nil,
	)

	return g != nil, err
}

func resourceAppUserCreate(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	u, _, err := client.Application.AssignUserToApplication(
		context.Background(),
		d.Get("app_id").(string),
		*getAppUser(d),
	)

	if err != nil {
		return err
	}

	d.SetId(u.Id)

	return resourceAppUserRead(d, m)
}

func resourceAppUserUpdate(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	_, _, err := client.Application.UpdateApplicationUser(
		context.Background(),
		d.Get("app_id").(string),
		d.Get("user_id").(string),
		*getAppUser(d),
	)

	if err != nil {
		return err
	}

	return resourceAppUserRead(d, m)
}

func resourceAppUserRead(d *schema.ResourceData, m interface{}) error {
	u, resp, err := getOktaClientFromMetadata(m).Application.GetApplicationUser(
		context.Background(),
		d.Get("app_id").(string),
		d.Get("user_id").(string),
		nil,
	)

	if resp != nil && is404(resp.StatusCode) {
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	_ = d.Set("profile", u.Profile)
	_ = d.Set("username", u.Credentials.UserName)

	return nil
}

func resourceAppUserDelete(d *schema.ResourceData, m interface{}) error {
	_, err := getOktaClientFromMetadata(m).Application.DeleteApplicationUser(
		context.Background(),
		d.Get("app_id").(string),
		d.Get("user_id").(string),
		nil,
	)
	return err
}

func getAppUser(d *schema.ResourceData) *okta.AppUser {
	var profile interface{}

	rawProfile := d.Get("profile").(string)
	// JSON is already validated
	_ = json.Unmarshal([]byte(rawProfile), &profile)

	return &okta.AppUser{
		Id: d.Get("user_id").(string),
		Credentials: &okta.AppUserCredentials{
			UserName: d.Get("username").(string),
			Password: &okta.AppUserPasswordCredential{
				Value: d.Get("password").(string),
			},
		},
		Profile: profile,
	}
}
