package okta

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/okta/terraform-provider-okta/sdk/query"
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
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "A human-readable string that identifies the Authenticator",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name of the Authenticator",
			},
			"settings": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Authenticator settings in JSON format",
				ValidateDiagFunc: stringIsJSON,
				StateFunc:        normalizeDataJSON,
				DiffSuppressFunc: noChangeInObjectFromUnmarshaledJSON,
			},
			"provider_json": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Provider in JSON format",
				ValidateDiagFunc: stringIsJSON,
				StateFunc:        normalizeDataJSON,
				DiffSuppressFunc: noChangeInObjectFromUnmarshaledJSON,
				ConflictsWith: []string{
					// general
					"provider_auth_port",
					"provider_hostname",
					"provider_shared_secret",
					"provider_user_name_template",
					// duo
					"provider_host",
					"provider_integration_key",
					"provider_secret_key",
				},
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     statusActive,
				Description: "Authenticator status: ACTIVE or INACTIVE",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of Authenticator",
			},
			// General Provider Arguments
			"provider_auth_port": {
				Type:          schema.TypeInt,
				Optional:      true,
				Description:   "The RADIUS server port (for example 1812). This is defined when the On-Prem RADIUS server is configured",
				RequiredWith:  []string{"provider_hostname"},
				ConflictsWith: []string{"provider_json"},
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					if _, ok := d.GetOk("provider_json"); ok {
						return true
					}
					return false
				},
			},
			"provider_hostname": {
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "localhost",
				Description:   "Server host name or IP address",
				ConflictsWith: []string{"provider_json"},
			},
			"provider_shared_secret": {
				Type:          schema.TypeString,
				Sensitive:     true,
				Optional:      true,
				Description:   "An authentication key that must be defined when the RADIUS server is configured, and must be the same on both the RADIUS client and server.",
				RequiredWith:  []string{"provider_hostname"},
				ConflictsWith: []string{"provider_json"},
			},
			"provider_user_name_template": {
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "global.assign.userName.login",
				Description:   "Format expected by the provider",
				RequiredWith:  []string{"provider_hostname"},
				ConflictsWith: []string{"provider_json"},
			},
			// DUO specific provider arguments
			"provider_host": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "The Duo Security API hostname",
				ConflictsWith: []string{"provider_json"},
			},
			"provider_integration_key": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "The Duo Security integration key",
				ConflictsWith: []string{"provider_json"},
			},
			"provider_secret_key": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "The Duo Security secret key",
				ConflictsWith: []string{"provider_json"},
			},
			// General Provider Attributes
			"provider_instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "App Instance ID.",
			},
			"provider_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Provider type. Supported value for Duo: `DUO`. Supported value for Custom App: `PUSH`",
			},
		},
	}
}

// resourceAuthenticatorCreate Okta API has an odd notion of create for
// authenticators. If the authenticator doesn't exist then a one time `POST
// /api/v1/authenticators` to create the authenticator (hard create) is to be
// performed. Thereafter, that authenticator is never deleted, it is only
// deactivated (soft delete). Therefore, if the authenticator already exists
// create is just a soft import of an existing authenticator.
func resourceAuthenticatorCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if isClassicOrg(ctx, m) {
		return resourceOIEOnlyFeatureError(authenticator)
	}

	var err error
	// soft create if the authenticator already exists
	authenticator, _ := findAuthenticator(ctx, m, d.Get("name").(string), d.Get("key").(string))
	if authenticator == nil {
		// otherwise hard create
		authenticator, err = buildAuthenticator(d)
		if err != nil {
			return diag.FromErr(err)
		}
		activate := (d.Get("status").(string) == statusActive)
		qp := &query.Params{
			Activate: boolPtr(activate),
		}
		if(d.Get("key").(string) == "custom_otp"){
			qp = &query.Params{
				Activate: boolPtr(false),
			}
		}
		authenticator, _, err = getOktaClientFromMetadata(m).Authenticator.CreateAuthenticator(ctx, *authenticator, qp)
		if err != nil {
			return diag.FromErr(err)
		}
		if(d.Get("key").(string) == "custom_otp"){
			var otp *sdk.OTP
			otp, err = buildOTP(d)
			if err != nil {
				return diag.FromErr(err)
			}
			_, err = getOktaClientFromMetadata(m).Authenticator.SetSettingsOTP(ctx, *otp, authenticator.Id)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	d.SetId(authenticator.Id)

	// If status is defined in the config, and the actual status reported by the
	// API is not the same, then toggle the status. Soft update.
	status, ok := d.GetOk("status")
	if ok && authenticator.Status != status.(string) {
		var err error
		if status.(string) == statusInactive {
			authenticator, _, err = getOktaClientFromMetadata(m).Authenticator.DeactivateAuthenticator(ctx, d.Id())
		} else {
			if(d.Get("key").(string) == "custom_otp"){
				var otp *sdk.OTP
				otp, err = buildOTP(d)
				if err != nil {
					return diag.FromErr(err)
				}
				_, err = getOktaClientFromMetadata(m).Authenticator.SetSettingsOTP(ctx, *otp, d.Id())
				if err != nil {
					return diag.FromErr(err)
				}
			}
			authenticator, _, err = getOktaClientFromMetadata(m).Authenticator.ActivateAuthenticator(ctx, d.Id())
		}
		if err != nil {
			return diag.Errorf("failed to change authenticator status: %v", err)
		}
	}

	establishAuthenticator(authenticator, d)
	return nil
}

func resourceAuthenticatorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if isClassicOrg(ctx, m) {
		return resourceOIEOnlyFeatureError(authenticator)
	}

	authenticator, _, err := getOktaClientFromMetadata(m).Authenticator.GetAuthenticator(ctx, d.Id())
	if err != nil {
		return diag.Errorf("failed to get authenticator: %v", err)
	}
	establishAuthenticator(authenticator, d)

	return nil
}

func resourceAuthenticatorUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if isClassicOrg(ctx, m) {
		return resourceOIEOnlyFeatureError(authenticator)
	}

	err := validateAuthenticator(d)
	if err != nil {
		return diag.FromErr(err)
	}
	authenticator, err := buildAuthenticator(d)
	if err != nil {
		return diag.Errorf("failed to update authenticator: %v", err)
	}
	_, _, err = getOktaClientFromMetadata(m).Authenticator.UpdateAuthenticator(ctx, d.Id(), *authenticator)
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

// resourceAuthenticatorDelete Delete is soft, authenticators are immutable for
// true delete. However, deactivate the authenticator as a stand in for delete.
// Authenticators that are utilized by existing policies can not be deactivated.
func resourceAuthenticatorDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if isClassicOrg(ctx, m) {
		return resourceOIEOnlyFeatureError(authenticator)
	}

	_, _, err := getOktaClientFromMetadata(m).Authenticator.DeactivateAuthenticator(ctx, d.Id())
	if err != nil {
		logger(m).Warn(fmt.Sprintf("Attempted to deactivate authenticator %q as soft delete and received error: %s", d.Get("key"), err))
	}

	return nil
}

func buildAuthenticator(d *schema.ResourceData) (*sdk.Authenticator, error) {
	authenticator := sdk.Authenticator{
		Type: d.Get("type").(string),
		Id:   d.Id(),
		Key:  d.Get("key").(string),
		Name: d.Get("name").(string),
	}
	if d.Get("key").(string) == "custom_otp" {

	} else if d.Get("type").(string) == "security_key" {
		authenticator.Provider = &sdk.AuthenticatorProvider{
			Type: d.Get("provider_type").(string),
			Configuration: &sdk.AuthenticatorProviderConfiguration{
				HostName:     d.Get("provider_hostname").(string),
				AuthPortPtr:  int64Ptr(d.Get("provider_auth_port").(int)),
				InstanceId:   d.Get("provider_instance_id").(string),
				SharedSecret: d.Get("provider_shared_secret").(string),
				UserNameTemplate: &sdk.AuthenticatorProviderConfigurationUserNamePlate{
					Template: "",
				},
			},
		}
	} else if d.Get("type").(string) == "DUO" {
		authenticator.Provider = &sdk.AuthenticatorProvider{
			Type: d.Get("provider_type").(string),
			Configuration: &sdk.AuthenticatorProviderConfiguration{
				Host:           d.Get("provider_host").(string),
				SecretKey:      d.Get("provider_secret_key").(string),
				IntegrationKey: d.Get("provider_integration_key").(string),
				UserNameTemplate: &sdk.AuthenticatorProviderConfigurationUserNamePlate{
					Template: d.Get("provider_user_name_template").(string),
				},
			},
		}
	} else {
		if s, ok := d.GetOk("settings"); ok {
			var settings sdk.AuthenticatorSettings
			err := json.Unmarshal([]byte(s.(string)), &settings)
			if err != nil {
				return nil, err
			}
			authenticator.Settings = &settings
		}
	}

	if p, ok := d.GetOk("provider_json"); ok {
		var provider sdk.AuthenticatorProvider
		err := json.Unmarshal([]byte(p.(string)), &provider)
		if err != nil {
			return nil, err
		}
		authenticator.Provider = &provider
	}

	return &authenticator, nil
}

func buildOTP(d *schema.ResourceData) (*sdk.OTP, error) {
	otp := sdk.OTP{}
	if s, ok := d.GetOk("settings"); ok {
		var settings sdk.AuthenticatorSettingsOTP
		err := json.Unmarshal([]byte(s.(string)), &settings)
		if err != nil {
			return nil, err
		}
		otp.Settings = &settings
	}

	return &otp, nil
}

func validateAuthenticator(d *schema.ResourceData) error {
	typ := d.Get("type").(string)
	if typ == "security_key" {
		h := d.Get("provider_hostname").(string)
		_, pok := d.GetOk("provider_auth_port")
		s := d.Get("provider_shared_secret").(string)
		templ := d.Get("provider_user_name_template").(string)
		if h == "" || s == "" || templ == "" || !pok {
			return fmt.Errorf("for authenticator type '%s' fields 'provider_hostname', "+
				"'provider_auth_port', 'provider_shared_secret' and 'provider_user_name_template' are required", typ)
		}
	}

	typ = d.Get("provider_type").(string)
	if typ == "DUO" {
		h := d.Get("provider_host").(string)
		sk := d.Get("provider_secret_key").(string)
		ik := d.Get("provider_integration_key").(string)
		templ := d.Get("provider_user_name_template").(string)
		if h == "" || sk == "" || ik == "" || templ == "" {
			return fmt.Errorf("for authenticator type '%s' fields 'provider_host', "+
				"'provider_secret_key', 'provider_integration_key' and 'provider_user_name_template' are required", typ)
		}
	}
	return nil
}

func establishAuthenticator(authenticator *sdk.Authenticator, d *schema.ResourceData) {
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

		if authenticator.Type == "security_key" {
			_ = d.Set("provider_hostname", authenticator.Provider.Configuration.HostName)
			if authenticator.Provider.Configuration.AuthPortPtr != nil {
				_ = d.Set("provider_auth_port", *authenticator.Provider.Configuration.AuthPortPtr)
			}
			_ = d.Set("provider_instance_id", authenticator.Provider.Configuration.InstanceId)
		}

		if authenticator.Provider.Configuration.UserNameTemplate != nil {
			_ = d.Set("provider_user_name_template", authenticator.Provider.Configuration.UserNameTemplate.Template)
		}

		if authenticator.Provider.Type == "DUO" {
			_ = d.Set("provider_host", authenticator.Provider.Configuration.Host)
			_ = d.Set("provider_secret_key", authenticator.Provider.Configuration.SecretKey)
			_ = d.Set("provider_integration_key", authenticator.Provider.Configuration.IntegrationKey)
		}
	}
}
