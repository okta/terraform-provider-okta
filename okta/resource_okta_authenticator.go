package okta

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourceAuthenticator() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAuthenticatorCreate,
		ReadContext:   resourceAuthenticatorRead,
		UpdateContext: resourceAuthenticatorUpdate,
		DeleteContext: resourceAuthenticatorDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"key": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "A human-readable string that identifies the Authenticator",
				ValidateDiagFunc: elemInSlice(sdk.AuthenticatorProviders),
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name of the Authenticator",
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
					return new == ""
				},
			},
			"status": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          statusActive,
				ValidateDiagFunc: elemInSlice([]string{statusActive, statusInactive}),
				Description:      "Authenticator status: ACTIVE or INACTIVE",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of Authenticator",
			},
			"provider_hostname": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "localhost",
				Description: "Server host name or IP address",
			},
			"provider_auth_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      9000,
				Description:  "The RADIUS server port (for example 1812). This is defined when the On-Prem RADIUS server is configured",
				RequiredWith: []string{"provider_hostname"},
			},
			"provider_shared_secret": {
				Type:         schema.TypeString,
				Sensitive:    true,
				Optional:     true,
				Description:  "An authentication key that must be defined when the RADIUS server is configured, and must be the same on both the RADIUS client and server.",
				RequiredWith: []string{"provider_hostname"},
			},
			"provider_instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provider_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provider_user_name_template": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "global.assign.userName.login",
				Description:  "Format expected by the provider",
				RequiredWith: []string{"provider_hostname"},
			},
		},
	}
}

// authenticator API is immutable, create is just a read of the key set on the resource
func resourceAuthenticatorCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	authenticator, err := findAuthenticator(ctx, m, "", d.Get("key").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(authenticator.Id)
	status, ok := d.GetOk("status")
	if ok && authenticator.Status != status.(string) {
		if status.(string) == statusInactive {
			_, _, err = getOktaClientFromMetadata(m).Authenticator.DeactivateAuthenticator(ctx, d.Id())
		} else {
			_, _, err = getOktaClientFromMetadata(m).Authenticator.ActivateAuthenticator(ctx, d.Id())
		}
		if err != nil {
			return diag.Errorf("failed to change authenticator status: %v", err)
		}
	}
	return resourceAuthenticatorRead(ctx, d, m)
}

func resourceAuthenticatorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	authenticator, resp, err := getOktaClientFromMetadata(m).Authenticator.GetAuthenticator(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get authenticator: %v", err)
	}
	_ = d.Set("key", authenticator.Key)
	_ = d.Set("name", authenticator.Name)
	_ = d.Set("status", authenticator.Status)
	_ = d.Set("type", authenticator.Type)
	if authenticator.Settings != nil {
		b, _ := json.Marshal(authenticator.Settings)

		dataMap := map[string]interface{}{}
		_ = json.Unmarshal([]byte(string(b)), &dataMap)
		b, _ = json.Marshal(dataMap)

		_ = d.Set("settings", string(b))
	}
	if authenticator.Provider != nil {
		_ = d.Set("provider_type", authenticator.Provider.Type)
		_ = d.Set("provider_hostname", authenticator.Provider.Configuration.HostName)
		_ = d.Set("provider_auth_port", authenticator.Provider.Configuration.AuthPort)
		_ = d.Set("provider_instance_id", authenticator.Provider.Configuration.InstanceId)
		if authenticator.Provider.Configuration.UserNameTemplate != nil {
			_ = d.Set("provider_user_name_template", authenticator.Provider.Configuration.UserNameTemplate.Template)
		}
	}
	return nil
}

func resourceAuthenticatorUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := validateAuthenticator(d)
	if err != nil {
		return diag.FromErr(err)
	}
	_, _, err = getOktaClientFromMetadata(m).Authenticator.UpdateAuthenticator(ctx, d.Id(), *buildAuthenticator(d))
	if err != nil {
		return diag.Errorf("failed to update authenticator: %v", err)
	}
	oldStatus, newStatus := d.GetChange("status")
	if oldStatus != newStatus {
		if newStatus == statusActive {
			_, _, err = getOktaClientFromMetadata(m).Authenticator.ActivateAuthenticator(ctx, d.Id())
		} else {
			_, _, err = getOktaClientFromMetadata(m).Authenticator.DeactivateAuthenticator(ctx, d.Id())
		}
		if err != nil {
			return diag.Errorf("failed to change authenticator status: %v", err)
		}
	}
	return resourceAuthenticatorRead(ctx, d, m)
}

// delete is NOOP, authenticators are immutable for create and delete
func resourceAuthenticatorDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func buildAuthenticator(d *schema.ResourceData) *okta.Authenticator {
	authenticator := okta.Authenticator{
		Type: d.Get("type").(string),
		Id:   d.Id(),
		Key:  d.Get("key").(string),
		Name: d.Get("name").(string),
	}
	if d.Get("type").(string) == "security_key" {
		authenticator.Provider = &okta.AuthenticatorProvider{
			Type: d.Get("provider_type").(string),
			Configuration: &okta.AuthenticatorProviderConfiguration{
				HostName:     d.Get("provider_hostname").(string),
				AuthPort:     d.Get("provider_auth_port").(int64),
				InstanceId:   d.Get("provider_instance_id").(string),
				SharedSecret: d.Get("provider_shared_secret").(string),
				UserNameTemplate: &okta.AuthenticatorProviderConfigurationUserNamePlate{
					Template: "",
				},
			},
		}
	} else {
		var settings okta.AuthenticatorSettings
		if s, ok := d.GetOk("settings"); ok {
			_ = json.Unmarshal([]byte(s.(string)), &settings)
		}
		authenticator.Settings = &settings
	}
	return &authenticator
}

func validateAuthenticator(d *schema.ResourceData) error {
	typ := d.Get("type").(string)
	if typ != "security_key" {
		return nil
	}
	h := d.Get("provider_hostname").(string)
	_, pok := d.GetOk("provider_auth_port")
	s := d.Get("provider_shared_secret").(string)
	templ := d.Get("provider_user_name_template").(string)
	if h == "" || s == "" || templ == "" || !pok {
		return fmt.Errorf("for authenticator type '%s' fields 'provider_hostname', "+
			"'provider_auth_port', 'provider_shared_secret' and 'provider_user_name_template' are required", typ)
	}
	return nil
}
