package idaas

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	v6okta "github.com/okta/okta-sdk-golang/v6/okta"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/utils"
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
				if keyAttrExists && req.RawConfig.GetAttr("key").IsKnown() && req.RawConfig.GetAttr("key").AsString() == "custom_app" {
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
deactivated if it's not in use by any other policy.

-> **Temporary Access Code (TAC):** The TAC authenticator (key = 'tac') is
configured via the 'provider_json' argument. The provider JSON must contain
'type' (value: 'TAC', uppercase) and 'configuration' fields. The configuration
supports: 'minTtl', 'maxTtl', 'defaultTtl' (minutes), 'length' (code length),
'complexity' (object with 'numbers', 'letters', 'specialCharacters' booleans),
and 'multiUseAllowed' (boolean).

-> All authenticator operations use the Okta SDK v6.`,
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "A human-readable string that identifies the authenticator. Some authenticators are available by feature flag on the organization. Possible values include: `duo`, `external_idp`, `google_otp`, `okta_email`, `okta_password`, `okta_verify`, `onprem_mfa`, `phone_number`, `rsa_token`, `security_question`, `tac`, `webauthn`",
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
				Description:      "Settings for the authenticator. The settings JSON contains values based on Authenticator key. It is not used for authenticators with type `security_key` or key `tac`",
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

	v6Client := getOktaV6ClientFromMetadata(meta)

	// soft create if the authenticator already exists
	authenticator, _ := findAuthenticatorV6(ctx, v6Client, d.Get("name").(string), d.Get("key").(string))
	if authenticator == nil {
		// otherwise hard create
		auth, err := buildAuthenticatorV6(d)
		if err != nil {
			return diag.FromErr(err)
		}
		activate := d.Get("status").(string) == StatusActive
		authenticator, _, err = v6Client.AuthenticatorAPI.
			CreateAuthenticator(ctx).
			Authenticator(auth).
			Activate(activate).
			Execute()
		if err != nil {
			return diag.FromErr(authenticatorAPIError("create", err))
		}

		if d.Get("key").(string) == "custom_otp" {
			if err := setOTPMethodSettings(ctx, v6Client, authenticator.GetId(), d); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	d.SetId(authenticator.GetId())

	// If status is defined in the config, and the actual status reported by the
	// API is not the same, then toggle the status. Soft update.
	status, ok := d.GetOk("status")
	if ok && authenticator.GetStatus() != status.(string) {
		var err error
		if status.(string) == StatusInactive {
			authenticator, _, err = v6Client.AuthenticatorAPI.DeactivateAuthenticator(ctx, d.Id()).Execute()
		} else {
			authenticator, _, err = v6Client.AuthenticatorAPI.ActivateAuthenticator(ctx, d.Id()).Execute()
		}
		if err != nil {
			return diag.Errorf("failed to change authenticator status: %v", err)
		}
	}

	establishAuthenticatorV6(authenticator, d)
	return nil
}

func resourceAuthenticatorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if providerIsClassicOrg(ctx, meta) {
		return resourceOIEOnlyFeatureError(resources.OktaIDaaSAuthenticator)
	}

	authenticator, _, err := getOktaV6ClientFromMetadata(meta).AuthenticatorAPI.GetAuthenticator(ctx, d.Id()).Execute()
	if err != nil {
		return diag.Errorf("failed to get authenticator: %v", err)
	}
	establishAuthenticatorV6(authenticator, d)

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
	auth, err := buildAuthenticatorV6(d)
	if err != nil {
		return diag.Errorf("failed to update authenticator: %v", err)
	}
	v6Client := getOktaV6ClientFromMetadata(meta)
	_, _, err = v6Client.AuthenticatorAPI.ReplaceAuthenticator(ctx, d.Id()).Authenticator(auth).Execute()
	if err != nil {
		return diag.Errorf("failed to update authenticator: %v", err)
	}
	oldStatus, newStatus := d.GetChange("status")
	if oldStatus != newStatus {
		if newStatus == StatusActive {
			_, _, err = v6Client.AuthenticatorAPI.ActivateAuthenticator(ctx, d.Id()).Execute()
		} else {
			_, _, err = v6Client.AuthenticatorAPI.DeactivateAuthenticator(ctx, d.Id()).Execute()
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

	_, _, err := getOktaV6ClientFromMetadata(meta).AuthenticatorAPI.DeactivateAuthenticator(ctx, d.Id()).Execute()
	if err != nil {
		logger(meta).Warn(fmt.Sprintf("Attempted to deactivate authenticator %q as soft delete and received error: %s", d.Get("key"), err))
	}

	return nil
}

// authenticatorProviderV6 and authenticatorProviderConfigV6 mirror the JSON
// shape of an authenticator's `provider` object. The v6 SDK's AuthenticatorBase
// does not model `provider` as a typed field (it is only typed on the
// ListAuthenticators discriminated-union response, which isn't used for
// create/replace/get) so these local types exist purely to get correct
// `omitempty` wire behavior when building the provider payload for the
// security_key/RADIUS and DUO cases below; provider_json is handled
// generically without any typed struct.
type authenticatorProviderV6 struct {
	Type          string                         `json:"type,omitempty"`
	Configuration *authenticatorProviderConfigV6 `json:"configuration,omitempty"`
}

type authenticatorProviderConfigV6 struct {
	HostName         string                                   `json:"hostName,omitempty"`
	AuthPort         *int64                                   `json:"authPort,omitempty"`
	InstanceId       string                                   `json:"instanceId,omitempty"`
	SharedSecret     string                                   `json:"sharedSecret,omitempty"`
	Host             string                                   `json:"host,omitempty"`
	SecretKey        string                                   `json:"secretKey,omitempty"`
	IntegrationKey   string                                   `json:"integrationKey,omitempty"`
	UserNameTemplate *authenticatorProviderUserNameTemplateV6 `json:"userNameTemplate,omitempty"`
}

type authenticatorProviderUserNameTemplateV6 struct {
	Template string `json:"template,omitempty"`
}

// authenticatorSettingsV6 mirrors the fixed field set the legacy SDK's
// AuthenticatorSettings modeled. Settings are round-tripped through this type
// (rather than passed through as an arbitrary map) on both build and establish
// so unrecognized/zero-value fields the Okta API happens to return (e.g.
// oauthClientId, an empty appInstanceId) are dropped exactly as before,
// preserving existing plan/state diff behavior.
type authenticatorSettingsV6 struct {
	AllowedFor              string                         `json:"allowedFor,omitempty"`
	AppInstanceId           string                         `json:"appInstanceId,omitempty"`
	ChannelBinding          *authenticatorChannelBindingV6 `json:"channelBinding,omitempty"`
	Compliance              *authenticatorComplianceV6     `json:"compliance,omitempty"`
	TokenLifetimeInMinutes  *int64                         `json:"tokenLifetimeInMinutes,omitempty"`
	UserVerification        string                         `json:"userVerification,omitempty"`
	EnrollmentSecurityLevel string                         `json:"enrollmentSecurityLevel,omitempty"`
	UserVerificationMethods []string                       `json:"userVerificationMethods,omitempty"`
}

type authenticatorChannelBindingV6 struct {
	Required string `json:"required,omitempty"`
	Style    string `json:"style,omitempty"`
}

type authenticatorComplianceV6 struct {
	Fips string `json:"fips,omitempty"`
}

// marshalAuthenticatorSettings filters a raw `settings` value (as decoded
// generically from AdditionalProperties) through authenticatorSettingsV6,
// matching the legacy behavior of reading settings back through a fixed,
// typed field set.
func marshalAuthenticatorSettings(raw interface{}) (string, bool) {
	b, err := json.Marshal(raw)
	if err != nil {
		return "", false
	}
	var settings authenticatorSettingsV6
	if json.Unmarshal(b, &settings) != nil {
		return "", false
	}
	b, err = json.Marshal(settings)
	if err != nil {
		return "", false
	}
	return string(b), true
}

// buildAuthenticatorV6 constructs a v6 AuthenticatorBase for any authenticator
// key. `provider`, `settings`, and `agreeToTerms` aren't typed fields on
// AuthenticatorBase, so they're carried in AdditionalProperties, which the v6
// SDK serializes at the top level of the request body.
func buildAuthenticatorV6(d *schema.ResourceData) (v6okta.AuthenticatorBase, error) {
	auth := v6okta.AuthenticatorBase{}
	if id := d.Id(); id != "" {
		auth.SetId(id)
	}
	if typ := d.Get("type").(string); typ != "" {
		auth.SetType(typ)
	}
	auth.SetKey(d.Get("key").(string))
	auth.SetName(d.Get("name").(string))

	additionalProps := map[string]interface{}{}

	switch {
	case d.Get("type").(string) == "security_key" && d.Get("key").(string) != "webauthn":
		// WebAuthn is a built-in authenticator and doesn't need provider configuration
		authPort := int64(d.Get("provider_auth_port").(int))
		additionalProps["provider"] = authenticatorProviderV6{
			Type: d.Get("provider_type").(string),
			Configuration: &authenticatorProviderConfigV6{
				HostName:         d.Get("provider_hostname").(string),
				AuthPort:         &authPort,
				InstanceId:       d.Get("provider_instance_id").(string),
				SharedSecret:     d.Get("provider_shared_secret").(string),
				UserNameTemplate: &authenticatorProviderUserNameTemplateV6{},
			},
		}
	case d.Get("type").(string) == "DUO":
		additionalProps["provider"] = authenticatorProviderV6{
			Type: d.Get("provider_type").(string),
			Configuration: &authenticatorProviderConfigV6{
				Host:           d.Get("provider_host").(string),
				SecretKey:      d.Get("provider_secret_key").(string),
				IntegrationKey: d.Get("provider_integration_key").(string),
				UserNameTemplate: &authenticatorProviderUserNameTemplateV6{
					Template: d.Get("provider_user_name_template").(string),
				},
			},
		}
	case d.Get("key").(string) == "custom_app":
		agreeToTerms, ok := d.Get("agree_to_terms").(bool)
		if !ok {
			return auth, fmt.Errorf("unable to parse agree_to_terms as a boolean value, valid values are true/false")
		}
		additionalProps["agreeToTerms"] = agreeToTerms
		if s, ok := d.GetOk("settings"); ok {
			var settings authenticatorSettingsV6
			if err := json.Unmarshal([]byte(s.(string)), &settings); err != nil {
				return auth, err
			}
			additionalProps["settings"] = settings
		}
		additionalProps["provider"] = authenticatorProviderV6{Type: d.Get("provider_type").(string)}
	case d.Get("key").(string) != "custom_otp": // does not include custom_app
		if s, ok := d.GetOk("settings"); ok {
			var settings authenticatorSettingsV6
			if err := json.Unmarshal([]byte(s.(string)), &settings); err != nil {
				return auth, err
			}
			additionalProps["settings"] = settings
		}
	}

	if p, ok := d.GetOk("provider_json"); ok {
		var provider interface{}
		if err := json.Unmarshal([]byte(p.(string)), &provider); err != nil {
			return auth, err
		}
		additionalProps["provider"] = provider
	}

	if len(additionalProps) > 0 {
		auth.AdditionalProperties = additionalProps
	}
	return auth, nil
}

// setOTPMethodSettings configures custom_otp settings via the authenticator
// methods API (`PUT /api/v1/authenticators/{id}/methods/otp`). OTP settings
// aren't part of the base authenticator body - confirmed against Okta's
// management OpenAPI spec, whose replaceAuthenticatorMethod operation accepts
// a flat AuthenticatorMethodOtp body (acceptableAdjacentIntervals, algorithm,
// encoding, passCodeLength, protocol, timeIntervalInSeconds), not a
// settings-wrapped object.
func setOTPMethodSettings(ctx context.Context, v6Client *v6okta.APIClient, authenticatorId string, d *schema.ResourceData) error {
	otp := v6okta.AuthenticatorMethodOtp{}
	if s, ok := d.GetOk("settings"); ok {
		if err := json.Unmarshal([]byte(s.(string)), &otp); err != nil {
			return err
		}
	}
	_, _, err := v6Client.AuthenticatorAPI.
		ReplaceAuthenticatorMethod(ctx, authenticatorId, "otp").
		ListAuthenticatorMethods200ResponseInner(v6okta.AuthenticatorMethodOtpAsListAuthenticatorMethods200ResponseInner(&otp)).
		Execute()
	if err != nil {
		return fmt.Errorf("failed to set custom_otp settings: %v", err)
	}
	return nil
}

func validateAuthenticator(d *schema.ResourceData) error {
	typ := d.Get("type").(string)
	key := d.Get("key").(string)
	if typ == "security_key" {
		// WebAuthn is a built-in authenticator and doesn't need provider configuration
		if key == "webauthn" {
			// WebAuthn authenticators don't require provider fields
			return nil
		}
		if key != "custom_otp" {
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

// findAuthenticatorV6 looks for an existing authenticator by name and/or key
// via the v6 SDK's ListAuthenticators, replicating the legacy findAuthenticator
// matching semantics (custom_app can have multiple authenticators sharing a
// key, so it matches by name only; custom_otp requires an exact name+key match
// on the first candidate found, erroring otherwise; everything else matches by
// name OR key).
func findAuthenticatorV6(ctx context.Context, v6Client *v6okta.APIClient, name, key string) (*v6okta.AuthenticatorBase, error) {
	authenticators, _, err := v6Client.AuthenticatorAPI.ListAuthenticators(ctx).Execute()
	if err != nil {
		return nil, err
	}
	for _, item := range authenticators {
		authenticator, err := normalizeListedAuthenticator(item)
		if err != nil {
			continue
		}
		switch {
		case key == "custom_app":
			if authenticator.GetName() == name { // there can be more than 1 custom_app type authenticator, return nil in the end if we can't find by name.
				return authenticator, nil
			}
		case key != "custom_otp":
			if authenticator.GetName() == name {
				return authenticator, nil
			}
			if authenticator.GetKey() == key {
				return authenticator, nil
			}
		default:
			if authenticator.GetName() == name && authenticator.GetKey() == key {
				return authenticator, nil
			}
			return nil, fmt.Errorf("authenticator with name '%s' and/or key '%s' does not exist", name, key)
		}
	}
	if key != "" {
		return nil, fmt.Errorf("authenticator with key '%s' does not exist", key)
	}
	return nil, fmt.Errorf("authenticator with name '%s' does not exist", name) // authenticator names must be unique.
}

// normalizeListedAuthenticator converts one ListAuthenticators union item
// (whose MarshalJSON serializes whichever typed key-variant is set) back into
// a generic AuthenticatorBase, so callers can read id/key/name/status/provider/
// settings without a type switch over every authenticator key variant.
func normalizeListedAuthenticator(item v6okta.ListAuthenticators200ResponseInner) (*v6okta.AuthenticatorBase, error) {
	b, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}
	var base v6okta.AuthenticatorBase
	if err := json.Unmarshal(b, &base); err != nil {
		return nil, err
	}
	return &base, nil
}

// authenticatorAPIError wraps a v6 SDK error for an authenticator operation,
// surfacing the Okta API response body (errorSummary/errorCauses) when
// available so configuration problems are visible to the user.
func authenticatorAPIError(action string, err error) error {
	var apiErr *v6okta.GenericOpenAPIError
	if errors.As(err, &apiErr) && len(apiErr.Body()) > 0 {
		return fmt.Errorf("failed to %s authenticator: %v: %s", action, err, string(apiErr.Body()))
	}
	return fmt.Errorf("failed to %s authenticator: %v", action, err)
}

// establishAuthenticatorV6 populates resource attributes from a v6
// AuthenticatorBase. `settings`, `provider`, and `agreeToTerms` are read out of
// AdditionalProperties generically (they round-trip as raw JSON) rather than
// through any typed model.
//
// provider_json is populated here for any authenticator whose provider
// configuration isn't otherwise captured by the dedicated provider_* attributes
// below (e.g. `tac`), matching the pre-v6-migration behavior where provider_json
// was write-only for DUO/security_key (those types' current provider values are
// reflected via provider_host/provider_hostname/etc. instead, to avoid a
// perpetual diff against the Optional, non-Computed provider_json attribute).
func establishAuthenticatorV6(authenticator *v6okta.AuthenticatorBase, d *schema.ResourceData) {
	_ = d.Set("key", authenticator.GetKey())
	_ = d.Set("name", authenticator.GetName())
	_ = d.Set("status", authenticator.GetStatus())
	_ = d.Set("type", authenticator.GetType())

	if raw, ok := authenticator.AdditionalProperties["settings"]; ok && raw != nil {
		if s, ok := marshalAuthenticatorSettings(raw); ok {
			_ = d.Set("settings", s)
		}
	}

	raw, ok := authenticator.AdditionalProperties["provider"]
	if !ok || raw == nil {
		return
	}
	b, err := json.Marshal(raw)
	if err != nil {
		return
	}

	var provider authenticatorProviderV6
	if json.Unmarshal(b, &provider) != nil {
		return
	}
	_ = d.Set("provider_type", provider.Type)

	authType := authenticator.GetType()
	authKey := authenticator.GetKey()
	hasDedicatedProviderAttrs := (authType == "security_key" && authKey != "webauthn") || provider.Type == "DUO"
	if !hasDedicatedProviderAttrs {
		_ = d.Set("provider_json", string(b))
	}

	if provider.Configuration == nil {
		return
	}
	if authType == "security_key" && authKey != "webauthn" {
		_ = d.Set("provider_hostname", provider.Configuration.HostName)
		if provider.Configuration.AuthPort != nil {
			_ = d.Set("provider_auth_port", int(*provider.Configuration.AuthPort))
		}
		_ = d.Set("provider_instance_id", provider.Configuration.InstanceId)
	}
	if provider.Configuration.UserNameTemplate != nil {
		_ = d.Set("provider_user_name_template", provider.Configuration.UserNameTemplate.Template)
	}
	if provider.Type == "DUO" {
		_ = d.Set("provider_host", provider.Configuration.Host)
		_ = d.Set("provider_secret_key", provider.Configuration.SecretKey)
		_ = d.Set("provider_integration_key", provider.Configuration.IntegrationKey)
	}
}
