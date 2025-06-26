package idaas

import (
	"context"
	"net/http"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourceUserType() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserTypeCreate,
		ReadContext:   resourceUserTypeRead,
		UpdateContext: resourceUserTypeUpdate,
		DeleteContext: resourceUserTypeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				return []*schema.ResourceData{d}, nil
			},
		},
		Description: "Creates a User type. This resource allows you to create and configure a User Type.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the User Type.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display Name of the User Type.",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Description of the User Type.",
			},
		},
	}
}

func resourceUserTypeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	newUserType, _, err := getOktaClientFromMetadata(meta).UserType.CreateUserType(ctx, buildUserType(d))
	if err != nil {
		return diag.Errorf("failed to create user type: %v", err)
	}
	d.SetId(newUserType.Id)
	return resourceUserTypeRead(ctx, d, meta)
}

func resourceUserTypeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	userType := buildUserType(d)
	_, _, err := getOktaClientFromMetadata(meta).UserType.UpdateUserType(ctx, d.Id(), userType)
	if err != nil {
		return diag.Errorf("failed to update user type: %v", err)
	}
	return resourceUserTypeRead(ctx, d, meta)
}

func resourceUserTypeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	userType, resp, err := getOktaClientFromMetadata(meta).UserType.GetUserType(ctx, d.Id())
	if err := utils.SuppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get user type: %v", err)
	}
	if userType == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("name", userType.Name)
	_ = d.Set("display_name", userType.DisplayName)
	_ = d.Set("description", userType.Description)
	return nil
}

func resourceUserTypeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// deleting a user type quickly (like in an acceptance test, not typical of
	// "real" usage) can result in a 500 from Okta service. This is an eventual
	// consistency issue on its part and will resolve in under 30 seconds run
	// delete in a backoff.
	boc := utils.NewExponentialBackOffWithContext(ctx, 30*time.Second)
	err := backoff.Retry(func() error {
		resp, err := getOktaClientFromMetadata(meta).UserType.DeleteUserType(ctx, d.Id())
		if err != nil {
			if resp.StatusCode == http.StatusInternalServerError {
				return err
			}
			return backoff.Permanent(err)
		}
		return nil
	}, boc)
	if err != nil {
		return diag.Errorf("failed to delete user type: %v", err)
	}
	return nil
}

func buildUserType(d *schema.ResourceData) sdk.UserType {
	return sdk.UserType{
		Name:        d.Get("name").(string),
		DisplayName: d.Get("display_name").(string),
		Description: d.Get("description").(string),
	}
}
