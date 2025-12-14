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

-> **Methods:** Some authenticators support multiple methods (e.g., Phone supports SMS and Voice, Okta Verify supports Push, TOTP, and Signed Nonce). Use the 'method' block to control individual methods. If no method blocks are specified, all methods will remain in their current state.

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
			"method": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Configuration block for authenticator methods. Only applicable for authenticators that support multiple methods (e.g., `phone_number`, `okta_verify`)",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Type of the authenticator method. For `phone_number`: `sms`, `voice`. For `okta_verify`: `push`, `totp`, `signed_nonce`. For `custom_otp`: `otp`",
						},
						"status": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     StatusActive,
							Description: "Status of the method: `ACTIVE` or `INACTIVE`. Default: `ACTIVE`",
						},
						"settings": {
							Type:             schema.TypeString,
							Optional:         true,
							Description:      "Method-specific settings in JSON format. Required settings vary by method type. See Okta API documentation for details",
							ValidateDiagFunc: stringIsJSON,
							StateFunc:        utils.NormalizeDataJSON,
							DiffSuppressFunc: utils.NoChangeInObjectFromUnmarshaledJSON,
						},
					},
				},
				Set: func(i interface{}) int {
					// Hash function for TypeSet - use method type as unique identifier
					m := i.(map[string]interface{})
					return schema.HashString(m["type"].(string))
				},
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
			_, resp, err := client.AuthenticatorAPI.DeactivateAuthenticator(ctx, d.Id()).Execute()
			if err != nil {
				return diag.Errorf("failed to change authenticator status: %v", err)
			}
			defer resp.Body.Close()
		} else {
			_, resp, err := client.AuthenticatorAPI.ActivateAuthenticator(ctx, d.Id()).Execute()
			if err != nil {
				return diag.Errorf("failed to change authenticator status: %v", err)
			}
			defer resp.Body.Close()
		}
	}

	// Manage authenticator methods if specified and supported
	if supportsAuthenticatorMethods(d.Get("key").(string)) {
		if _, ok := d.GetOk("method"); ok {
			desiredMethods := getMethodsFromSchema(d)
			if len(desiredMethods) > 0 {
				if err := syncAuthenticatorMethods(ctx, client, d.Id(), desiredMethods, meta); err != nil {
					return diag.Errorf("failed to sync authenticator methods: %v", err)
				}
			}
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

	// Read authenticator methods if supported
	if supportsAuthenticatorMethods(authenticator.GetKey()) {
		methods, err := listAuthenticatorMethodsV6(ctx, client, d.Id(), meta)
		if err != nil {
			logger(meta).Warn(fmt.Sprintf("Failed to list authenticator methods: %v", err))
		} else if len(methods) > 0 {
			if err := d.Set("method", flattenAuthenticatorMethods(methods)); err != nil {
				return diag.Errorf("failed to set method: %v", err)
			}
		}
	}

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

	// Handle method changes if supported
	if d.HasChange("method") && supportsAuthenticatorMethods(d.Get("key").(string)) {
		desiredMethods := getMethodsFromSchema(d)
		if len(desiredMethods) > 0 {
			if err := syncAuthenticatorMethods(ctx, client, d.Id(), desiredMethods, meta); err != nil {
				return diag.Errorf("failed to sync authenticator methods during update: %v", err)
			}
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

	// Validate method blocks if present
	if err := validateAuthenticatorMethods(d, key, meta); err != nil {
		return err
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

// authenticatorMethod represents a method configuration
type authenticatorMethod struct {
	Type     string
	Status   string
	Settings map[string]interface{}
}

// supportsAuthenticatorMethods checks if the authenticator supports method-level configuration
func supportsAuthenticatorMethods(key string) bool {
	supportedKeys := map[string]bool{
		"phone_number": true,
		"okta_verify":  true,
		"custom_otp":   true,
	}
	return supportedKeys[key]
}

// getMethodsFromSchema extracts method blocks from Terraform schema
func getMethodsFromSchema(d *schema.ResourceData) []authenticatorMethod {
	var methods []authenticatorMethod

	if methodSet, ok := d.GetOk("method"); ok {
		for _, m := range methodSet.(*schema.Set).List() {
			methodMap := m.(map[string]interface{})
			method := authenticatorMethod{
				Type:   methodMap["type"].(string),
				Status: methodMap["status"].(string),
			}

			// Parse settings if present
			if settingsStr, ok := methodMap["settings"].(string); ok && settingsStr != "" {
				var settings map[string]interface{}
				if err := json.Unmarshal([]byte(settingsStr), &settings); err == nil {
					method.Settings = settings
				}
			}

			methods = append(methods, method)
		}
	}

	return methods
}

// listAuthenticatorMethodsV6 fetches all methods for an authenticator
func listAuthenticatorMethodsV6(ctx context.Context, client *v6okta.APIClient, authenticatorId string, meta interface{}) ([]authenticatorMethod, error) {
	methodsResp, resp, err := client.AuthenticatorAPI.ListAuthenticatorMethods(ctx, authenticatorId).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to list authenticator methods: %v", err)
	}
	defer resp.Body.Close()

	var methods []authenticatorMethod
	for _, methodUnion := range methodsResp {
		// Convert union type to map via JSON marshaling
		methodBytes, err := json.Marshal(methodUnion)
		if err != nil {
			logger(meta).Warn(fmt.Sprintf("Failed to marshal method: %v", err))
			continue
		}

		var methodMap map[string]interface{}
		if err := json.Unmarshal(methodBytes, &methodMap); err != nil {
			logger(meta).Warn(fmt.Sprintf("Failed to unmarshal method: %v", err))
			continue
		}

		method := authenticatorMethod{}
		if t, ok := methodMap["type"].(string); ok {
			method.Type = t
		}
		if s, ok := methodMap["status"].(string); ok {
			method.Status = s
		}
		if settings, ok := methodMap["settings"].(map[string]interface{}); ok {
			method.Settings = settings
		}

		methods = append(methods, method)
	}

	return methods, nil
}

// syncAuthenticatorMethods manages method activation/deactivation and settings
func syncAuthenticatorMethods(ctx context.Context, client *v6okta.APIClient, authenticatorId string, desiredMethods []authenticatorMethod, meta interface{}) error {
	// Get current methods from API
	currentMethods, err := listAuthenticatorMethodsV6(ctx, client, authenticatorId, meta)
	if err != nil {
		return err
	}

	// Create maps for easy lookup
	currentMethodMap := make(map[string]authenticatorMethod)
	for _, m := range currentMethods {
		currentMethodMap[m.Type] = m
	}

	desiredMethodMap := make(map[string]authenticatorMethod)
	for _, m := range desiredMethods {
		desiredMethodMap[m.Type] = m
	}

	// Process each desired method
	for _, desired := range desiredMethods {
		current, exists := currentMethodMap[desired.Type]

		// Determine if we need to update the method
		needsUpdate := false
		if !exists {
			needsUpdate = true
		} else {
			// Check if status changed
			if current.Status != desired.Status {
				needsUpdate = true
			}
			// Check if settings changed (if settings are provided)
			if desired.Settings != nil {
				needsUpdate = true
			}
		}

		if needsUpdate {
			// Only update status if method exists and status changed
			if exists && current.Status != desired.Status {
				if desired.Status == StatusActive {
					if err := activateAuthenticatorMethodV6(ctx, client, authenticatorId, desired.Type, meta); err != nil {
						return err
					}
				} else if desired.Status == StatusInactive {
					if err := deactivateAuthenticatorMethodV6(ctx, client, authenticatorId, desired.Type, meta); err != nil {
						return err
					}
				}
			}

			// Update method settings if provided (only if method exists or will be activated)
			if desired.Settings != nil && (exists || desired.Status == StatusActive) {
				if err := updateAuthenticatorMethodV6(ctx, client, authenticatorId, desired.Type, desired.Settings, meta); err != nil {
					// Log warning but don't fail - some methods may not support all settings
					logger(meta).Warn(fmt.Sprintf("Failed to update settings for method %s: %v", desired.Type, err))
				}
			}
		}
	}

	return nil
}

// activateAuthenticatorMethodV6 activates a specific method
func activateAuthenticatorMethodV6(ctx context.Context, client *v6okta.APIClient, authenticatorId, methodType string, meta interface{}) error {
	_, resp, err := client.AuthenticatorAPI.ActivateAuthenticatorMethod(ctx, authenticatorId, methodType).Execute()
	if err != nil {
		return fmt.Errorf("failed to activate method %s: %v", methodType, err)
	}
	defer resp.Body.Close()

	logger(meta).Info(fmt.Sprintf("Activated method %s for authenticator %s", methodType, authenticatorId))
	return nil
}

// deactivateAuthenticatorMethodV6 deactivates a specific method
func deactivateAuthenticatorMethodV6(ctx context.Context, client *v6okta.APIClient, authenticatorId, methodType string, meta interface{}) error {
	_, resp, err := client.AuthenticatorAPI.DeactivateAuthenticatorMethod(ctx, authenticatorId, methodType).Execute()
	if err != nil {
		return fmt.Errorf("failed to deactivate method %s: %v", methodType, err)
	}
	defer resp.Body.Close()

	logger(meta).Info(fmt.Sprintf("Deactivated method %s for authenticator %s", methodType, authenticatorId))
	return nil
}

// updateAuthenticatorMethodV6 updates method settings
func updateAuthenticatorMethodV6(ctx context.Context, client *v6okta.APIClient, authenticatorId, methodType string, settings map[string]interface{}, meta interface{}) error {
	// Build the method payload
	methodPayload := map[string]interface{}{
		"type":     methodType,
		"settings": settings,
	}

	methodReq := client.AuthenticatorAPI.ReplaceAuthenticatorMethod(ctx, authenticatorId, methodType)
	_, resp, err := methodReq.Execute()
	if err != nil {
		methodBytes, _ := json.Marshal(methodPayload)
		return fmt.Errorf("failed to update method %s settings: %v, payload: %s", methodType, err, string(methodBytes))
	}
	defer resp.Body.Close()

	logger(meta).Info(fmt.Sprintf("Updated settings for method %s on authenticator %s", methodType, authenticatorId))
	return nil
}

// flattenAuthenticatorMethods converts API methods to Terraform state format
func flattenAuthenticatorMethods(methods []authenticatorMethod) []interface{} {
	result := make([]interface{}, 0, len(methods))

	for _, method := range methods {
		m := make(map[string]interface{})
		m["type"] = method.Type
		m["status"] = method.Status

		if method.Settings != nil && len(method.Settings) > 0 {
			settingsBytes, err := json.Marshal(method.Settings)
			if err == nil {
				m["settings"] = string(settingsBytes)
			}
		}

		result = append(result, m)
	}

	return result
}

// validateAuthenticatorMethods validates method blocks for an authenticator
func validateAuthenticatorMethods(d *schema.ResourceData, key string, meta interface{}) error {
	methodSet, ok := d.GetOk("method")
	if !ok {
		// No methods defined, validation passes
		return nil
	}

	// Check if this authenticator supports methods
	if !supportsAuthenticatorMethods(key) {
		return fmt.Errorf("authenticator with key '%s' does not support method blocks. Only 'phone_number', 'okta_verify', and 'custom_otp' authenticators support methods", key)
	}

	methods := methodSet.(*schema.Set).List()
	if len(methods) == 0 {
		// Empty method set, validation passes
		return nil
	}

	// Define valid method types for each authenticator
	validMethodTypes := map[string][]string{
		"phone_number": {"sms", "voice"},
		"okta_verify":  {"push", "totp", "signed_nonce"},
		"custom_otp":   {"otp"},
	}

	allowedTypes, exists := validMethodTypes[key]
	if !exists {
		logger(meta).Warn(fmt.Sprintf("No method type validation rules defined for authenticator key: %s", key))
		return nil
	}

	// Validate each method
	seenTypes := make(map[string]bool)
	for _, m := range methods {
		methodMap := m.(map[string]interface{})
		methodType := methodMap["type"].(string)
		methodStatus := methodMap["status"].(string)

		// Check for duplicate method types (should be prevented by Set, but double-check)
		if seenTypes[methodType] {
			return fmt.Errorf("duplicate method type '%s' found. Each method type can only be specified once", methodType)
		}
		seenTypes[methodType] = true

		// Validate method type is allowed for this authenticator
		validType := false
		for _, allowed := range allowedTypes {
			if methodType == allowed {
				validType = true
				break
			}
		}
		if !validType {
			return fmt.Errorf("invalid method type '%s' for authenticator '%s'. Valid types are: %v", methodType, key, allowedTypes)
		}

		// Validate status values
		if methodStatus != StatusActive && methodStatus != StatusInactive {
			return fmt.Errorf("invalid status '%s' for method '%s'. Status must be either 'ACTIVE' or 'INACTIVE'", methodStatus, methodType)
		}

		// Validate settings if present
		if settingsStr, ok := methodMap["settings"].(string); ok && settingsStr != "" {
			var settings map[string]interface{}
			if err := json.Unmarshal([]byte(settingsStr), &settings); err != nil {
				return fmt.Errorf("invalid JSON in settings for method '%s': %v", methodType, err)
			}

			// Perform method-specific settings validation
			if err := validateMethodSettings(key, methodType, settings, meta); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateMethodSettings performs method-specific settings validation
func validateMethodSettings(authenticatorKey, methodType string, settings map[string]interface{}, meta interface{}) error {
	// Define required/optional settings for each method type
	// This is a basic validation - the API will perform more detailed validation

	switch authenticatorKey {
	case "okta_verify":
		switch methodType {
		case "push", "signed_nonce":
			// Validate algorithms if present
			if algorithms, ok := settings["algorithms"]; ok {
				algList, ok := algorithms.([]interface{})
				if !ok {
					return fmt.Errorf("'algorithms' in method '%s' settings must be an array", methodType)
				}
				validAlgorithms := map[string]bool{"ES256": true, "RS256": true, "ES384": true, "ES512": true, "RS384": true, "RS512": true, "EdDSA": true}
				for _, alg := range algList {
					algStr, ok := alg.(string)
					if !ok || !validAlgorithms[algStr] {
						return fmt.Errorf("invalid algorithm '%v' in method '%s'. Valid algorithms: ES256, RS256, ES384, ES512, RS384, RS512, EdDSA", alg, methodType)
					}
				}
			}

			// Validate keyProtection if present
			if keyProtection, ok := settings["keyProtection"]; ok {
				kpStr, ok := keyProtection.(string)
				if !ok {
					return fmt.Errorf("'keyProtection' in method '%s' settings must be a string", methodType)
				}
				validKeyProtection := map[string]bool{"ANY": true, "SOFTWARE": true, "HARDWARE": true}
				if !validKeyProtection[kpStr] {
					return fmt.Errorf("invalid keyProtection '%s' in method '%s'. Valid values: ANY, SOFTWARE, HARDWARE", kpStr, methodType)
				}
			}

		case "totp":
			// TOTP settings validation
			if timeInterval, ok := settings["timeIntervalInSeconds"]; ok {
				if _, ok := timeInterval.(float64); !ok {
					return fmt.Errorf("'timeIntervalInSeconds' in method '%s' settings must be a number", methodType)
				}
			}
			if encoding, ok := settings["encoding"]; ok {
				encodingStr, ok := encoding.(string)
				if !ok {
					return fmt.Errorf("'encoding' in method '%s' settings must be a string", methodType)
				}
				validEncodings := map[string]bool{"base32": true, "base64": true}
				if !validEncodings[encodingStr] {
					return fmt.Errorf("invalid encoding '%s' in method '%s'. Valid values: base32, base64", encodingStr, methodType)
				}
			}
		}

	case "custom_otp":
		if methodType == "otp" {
			// Validate protocol
			if protocol, ok := settings["protocol"]; ok {
				protocolStr, ok := protocol.(string)
				if !ok {
					return fmt.Errorf("'protocol' in method '%s' settings must be a string", methodType)
				}
				validProtocols := map[string]bool{"TOTP": true, "HOTP": true}
				if !validProtocols[protocolStr] {
					return fmt.Errorf("invalid protocol '%s' in method '%s'. Valid values: TOTP, HOTP", protocolStr, methodType)
				}
			}
		}

	case "phone_number":
		// Phone methods (sms, voice) typically don't have complex settings to validate
		logger(meta).Debug(fmt.Sprintf("No specific settings validation for phone method '%s'", methodType))
	}

	return nil
}
