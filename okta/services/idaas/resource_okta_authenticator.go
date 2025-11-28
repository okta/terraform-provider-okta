package idaas

import (
	"context"
	"encoding/json"
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
		Description: `~> **WARNING:** This feature is only available as a part of the Identity Engine. [Contact support](mailto:dev-inquiries@okta.com) for further information.

This resource allows you to configure different authenticators.

-> **Create:** The Okta API has an odd notion of create for authenticators. If
the authenticator doesn't exist then a one time 'POST /api/v1/authenticators' to
create the authenticator (hard create) will be performed. Thereafter, that
authenticator is never deleted, it is only deactivated (soft delete). Therefore,
if the authenticator already exists create is just a soft import of an existing
authenticator. This does not apply to custom_otp authenticator. There can be 
multiple custom_otp authenticator. To create new custom_otp authenticator, a new 
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

	client := getOktaV6ClientFromMetadata(meta)

	// soft create if the authenticator already exists
	authenticator, _ := findAuthenticatorV6(ctx, client, d.Get("name").(string), d.Get("key").(string))
	if authenticator == nil {
		// otherwise hard create
		authenticatorReq, err := buildAuthenticatorV6(d, meta)
		if err != nil {
			return diag.FromErr(err)
		}

		activate := (d.Get("status").(string) == StatusActive)
		req := client.AuthenticatorAPI.CreateAuthenticator(ctx).Authenticator(*authenticatorReq)
		if activate {
			req = req.Activate(activate)
		}

		createdAuth, resp, err := req.Execute()
		if err != nil {
			return diag.Errorf("failed to create authenticator: %v", err)
		}
		defer resp.Body.Close()
		authenticator = createdAuth

		// Handle custom_otp special case
		if d.Get("key").(string) == "custom_otp" {
			if s, ok := d.GetOk("settings"); ok {
				var settingsMap map[string]interface{}
				if err := json.Unmarshal([]byte(s.(string)), &settingsMap); err != nil {
					return diag.FromErr(err)
				}

				// Update OTP method settings using ReplaceAuthenticatorMethod
				methodBytes, _ := json.Marshal(map[string]interface{}{
					"type":     "otp",
					"settings": settingsMap,
				})

				methodReq := client.AuthenticatorAPI.ReplaceAuthenticatorMethod(ctx, authenticator.GetId(), "otp")
				_, methodResp, err := methodReq.Execute()
				if err != nil {
					logger(meta).Warn(fmt.Sprintf("Failed to set OTP settings: %v, request body: %s", err, string(methodBytes)))
				} else if methodResp != nil {
					defer methodResp.Body.Close()
				}
			}
		}
	}

	d.SetId(authenticator.GetId())

	// If status is defined in the config, and the actual status reported by the
	// API is not the same, then toggle the status. Soft update.
	status, ok := d.GetOk("status")
	if ok && authenticator.GetStatus() != status.(string) {
		if status.(string) == StatusInactive {
			authenticator, resp, err := client.AuthenticatorAPI.DeactivateAuthenticator(ctx, d.Id()).Execute()
			if err != nil {
				return diag.Errorf("failed to change authenticator status: %v", err)
			}
			defer resp.Body.Close()
			return establishAuthenticatorV6(authenticator, d, meta)
		} else {
			authenticator, resp, err := client.AuthenticatorAPI.ActivateAuthenticator(ctx, d.Id()).Execute()
			if err != nil {
				return diag.Errorf("failed to change authenticator status: %v", err)
			}
			defer resp.Body.Close()
			return establishAuthenticatorV6(authenticator, d, meta)
		}
	}

	return establishAuthenticatorV6(authenticator, d, meta)
}

func resourceAuthenticatorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if providerIsClassicOrg(ctx, meta) {
		return resourceOIEOnlyFeatureError(resources.OktaIDaaSAuthenticator)
	}

	client := getOktaV6ClientFromMetadata(meta)

	authenticator, resp, err := client.AuthenticatorAPI.GetAuthenticator(ctx, d.Id()).Execute()
	if err != nil {
		return diag.Errorf("failed to get authenticator: %v", err)
	}
	defer resp.Body.Close()

	return establishAuthenticatorV6(authenticator, d, meta)
}

func resourceAuthenticatorUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if providerIsClassicOrg(ctx, meta) {
		return resourceOIEOnlyFeatureError(resources.OktaIDaaSAuthenticator)
	}

	client := getOktaV6ClientFromMetadata(meta)

	if err := validateAuthenticatorV6(d, meta); err != nil {
		return diag.FromErr(err)
	}

	authenticator, err := buildAuthenticatorV6(d, meta)
	if err != nil {
		return diag.Errorf("failed to build authenticator: %v", err)
	}

	_, resp, err := client.AuthenticatorAPI.ReplaceAuthenticator(ctx, d.Id()).Authenticator(*authenticator).Execute()
	if err != nil {
		return diag.Errorf("failed to update authenticator: %v", err)
	}
	defer resp.Body.Close()

	oldStatus, newStatus := d.GetChange("status")
	if oldStatus != newStatus {
		if newStatus == StatusActive {
			_, resp, err := client.AuthenticatorAPI.ActivateAuthenticator(ctx, d.Id()).Execute()
			if err != nil {
				return diag.Errorf("failed to change authenticator status: %v", err)
			}
			defer resp.Body.Close()
		} else {
			_, resp, err := client.AuthenticatorAPI.DeactivateAuthenticator(ctx, d.Id()).Execute()
			if err != nil {
				return diag.Errorf("failed to change authenticator status: %v", err)
			}
			defer resp.Body.Close()
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

	client := getOktaV6ClientFromMetadata(meta)

	_, resp, err := client.AuthenticatorAPI.DeactivateAuthenticator(ctx, d.Id()).Execute()
	if err != nil {
		logger(meta).Warn(fmt.Sprintf("Attempted to deactivate authenticator %q as soft delete and received error: %s", d.Get("key"), err))
	}
	if resp != nil {
		defer resp.Body.Close()
	}

	return nil
}

// Helper functions for v6 SDK

func findAuthenticatorV6(ctx context.Context, client *v6okta.APIClient, name, key string) (*v6okta.AuthenticatorBase, error) {
	authenticators, resp, err := client.AuthenticatorAPI.ListAuthenticators(ctx).Execute()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	for _, authUnion := range authenticators {
		// Convert union type to AuthenticatorBase via JSON marshaling
		authBytes, err := json.Marshal(authUnion)
		if err != nil {
			continue
		}
		var authenticator v6okta.AuthenticatorBase
		if err := json.Unmarshal(authBytes, &authenticator); err != nil {
			continue
		}

		// For custom_otp, both name and key must match
		if key == "custom_otp" {
			if authenticator.GetName() == name && authenticator.GetKey() == key {
				return &authenticator, nil
			}
		} else {
			// For other authenticators, match by name or key
			if authenticator.GetName() == name {
				return &authenticator, nil
			}
			if authenticator.GetKey() == key {
				return &authenticator, nil
			}
		}
	}

	if key != "" && key != "custom_otp" {
		return nil, fmt.Errorf("authenticator with key '%s' does not exist", key)
	}
	if key == "custom_otp" {
		return nil, fmt.Errorf("authenticator with name '%s' and key '%s' does not exist", name, key)
	}
	return nil, fmt.Errorf("authenticator with name '%s' does not exist", name)
}

func buildAuthenticatorV6(d *schema.ResourceData, meta interface{}) (*v6okta.AuthenticatorBase, error) {
	authenticator := v6okta.NewAuthenticatorBase()

	if d.Id() != "" {
		authenticator.SetId(d.Id())
	}
	authenticator.SetKey(d.Get("key").(string))
	authenticator.SetName(d.Get("name").(string))

	if typ, ok := d.GetOk("type"); ok {
		authenticator.SetType(typ.(string))
	}

	// Handle settings - stored in AdditionalProperties
	if d.Get("key").(string) != "custom_otp" {
		if s, ok := d.GetOk("settings"); ok {
			var settingsMap map[string]interface{}
			if err := json.Unmarshal([]byte(s.(string)), &settingsMap); err != nil {
				return nil, err
			}
			if authenticator.AdditionalProperties == nil {
				authenticator.AdditionalProperties = make(map[string]interface{})
			}
			authenticator.AdditionalProperties["settings"] = settingsMap
		}
	}

	// Handle provider configuration - stored in AdditionalProperties
	if p, ok := d.GetOk("provider_json"); ok {
		var providerMap map[string]interface{}
		if err := json.Unmarshal([]byte(p.(string)), &providerMap); err != nil {
			return nil, err
		}
		if authenticator.AdditionalProperties == nil {
			authenticator.AdditionalProperties = make(map[string]interface{})
		}
		authenticator.AdditionalProperties["provider"] = providerMap
	} else {
		// Build provider from individual fields based on authenticator type
		authType := d.Get("type").(string)
		var provider map[string]interface{}

		switch authType {
		case "security_key":
			provider = buildSecurityKeyProvider(d)
		default:
			// Handle provider type for non-security_key authenticators
			if providerType, ok := d.GetOk("provider_type"); ok {
				switch providerType.(string) {
				case "DUO":
					provider = buildDuoProvider(d)
				default:
					logger(meta).Warn("Unknown provider type - using default configuration",
						"provider_type", providerType.(string),
						"authenticator_key", d.Get("key").(string),
						"supported_types", []string{"DUO"})
				}
			}
		}

		if provider != nil {
			if authenticator.AdditionalProperties == nil {
				authenticator.AdditionalProperties = make(map[string]interface{})
			}
			authenticator.AdditionalProperties["provider"] = provider
		}
	}

	return authenticator, nil
}

func validateAuthenticatorV6(d *schema.ResourceData, meta interface{}) error {
	typ := d.Get("type").(string)
	key := d.Get("key").(string)

	switch typ {
	case "security_key":
		if key == "custom_otp" {
			return fmt.Errorf("custom_otp is not updatable")
		}

		h := d.Get("provider_hostname").(string)
		_, pok := d.GetOk("provider_auth_port")
		s := d.Get("provider_shared_secret").(string)
		templ := d.Get("provider_user_name_template").(string)
		if h == "" || s == "" || templ == "" || !pok {
			return fmt.Errorf("for authenticator type '%s' fields 'provider_hostname', "+
				"'provider_auth_port', 'provider_shared_secret' and 'provider_user_name_template' are required", typ)
		}
	default:
		logger(meta).Debug("Authenticator type using standard validation",
			"authenticator_type", typ)
	}

	// Validate provider-specific requirements
	if providerType, ok := d.GetOk("provider_type"); ok {
		switch providerType.(string) {
		case "DUO":
			h := d.Get("provider_host").(string)
			sk := d.Get("provider_secret_key").(string)
			ik := d.Get("provider_integration_key").(string)
			templ := d.Get("provider_user_name_template").(string)
			if h == "" || sk == "" || ik == "" || templ == "" {
				return fmt.Errorf("for authenticator type 'DUO' fields 'provider_host', " +
					"'provider_secret_key', 'provider_integration_key' and 'provider_user_name_template' are required")
			}
		default:
			logger(meta).Warn("Unknown provider type in validation - using default behavior",
				"provider_type", providerType.(string),
				"authenticator_key", d.Get("key").(string),
				"supported_types", []string{"DUO"})
		}
	}

	return nil
}

func establishAuthenticatorV6(authenticator *v6okta.AuthenticatorBase, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_ = d.Set("key", authenticator.GetKey())
	_ = d.Set("name", authenticator.GetName())
	_ = d.Set("status", authenticator.GetStatus())
	_ = d.Set("type", authenticator.GetType())

	// Extract settings from AdditionalProperties
	if authenticator.AdditionalProperties != nil {
		if settings, ok := authenticator.AdditionalProperties["settings"]; ok && settings != nil {
			b, _ := json.Marshal(settings)
			var dataMap map[string]interface{}
			_ = json.Unmarshal(b, &dataMap)
			b, _ = json.Marshal(dataMap)
			_ = d.Set("settings", string(b))
		}

		// Extract provider from AdditionalProperties
		if providerRaw, ok := authenticator.AdditionalProperties["provider"]; ok && providerRaw != nil {
			providerMap, ok := providerRaw.(map[string]interface{})
			if ok {
				if provType, ok := providerMap["type"].(string); ok {
					_ = d.Set("provider_type", provType)
				}

				if config, ok := providerMap["configuration"].(map[string]interface{}); ok {
					// Extract configuration based on authenticator type
					switch authenticator.GetType() {
					case "security_key":
						extractSecurityKeyConfig(config, d)
					default:
						logger(meta).Debug("Authenticator type using standard configuration extraction",
							"authenticator_type", authenticator.GetType())
					}

					// Extract configuration based on provider type
					if provType, ok := providerMap["type"].(string); ok {
						switch provType {
						case "DUO":
							extractDuoConfig(config, d)
						default:
							logger(meta).Warn("Unknown provider type in configuration extraction - using default behavior",
								"provider_type", provType,
								"authenticator_key", d.Get("key").(string),
								"supported_types", []string{"DUO"})
						}
					}

					// Extract common template configuration
					if template, ok := config["userNameTemplate"].(map[string]interface{}); ok {
						if t, ok := template["template"].(string); ok {
							_ = d.Set("provider_user_name_template", t)
						}
					}
				}
			}
		}
	}

	return nil
}

// Helper functions for provider building

func buildSecurityKeyProvider(d *schema.ResourceData) map[string]interface{} {
	provider := make(map[string]interface{})

	if provType, ok := d.GetOk("provider_type"); ok {
		provider["type"] = provType.(string)
	}

	config := make(map[string]interface{})
	if hostname, ok := d.GetOk("provider_hostname"); ok {
		config["hostName"] = hostname.(string)
	}
	if port, ok := d.GetOk("provider_auth_port"); ok {
		config["authPort"] = port.(int)
	}
	if secret, ok := d.GetOk("provider_shared_secret"); ok {
		config["sharedSecret"] = secret.(string)
	}
	if template, ok := d.GetOk("provider_user_name_template"); ok {
		config["userNameTemplate"] = map[string]interface{}{
			"template": template.(string),
		}
	}
	if instanceId, ok := d.GetOk("provider_instance_id"); ok {
		config["instanceId"] = instanceId.(string)
	}

	provider["configuration"] = config
	return provider
}

func buildDuoProvider(d *schema.ResourceData) map[string]interface{} {
	provider := make(map[string]interface{})
	provider["type"] = "DUO"

	config := make(map[string]interface{})
	if host, ok := d.GetOk("provider_host"); ok {
		config["host"] = host.(string)
	}
	if secretKey, ok := d.GetOk("provider_secret_key"); ok {
		config["secretKey"] = secretKey.(string)
	}
	if integrationKey, ok := d.GetOk("provider_integration_key"); ok {
		config["integrationKey"] = integrationKey.(string)
	}
	if template, ok := d.GetOk("provider_user_name_template"); ok {
		config["userNameTemplate"] = map[string]interface{}{
			"template": template.(string),
		}
	}

	provider["configuration"] = config
	return provider
}

// Helper functions for extracting provider configuration

func extractSecurityKeyConfig(config map[string]interface{}, d *schema.ResourceData) {
	if hostname, ok := config["hostName"].(string); ok {
		_ = d.Set("provider_hostname", hostname)
	}
	if authPort, ok := config["authPort"].(float64); ok {
		_ = d.Set("provider_auth_port", int(authPort))
	}
	if instanceId, ok := config["instanceId"].(string); ok {
		_ = d.Set("provider_instance_id", instanceId)
	}
}

func extractDuoConfig(config map[string]interface{}, d *schema.ResourceData) {
	if host, ok := config["host"].(string); ok {
		_ = d.Set("provider_host", host)
	}
	if secretKey, ok := config["secretKey"].(string); ok {
		_ = d.Set("provider_secret_key", secretKey)
	}
	if integrationKey, ok := config["integrationKey"].(string); ok {
		_ = d.Set("provider_integration_key", integrationKey)
	}
}
