package okta

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func resourceAppUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppUserCreate,
		ReadContext:   resourceAppUserRead,
		UpdateContext: resourceAppUserUpdate,
		DeleteContext: resourceAppUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				parts := strings.Split(d.Id(), "/")
				if len(parts) != 2 {
					return nil, errors.New("invalid resource import specifier. Use: terraform import <app_id>/<user_id>")
				}

				_ = d.Set("app_id", parts[0])
				_ = d.Set("user_id", parts[1])
				_ = d.Set("retain_assignment", false)

				assignment, _, err := getOktaClientFromMetadata(m).Application.
					GetApplicationUser(ctx, parts[0], parts[1], nil)
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
				Type:             schema.TypeString,
				ValidateDiagFunc: stringIsJSON,
				StateFunc:        normalizeDataJSON,
				Optional:         true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == ""
				},
			},
			"retain_assignment": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Retain the user assignment on destroy. If set to true, the resource will be removed from state but not from the Okta app.",
			},
		},
	}
}

func resourceAppUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	u, _, err := getOktaClientFromMetadata(m).Application.AssignUserToApplication(
		ctx,
		d.Get("app_id").(string),
		*getAppUser(d),
	)
	if err != nil {
		return diag.Errorf("failed to assign user to application: %v", err)
	}
	d.SetId(u.Id)
	return resourceAppUserRead(ctx, d, m)
}

func resourceAppUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	_, _, err := getOktaClientFromMetadata(m).Application.UpdateApplicationUser(
		ctx,
		d.Get("app_id").(string),
		d.Get("user_id").(string),
		*getAppUser(d),
	)
	if err != nil {
		return diag.Errorf("failed to update application's user: %v", err)
	}
	return resourceAppUserRead(ctx, d, m)
}

func resourceAppUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	u, resp, err := getOktaClientFromMetadata(m).Application.GetApplicationUser(
		ctx,
		d.Get("app_id").(string),
		d.Get("user_id").(string),
		nil,
	)
	if is404(resp) {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get application's user: %v", err)
	}
	var rawProfile string
	if u.Profile != nil {
		p, _ := json.Marshal(u.Profile)
		rawProfile = string(p)
	}
	_ = d.Set("profile", rawProfile)
	_ = d.Set("username", u.Credentials.UserName)
	return nil
}

func resourceAppUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	retain := d.Get("retain_assignment").(bool)
	if retain {
		// The assignment should be retained, bail before DeleteApplicationUser is called
		return nil
	}

	_, err := getOktaClientFromMetadata(m).Application.DeleteApplicationUser(
		ctx,
		d.Get("app_id").(string),
		d.Get("user_id").(string),
		nil,
	)
	if err != nil {
		return diag.Errorf("failed to delete application's user: %v", err)
	}
	return nil
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
