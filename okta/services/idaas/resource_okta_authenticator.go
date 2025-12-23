package idaas

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/utils"
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
		ValidateRawResourceConfigFuncs: []schema.ValidateRawResourceConfigFunc{
			func(ctx context.Context, req schema.ValidateResourceConfigFuncRequest, resp *schema.ValidateResourceConfigFuncResponse) {
				keyAttrExists := !req.RawConfig.GetAttr("key").IsNull()
				legacyIgnoreNameAttrExists := !req.RawConfig.GetAttr("legacy_ignore_name").IsNull()
				if keyAttrExists && req.RawConfig.GetAttr("key").AsString() == "custom_app" {
					if !legacyIgnoreNameAttrExists || legacyIgnoreNameAttrExists && req.RawConfig.GetAttr("legacy_ignore_name").True() {
						resp.Diagnostics = append(resp.Diagnostics, diag.Errorf("legacy_ignore_name must be false when creating a custom_app type authenticator")...)
					}
				}
			},
		},

		Description: `~> **WARNING:** This feature is only available as a part of the Identity Engine. [Contact support](mailto:dev-inquiries@okta.com) for further information.

This resource allows you to configure different authenticators.

-> **Create:** The Okta API has an odd notion of create for authenticators. If
the authenticator doesn't exist then a one time 'POST /api/v1/authenticators' to
create the authenticator (hard create) will be performed. Thereafter, that
authenticator is never deleted, it is only deactivated (soft delete). Therefore,
if the authenticator already exists create is just a soft import of an existing
authenticator. This does not apply to custom_otp authenticator. There can be 
multiple custom_otp authenticator. To create new custom_otp authenticator, 
name and key = custom_otp is required. If an old name is used, it will simply 
reactivate the old custom_otp authenticator

-> **Update:** custom_otp authenticator cannot be updated

-> **Delete:** Authenticators can not be truly deleted therefore delete is soft.
Delete will attempt to deativate the authenticator. An authenticator can only be
deactivated if it's not in use by any other policy.`,
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "A human-readable string that identifies the authenticator. Some authenticators are available by feature flag on the organization. Possible values inclue: `duo`, `external_idp`, `google_otp`, `okta_email`, `okta_password`, `okta_verify`, `onprem_mfa`, `phone_number`, `rsa_token`, `security_question`, `webauthn`",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name of the Authenticator",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("legacy_ignore_name").(bool)
				},
			},
			"settings": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Settings for the authenticator. The settings JSON contains values based on Authenticator key. It is not used for authenticators with type `security_key`",
				ValidateDiagFunc: stringIsJSON,
				StateFunc:        utils.NormalizeDataJSON,
				DiffSuppressFunc: utils.NoChangeInObjectWithSortedSlicesFromUnmarshaledJSON,
			},
			"provider_json": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      `Provider JSON allows for expressive providervalues. This argument conflicts with the other 'provider_xxx' arguments. The [CreateProvider](https://developer.okta.com/docs/reference/api/authenticators-admin/#request) illustrates detailed provider values for a Duo authenticator. [Provider values](https://developer.okta.com/docs/reference/api/authenticators-admin/#authenticators-administration-api-object)are listed in Okta API.`,
				ValidateDiagFunc: stringIsJSON,
				StateFunc:        utils.NormalizeDataJSON,
				DiffSuppressFunc: utils.NoChangeInObjectFromUnmarshaledJSON,
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
				Default:     StatusActive,
				Description: "Authenticator status: `ACTIVE` or `INACTIVE`. Default: `ACTIVE`",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "he type of Authenticator. Values include: `password`, `security_question`, `phone`, `email`, `app`, `federated`, and `security_key`.",
			},
			// General Provider Arguments
			"provider_auth_port": {
				Type:          schema.TypeInt,
				Optional:      true,
				Description:   "The RADIUS server port (for example 1812). This is defined when the On-Prem RADIUS server is configured. Used only for authenticators with type `security_key`.  Conflicts with `provider_json` argument.",
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
				Description:   "Server host name or IP address. Default is `localhost`. Used only for authenticators with type `security_key`. Conflicts with `provider_json` argument.",
				ConflictsWith: []string{"provider_json"},
			},
			"provider_shared_secret": {
				Type:          schema.TypeString,
				Sensitive:     true,
				Optional:      true,
				Description:   "An authentication key that must be defined when the RADIUS server is configured, and must be the same on both the RADIUS client and server. Used only for authenticators with type `security_key`. Conflicts with `provider_json` argument.",
				RequiredWith:  []string{"provider_hostname"},
				ConflictsWith: []string{"provider_json"},
			},
			"provider_user_name_template": {
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "global.assign.userName.login",
				Description:   "Username template expected by the provider. Used only for authenticators with type `security_key`.  Conflicts with `provider_json` argument.",
				RequiredWith:  []string{"provider_hostname"},
				ConflictsWith: []string{"provider_json"},
			},
			// DUO specific provider arguments
			"provider_host": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "(DUO specific) - The Duo Security API hostname. Conflicts with `provider_json` argument.",
				ConflictsWith: []string{"provider_json"},
			},
			"provider_integration_key": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "(DUO specific) - The Duo Security integration key.  Conflicts with `provider_json` argument.",
				ConflictsWith: []string{"provider_json"},
			},
			"provider_secret_key": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "(DUO specific) - The Duo Security secret key.  Conflicts with `provider_json` argument.",
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
			"legacy_ignore_name": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Name does not trigger change detection (legacy behavior)",
			},
			"agree_to_terms": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "A value of true indicates that the administrator accepts the terms for creating a new authenticator. Okta requires that you accept the terms when creating a new custom_app authenticator. Other authenticators don't require this field.",
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
func resourceAuthenticatorCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if providerIsClassicOrg(ctx, meta) {
		return resourceOIEOnlyFeatureError(resources.OktaIDaaSAuthenticator)
	}

	var err error
	// soft create if the authenticator already exists
	authenticator, _ := findAuthenticator(ctx, meta, d.Get("name").(string), d.Get("key").(string))
	if authenticator == nil {
		// otherwise hard create
		authenticator, err = buildAuthenticator(d)
		if err != nil {
			return diag.FromErr(err)
		}
		activate := (d.Get("status").(string) == StatusActive)
		qp := &query.Params{
			Activate: utils.BoolPtr(activate),
		}
		authenticator, _, err = getOktaClientFromMetadata(meta).Authenticator.CreateAuthenticator(ctx, *authenticator, qp)
		if err != nil {
			return diag.FromErr(err)
		}
		if d.Get("key").(string) == "custom_otp" {
			var otp *sdk.OTP
			otp, err = buildOTP(d)
			if err != nil {
				return diag.FromErr(err)
			}
			_, err = getOktaClientFromMetadata(meta).Authenticator.SetSettingsOTP(ctx, *otp, authenticator.Id)
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
		if status.(string) == StatusInactive {
			authenticator, _, err = getOktaClientFromMetadata(meta).Authenticator.DeactivateAuthenticator(ctx, d.Id())
		} else {
			authenticator, _, err = getOktaClientFromMetadata(meta).Authenticator.ActivateAuthenticator(ctx, d.Id())
		}
		if err != nil {
			return diag.Errorf("failed to change authenticator status: %v", err)
		}
	}

	establishAuthenticator(authenticator, d)
	return nil
}

func resourceAuthenticatorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if providerIsClassicOrg(ctx, meta) {
		return resourceOIEOnlyFeatureError(resources.OktaIDaaSAuthenticator)
	}

	authenticator, _, err := getOktaClientFromMetadata(meta).Authenticator.GetAuthenticator(ctx, d.Id())
	if err != nil {
		return diag.Errorf("failed to get authenticator: %v", err)
	}
	establishAuthenticator(authenticator, d)

	return nil
}

func resourceAuthenticatorUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if providerIsClassicOrg(ctx, meta) {
		return resourceOIEOnlyFeatureError(resources.OktaIDaaSAuthenticator)
	}

	err := validateAuthenticator(d)
	if err != nil {
		return diag.FromErr(err)
	}
	authenticator, err := buildAuthenticator(d)
	if err != nil {
		return diag.Errorf("failed to update authenticator: %v", err)
	}
	_, _, err = getOktaClientFromMetadata(meta).Authenticator.UpdateAuthenticator(ctx, d.Id(), *authenticator)
	if err != nil {
		return diag.Errorf("failed to update authenticator: %v", err)
	}
	oldStatus, newStatus := d.GetChange("status")
	if oldStatus != newStatus {
		if newStatus == StatusActive {
			_, _, err = getOktaClientFromMetadata(meta).Authenticator.ActivateAuthenticator(ctx, d.Id())
		} else {
			_, _, err = getOktaClientFromMetadata(meta).Authenticator.DeactivateAuthenticator(ctx, d.Id())
		}
		if err != nil {
			return diag.Errorf("failed to change authenticator status: %v", err)
		}
	}
	return resourceAuthenticatorRead(ctx, d, meta)
}

// resourceAuthenticatorDelete Delete is soft, authenticators are immutable for
// true delete. However, deactivate the authenticator as a stand in for delete.
// Authenticators that are utilized by existing policies can not be deactivated.
func resourceAuthenticatorDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if providerIsClassicOrg(ctx, meta) {
		return resourceOIEOnlyFeatureError(resources.OktaIDaaSAuthenticator)
	}

	_, _, err := getOktaClientFromMetadata(meta).Authenticator.DeactivateAuthenticator(ctx, d.Id())
	if err != nil {
		logger(meta).Warn(fmt.Sprintf("Attempted to deactivate authenticator %q as soft delete and received error: %s", d.Get("key"), err))
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
	if d.Get("type").(string) == "security_key" {
		authenticator.Provider = &sdk.AuthenticatorProvider{
			Type: d.Get("provider_type").(string),
			Configuration: &sdk.AuthenticatorProviderConfiguration{
				HostName:     d.Get("provider_hostname").(string),
				AuthPortPtr:  utils.Int64Ptr(d.Get("provider_auth_port").(int)),
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
	} else if d.Get("key").(string) == "custom_app" {
		agreeToTerms, ok := d.Get("agree_to_terms").(bool)
		if !ok {
			return nil, fmt.Errorf("unable to parse agree_to_terms as a boolean value, valid values are true/false")
		}

		authenticator.AgreeToTerms = agreeToTerms
		if s, ok := d.GetOk("settings"); ok {
			var settings sdk.AuthenticatorSettings
			err := json.Unmarshal([]byte(s.(string)), &settings)
			if err != nil {
				return nil, err
			}
			authenticator.Settings = &settings
		}
		authenticator.Provider = &sdk.AuthenticatorProvider{
			Type: d.Get("provider_type").(string),
		}
	} else if d.Get("key").(string) != "custom_otp" { // does not include custom_app
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
		if d.Get("key").(string) != "custom_otp" {
			h := d.Get("provider_hostname").(string)
			_, pok := d.GetOk("provider_auth_port")
			s := d.Get("provider_shared_secret").(string)
			templ := d.Get("provider_user_name_template").(string)
			if h == "" || s == "" || templ == "" || !pok {
				return fmt.Errorf("for authenticator type '%s' fields 'provider_hostname', "+
					"'provider_auth_port', 'provider_shared_secret' and 'provider_user_name_template' are required", typ)
			}
		} else {
			return fmt.Errorf("custom_otp is not updatable")
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
				_ = d.Set("provider_auth_port", authenticator.Provider.Configuration.AuthPortPtr)
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
