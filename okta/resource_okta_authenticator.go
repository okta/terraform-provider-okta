package okta

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func resourceAuthenticator() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAuthenticatorCreate,
		ReadContext:   resourceAuthenticatorRead,
		UpdateContext: resourceAuthenticatorUpdate,
		DeleteContext: resourceAuthenticatorDelete,
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return true
				},
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return true
				},
			},
			"settings": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Authenticator settings in JSON format",
				ValidateDiagFunc: stringIsJSON,
				StateFunc:        normalizeDataJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// TODO implement settings diff when we are able to update settings.
					return true
				},
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if new == "" {
						return true
					}
					return old == new
				},
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of the authenticator. When specified in the terraform resource, will act as a filter when searching for the authenticator",
			},
		},
	}
}

func resourceAuthenticatorCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// authenticator API is immutable, create is just a read of the type set on the resource
	return resourceAuthenticatorRead(ctx, d, m)
}

func resourceAuthenticatorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	_type, ok := d.GetOk("type")
	if !ok {
		return diag.Errorf("`type` not present on authenticator resource")
	}
	return findDataSourceAuthenticator(ctx, _type.(string), d, m)
}

func resourceAuthenticatorUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// TODO handle updating settings when the okta-sdk-golang package adds
	// support for updating settings.

	_type, ok := d.GetOk("type")
	if !ok {
		return diag.Errorf("`type` not present on authenticator resource")
	}
	_status, ok := d.GetOk("status")
	if !ok {
		return diag.Errorf("`status` not present on authenticator resource")
	}

	authenticator, err := findOktaAuthenticator(ctx, _type.(string), m)
	if err != nil {
		return diag.Errorf("error finding authenticator resource: %+v", err)
	}

	// NOOP if current authenticator status is the same as the resource
	status := _status.(string)
	if status == authenticator.Status {
		setAuthenticatorValuesOnResourceData(authenticator, d)
		return nil
	}

	switch status {
	case "ACTIVE":
		authenticator, _, err = getOktaClientFromMetadata(m).Authenticator.ActivateAuthenticator(ctx, authenticator.Id)
		if err != nil {
			return diag.Errorf("error activating authenticator resource: %+v", err)
		}
	case "INACTIVE":
		authenticator, _, err = getOktaClientFromMetadata(m).Authenticator.DeactivateAuthenticator(ctx, authenticator.Id)
		if err != nil {
			return diag.Errorf("error deactivating authenticator resource: %+v", err)
		}
	default:
		return diag.Errorf("`status=%q` is invalid for resource", status)
	}

	setAuthenticatorValuesOnResourceData(authenticator, d)

	return nil
}

func resourceAuthenticatorDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// delete is NOOP, authenticators are immutable for create and delete
	return nil
}

func setAuthenticatorValuesOnResourceData(authenticator *okta.Authenticator, d *schema.ResourceData) {
	d.SetId(authenticator.Id)
	_ = d.Set("key", authenticator.Key)
	_ = d.Set("name", authenticator.Name)
	b, err := json.Marshal(authenticator.Settings)
	if err != nil {
		_ = d.Set("settings", string(b))
	}
	_ = d.Set("status", authenticator.Status)
	_ = d.Set("type", authenticator.Type)
}
