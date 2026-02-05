package idaas

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourceAppSignOnPolicyRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppSignOnPolicyRuleCreate,
		ReadContext:   resourceAppSignOnPolicyRuleRead,
		UpdateContext: resourceAppSignOnPolicyRuleUpdate,
		DeleteContext: resourceAppSignOnPolicyRuleDelete,
		ValidateRawResourceConfigFuncs: []schema.ValidateRawResourceConfigFunc{
			func(ctx context.Context, req schema.ValidateResourceConfigFuncRequest, resp *schema.ValidateResourceConfigFuncResponse) {
				reauthenticateInFreqPresent := !req.RawConfig.GetAttr("re_authentication_frequency").IsNull()
				chains := req.RawConfig.GetAttr("chains")
				reauthenticateInInChainsPresent := !chains.IsNull() && strings.Contains(chains.GoString(), "reauthenticateIn")
				if reauthenticateInFreqPresent && reauthenticateInInChainsPresent {
					resp.Diagnostics = append(resp.Diagnostics, diag.Errorf("CANNOT set re_authentication_frequency AND reauthenticateIn in one or more entries in chains at the same time")...)
				}
			},
		},
		Importer: createPolicyRuleImporter(),
		Description: ` Manages a sign-on policy rules for the application.
~> **WARNING:** This feature is only available as a part of the Identity Engine. [Contact support](mailto:dev-inquiries@okta.com) for further information.
This resource allows you to create and configure a sign-on policy rule for the application.
A default or 'Catch-all Rule' sign-on policy rule can be imported and managed as a custom rule.
The only difference is that these fields are immutable and can not be managed: 'network_connection', 'network_excludes', 
'network_includes', 'platform_include', 'custom_expression', 'device_is_registered', 'device_is_managed', 'users_excluded',
'users_included', 'groups_excluded', 'groups_included', 'user_types_excluded' and 'user_types_included'.
~> **PRIORITY MANAGEMENT:** The Okta API automatically shifts rule priorities when conflicts occur. If you assign a rule to a priority already taken by another rule, the existing rule shifts to the next priority. This means directly swapping priorities between rules will cause drift. Use a two-step approach: first move rules to temporary high priorities (100+), apply, then move to final priorities. Always use 'depends_on' to chain rules sequentially based on priority order (ascending).`,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Policy Rule Name",
			},
			"policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the policy",
			},
			"system": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Often the `Catch-all Rule` this rule is the system (default) rule for its associated policy",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Status of the rule",
				Default:     StatusActive,
			},
			"priority": {
				Type:     schema.TypeInt,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					p, n := d.GetChange("priority")
					return p == n && new == "0"
				},
				Description: "Priority of the rule.",
			},
			"groups_included": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of group IDs to include",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"groups_excluded": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of group IDs to exclude",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"users_excluded": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Set of User IDs to exclude",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"users_included": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Set of User IDs to include",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"network_connection": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Network selection mode: ANYWHERE, ZONE, ON_NETWORK, or OFF_NETWORK.",
				Default:     "ANYWHERE",
			},
			"network_includes": {
				Type:          schema.TypeList,
				Optional:      true,
				Description:   "The zones to include",
				ConflictsWith: []string{"network_excludes"},
				Elem:          &schema.Schema{Type: schema.TypeString},
			},
			"network_excludes": {
				Type:          schema.TypeList,
				Optional:      true,
				Description:   "The zones to exclude",
				ConflictsWith: []string{"network_includes"},
				Elem:          &schema.Schema{Type: schema.TypeString},
			},
			"device_is_registered": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If the device is registered. A device is registered if the User enrolls with Okta Verify that is installed on the device.",
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					// Note: Keep this validator as it is enforcing payload
					// format the API is expecting and the side effects related
					// to that.
					if i == nil {
						return nil
					}
					v := i.(bool)
					if !v {
						return diag.Errorf("'device_is_registered' can either be set to 'true' or should not be present in the configuration")
					}
					return nil
				},
			},
			"device_is_managed": {
				Type:         schema.TypeBool,
				Optional:     true,
				RequiredWith: []string{"device_is_registered"},
				Description:  "If the device is managed. A device is managed if it's managed by a device management system. When managed is passed, registered must also be included and must be set to true.",
			},
			"device_assurances_included": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of device assurance IDs to include",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"platform_include": {
				Type:     schema.TypeSet,
				Elem:     platformIncludeResource,
				Optional: true,
			},
			"custom_expression": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This is an optional advanced setting. If the expression is formatted incorrectly or conflicts with conditions set above, the rule may not match any users.",
			},
			"user_types_excluded": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Set of User Type IDs to exclude",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"user_types_included": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Set of User Type IDs to include",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"access": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Allow or deny access based on the rule conditions: ALLOW or DENY",
				Default:     "ALLOW",
			},
			"factor_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The number of factors required to satisfy this assurance level",
				Default:     "2FA",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Verification Method type",
				Default:     "ASSURANCE",
			},
			"re_authentication_frequency": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The duration after which the end user must re-authenticate, regardless of user activity. Use the ISO 8601 Period format for recurring time intervals. PT0S - Every sign-in attempt, PT43800H - Once per session",
				Default:     "PT2H",
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					// check if reauthenticateIn is present for any of the entries in chains.
					if vc, ok := d.GetOk("chains"); ok {
						valueList := vc.([]interface{})
						for _, item := range valueList {
							var chain sdk.AccessPolicyChains
							_ = json.Unmarshal([]byte(item.(string)), &chain)
							if chain.ReauthenticateIn != "" {
								return true // ignore re_authentication_frequency as it'll fail if it has been set in config and if it hasn't been (set) then it's default value
							}
						}
					}
					return oldValue == newValue
				},
			},
			"inactivity_period": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The inactivity duration after which the end user must re-authenticate. Use the ISO 8601 Period format for recurring time intervals.",
			},
			"constraints": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: stringIsJSON,
					StateFunc:        utils.NormalizeDataJSON,
					DiffSuppressFunc: utils.NoChangeInObjectFromUnmarshaledJSON,
				},
				Optional:    true,
				Description: "An array that contains nested Authenticator Constraint objects that are organized by the Authenticator class",
			},
			"risk_score": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The risk score specifies a particular level of risk to match on: ANY, LOW, MEDIUM, HIGH",
			},
			"chains": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: stringIsJSON,
					StateFunc:        utils.NormalizeDataJSON,
					DiffSuppressFunc: utils.NoChangeInObjectFromUnmarshaledJSON,
				},
				Optional:    true,
				Description: "Use with verification method = `AUTH_METHOD_CHAIN` only",
			},
		},
	}
}

//url := fmt.Sprintf("/api/v1/policies/%v/rules", policyID)
//req, err := m.RequestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodPost, url, body)
//if err != nil {
//return nil, nil, err
//}
//var appSignOnPolicyRule *AccessPolicyRule
//resp, err := m.RequestExecutor.Do(ctx, req, &appSignOnPolicyRule)
//if err != nil {
//return nil, resp, err
//}
//return appSignOnPolicyRule, resp, nil

func resourceAppSignOnPolicyRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// See - https://github.com/okta/terraform-provider-okta/pull/1248
	oktaMutexKV.Lock(resources.OktaIDaaSAppSignOnPolicyRule)
	defer oktaMutexKV.Unlock(resources.OktaIDaaSAppSignOnPolicyRule)

	if providerIsClassicOrg(ctx, meta) {
		return resourceOIEOnlyFeatureError(resources.OktaIDaaSAppSignOnPolicyRule)
	}

	rule, _, err := getAPISupplementFromMetadata(meta).CreateAppSignOnPolicyRule(ctx, d.Get("policy_id").(string), buildAppSignOnPolicyRule(d))
	if err != nil {
		return diag.Errorf("failed to create app sign on policy rule: %v", err)
	}
	d.SetId(rule.Id)
	if status, ok := d.GetOk("status"); ok {
		if status.(string) == StatusInactive {
			_, err = getAPISupplementFromMetadata(meta).DeactivateAppSignOnPolicyRule(ctx, d.Get("policy_id").(string), d.Id())
			if err != nil {
				return diag.Errorf("failed to deactivate app sign on policy rule: %v", err)
			}
		}
	}
	return resourceAppSignOnPolicyRuleRead(ctx, d, meta)
}

func resourceAppSignOnPolicyRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if providerIsClassicOrg(ctx, meta) {
		return resourceOIEOnlyFeatureError(resources.OktaIDaaSAppSignOnPolicyRule)
	}

	rule, resp, err := getAPISupplementFromMetadata(meta).GetAppSignOnPolicyRule(ctx, d.Get("policy_id").(string), d.Id())
	if err := utils.SuppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get app sign on policy rule: %v", err)
	}
	if rule == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("system", utils.BoolFromBoolPtr(rule.System))
	_ = d.Set("name", rule.Name)
	if rule.PriorityPtr != nil {
		_ = d.Set("priority", rule.PriorityPtr)
	}
	_ = d.Set("status", rule.Status)
	if rule.Actions.AppSignOn != nil {
		_ = d.Set("access", rule.Actions.AppSignOn.Access)
		if rule.Actions.AppSignOn.VerificationMethod != nil {
			_ = d.Set("type", rule.Actions.AppSignOn.VerificationMethod.Type)
			if rule.Actions.AppSignOn.VerificationMethod.FactorMode != "" {
				_ = d.Set("factor_mode", rule.Actions.AppSignOn.VerificationMethod.FactorMode)
			} else {
				_ = d.Set("factor_mode", "2FA")
			}
			_ = d.Set("re_authentication_frequency", rule.Actions.AppSignOn.VerificationMethod.ReauthenticateIn)
			if rule.Actions.AppSignOn.VerificationMethod.InactivityPeriod != "" {
				_ = d.Set("inactivity_period", rule.Actions.AppSignOn.VerificationMethod.InactivityPeriod)
			}
			constraintArr := make([]interface{}, len(rule.Actions.AppSignOn.VerificationMethod.Constraints))
			for i := range rule.Actions.AppSignOn.VerificationMethod.Constraints {
				b, _ := json.Marshal(rule.Actions.AppSignOn.VerificationMethod.Constraints[i])
				constraintArr[i] = string(b)
			}
			_ = d.Set("constraints", constraintArr)
			if len(rule.Actions.AppSignOn.VerificationMethod.Chains) != 0 {
				chainArr := make([]interface{}, len(rule.Actions.AppSignOn.VerificationMethod.Chains))
				for i := range rule.Actions.AppSignOn.VerificationMethod.Chains {
					b, _ := json.Marshal(rule.Actions.AppSignOn.VerificationMethod.Chains[i])
					chainArr[i] = string(b)
				}
				_ = d.Set("chains", chainArr)
			}
		}
	}
	if rule.Conditions != nil {
		if rule.Conditions.ElCondition != nil {
			_ = d.Set("custom_expression", rule.Conditions.ElCondition.Condition)
		}
		m := map[string]interface{}{
			"platform_include": flattenAccessPolicyPlatformInclude(rule.Conditions.Platform),
		}
		_ = d.Set("network_connection", rule.Conditions.Network.Connection)
		if len(rule.Conditions.Network.Include) > 0 {
			m["network_includes"] = utils.ConvertStringSliceToInterfaceSlice(rule.Conditions.Network.Include)
		}
		if len(rule.Conditions.Network.Exclude) > 0 {
			m["network_excludes"] = utils.ConvertStringSliceToInterfaceSlice(rule.Conditions.Network.Exclude)
		}
		if rule.Conditions.Device != nil {
			_ = d.Set("device_is_managed", rule.Conditions.Device.Managed)
			_ = d.Set("device_is_registered", rule.Conditions.Device.Registered)
			if rule.Conditions.Device.Assurance != nil {
				m["device_assurances_included"] = utils.ConvertStringSliceToSetNullable(rule.Conditions.Device.Assurance.Include)
			}
		}
		if rule.Conditions.People != nil {
			if rule.Conditions.People.Users != nil {
				m["users_excluded"] = utils.ConvertStringSliceToSetNullable(rule.Conditions.People.Users.Exclude)
				m["users_included"] = utils.ConvertStringSliceToSetNullable(rule.Conditions.People.Users.Include)
			}
			if rule.Conditions.People.Groups != nil {
				m["groups_excluded"] = utils.ConvertStringSliceToSetNullable(rule.Conditions.People.Groups.Exclude)
				m["groups_included"] = utils.ConvertStringSliceToSetNullable(rule.Conditions.People.Groups.Include)
			}
		}
		if rule.Conditions.UserType != nil {
			m["user_types_excluded"] = utils.ConvertStringSliceToSetNullable(rule.Conditions.UserType.Exclude)
			m["user_types_included"] = utils.ConvertStringSliceToSetNullable(rule.Conditions.UserType.Include)
		}
		if rule.Conditions.RiskScore != nil {
			_ = d.Set("risk_score", rule.Conditions.RiskScore.Level)
		}
		_ = utils.SetNonPrimitives(d, m)
	}
	return nil
}

func resourceAppSignOnPolicyRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	oktaMutexKV.Lock(resources.OktaIDaaSAppSignOnPolicyRule)
	defer oktaMutexKV.Unlock(resources.OktaIDaaSAppSignOnPolicyRule)

	if providerIsClassicOrg(ctx, meta) {
		return resourceOIEOnlyFeatureError(resources.OktaIDaaSAppSignOnPolicyRule)
	}

	rule := buildAppSignOnPolicyRule(d)
	if utils.BoolFromBoolPtr(rule.System) {
		// Conditions can't be set on the default/system rule
		rule.Conditions = nil
	}
	_, _, err := getAPISupplementFromMetadata(meta).UpdateAppSignOnPolicyRule(ctx, d.Get("policy_id").(string), d.Id(), rule)
	if err != nil {
		return diag.Errorf("failed to update app sign on policy rule: %v", err)
	}
	oldStatus, newStatus := d.GetChange("status")
	if oldStatus != newStatus {
		if newStatus == StatusActive {
			_, err = getAPISupplementFromMetadata(meta).ActivateAppSignOnPolicyRule(ctx, d.Get("policy_id").(string), d.Id())
		} else {
			_, err = getAPISupplementFromMetadata(meta).DeactivateAppSignOnPolicyRule(ctx, d.Get("policy_id").(string), d.Id())
		}
		if err != nil {
			return diag.Errorf("failed to change app sign on policy rule status: %v", err)
		}
	}
	return resourceAppSignOnPolicyRuleRead(ctx, d, meta)
}

func resourceAppSignOnPolicyRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if providerIsClassicOrg(ctx, meta) {
		return resourceOIEOnlyFeatureError(resources.OktaIDaaSAppSignOnPolicyRule)
	}

	oktaMutexKV.Lock(resources.OktaIDaaSAppSignOnPolicyRule)
	defer oktaMutexKV.Unlock(resources.OktaIDaaSAppSignOnPolicyRule)

	if d.Get("name") == "Catch-all Rule" {
		// You cannot delete a default rule in a policy
		return nil
	}
	_, err := getAPISupplementFromMetadata(meta).DeleteAppSignOnPolicyRule(ctx, d.Get("policy_id").(string), d.Id())
	if err != nil {
		return diag.Errorf("failed to delete app sign-on policy rule: %v", err)
	}
	return nil
}

func buildAppSignOnPolicyRule(d *schema.ResourceData) sdk.AccessPolicyRule {
	rule := sdk.AccessPolicyRule{
		Actions: &sdk.AccessPolicyRuleActions{
			AppSignOn: &sdk.AccessPolicyRuleApplicationSignOn{
				Access: d.Get("access").(string),
				VerificationMethod: &sdk.VerificationMethod{
					FactorMode:       d.Get("factor_mode").(string),
					ReauthenticateIn: d.Get("re_authentication_frequency").(string),
					InactivityPeriod: d.Get("inactivity_period").(string),
					Type:             d.Get("type").(string),
				},
			},
		},
		Name:        d.Get("name").(string),
		PriorityPtr: utils.Int64Ptr(d.Get("priority").(int)),
		Type:        "ACCESS_POLICY", // TODO New types of access policy rules like MFA_ENROLL etc. will be supported iff there is an ask.
	}

	// NOTE: Only the API read will be able to set the "system" boolean so it is
	// ok to inspect the resource data for its presence to set the bool pointer.
	// When buildAppSignOnPolicyRule is called from the create context the bool
	// pointer is effectively inert (nil) and we don't need additional logic
	// about if this is being called for create/read/update.
	if v, ok := d.GetOk("system"); ok {
		rule.System = utils.BoolPtr(v.(bool))
	}

	var constraints []*sdk.AccessPolicyConstraints
	v, ok := d.GetOk("constraints")
	if ok {
		valueList := v.([]interface{})
		for _, item := range valueList {
			var constraint sdk.AccessPolicyConstraints
			_ = json.Unmarshal([]byte(item.(string)), &constraint)
			constraints = append(constraints, &constraint)
		}
	}
	rule.Actions.AppSignOn.VerificationMethod.Constraints = constraints
	var chains []*sdk.AccessPolicyChains
	vc, ok := d.GetOk("chains")
	if ok { // condition should be true iff type is verification method = AUTH_METHOD_CHAIN
		valueList := vc.([]interface{})
		for _, item := range valueList {
			var chain sdk.AccessPolicyChains
			_ = json.Unmarshal([]byte(item.(string)), &chain)
			chains = append(chains, &chain)
			if chain.ReauthenticateIn != "" {
				// if ReauthenticateIn has been set in any chain, unset it in VerificationMethod as the combination isn't supported .
				rule.Actions.AppSignOn.VerificationMethod.ReauthenticateIn = ""
			}
		}
	}
	rule.Actions.AppSignOn.VerificationMethod.Chains = chains
	rule.Conditions = &sdk.AccessPolicyRuleConditions{
		Network: buildPolicyNetworkCondition(d),
		Platform: &sdk.PlatformPolicyRuleCondition{
			Include: buildAccessPolicyPlatformInclude(d),
		},
		ElCondition: &sdk.AccessPolicyRuleCustomCondition{
			Condition: d.Get("custom_expression").(string),
		},
	}
	riskScore, ok := d.GetOk("risk_score")
	if ok {
		rule.Conditions.RiskScore = &sdk.RiskScorePolicyRuleCondition{
			Level: riskScore.(string),
		}
	}
	isRegistered, ok := d.GetOk("device_is_registered")
	if ok && isRegistered.(bool) {
		rule.Conditions.Device = &sdk.DeviceAccessPolicyRuleCondition{
			Managed:    utils.BoolPtr(d.Get("device_is_managed").(bool)),
			Registered: utils.BoolPtr(isRegistered.(bool)),
		}
	}
	deviceAssurancesIncluded, deviceAssurancesIncludedOk := d.GetOk("device_assurances_included")
	if deviceAssurancesIncludedOk {
		if rule.Conditions.Device != nil {
			rule.Conditions.Device.Assurance = &sdk.DeviceAssurancePolicyRuleCondition{
				Include: utils.ConvertInterfaceToStringSetNullable(deviceAssurancesIncluded),
			}
		} else {
			rule.Conditions.Device = &sdk.DeviceAccessPolicyRuleCondition{
				Assurance: &sdk.DeviceAssurancePolicyRuleCondition{
					Include: utils.ConvertInterfaceToStringSetNullable(deviceAssurancesIncluded),
				},
			}
		}
	}

	usersExcluded, usersExcludedOk := d.GetOk("users_excluded")
	usersIncluded, usersIncludedOk := d.GetOk("users_included")
	if usersExcludedOk || usersIncludedOk {
		rule.Conditions.People = &sdk.PolicyPeopleCondition{
			Users: &sdk.UserCondition{
				Exclude: utils.ConvertInterfaceToStringSetNullable(usersExcluded),
				Include: utils.ConvertInterfaceToStringSetNullable(usersIncluded),
			},
		}
	}
	groupsExcluded, groupsExcludedOk := d.GetOk("groups_excluded")
	groupsIncluded, groupsIncludedOk := d.GetOk("groups_included")
	if groupsExcludedOk || groupsIncludedOk {
		if rule.Conditions.People == nil {
			rule.Conditions.People = &sdk.PolicyPeopleCondition{}
		}
		rule.Conditions.People.Groups = &sdk.GroupCondition{
			Exclude: utils.ConvertInterfaceToStringSetNullable(groupsExcluded),
			Include: utils.ConvertInterfaceToStringSetNullable(groupsIncluded),
		}
	}
	userTypesExcluded, userTypesExcludedOk := d.GetOk("user_types_excluded")
	userTypesIncluded, userTypesIncludedOk := d.GetOk("user_types_included")
	if userTypesExcludedOk || userTypesIncludedOk {
		rule.Conditions.UserType = &sdk.UserTypeCondition{
			Exclude: utils.ConvertInterfaceToStringSetNullable(userTypesExcluded),
			Include: utils.ConvertInterfaceToStringSetNullable(userTypesIncluded),
		}
	}
	return rule
}

func buildAccessPolicyPlatformInclude(d *schema.ResourceData) []*sdk.PlatformConditionEvaluatorPlatform {
	var includeList []*sdk.PlatformConditionEvaluatorPlatform
	v, ok := d.GetOk("platform_include")
	if !ok {
		return includeList
	}
	valueList := v.(*schema.Set).List()
	for _, item := range valueList {
		if value, ok := item.(map[string]interface{}); ok {
			var expr *string
			if typ := utils.GetMapString(value, "os_type"); typ == "OTHER" {
				if v, ok := value["os_expression"]; ok {
					if v != nil {
						res := v.(string)
						expr = &res
					}
				}
			}
			includeList = append(includeList, &sdk.PlatformConditionEvaluatorPlatform{
				Os: &sdk.PlatformConditionEvaluatorPlatformOperatingSystem{
					Expression: expr,
					Type:       utils.GetMapString(value, "os_type"),
				},
				Type: utils.GetMapString(value, "type"),
			})
		}
	}
	return includeList
}

func flattenAccessPolicyPlatformInclude(platform *sdk.PlatformPolicyRuleCondition) *schema.Set {
	var flattened []interface{}
	if platform != nil && platform.Include != nil {
		for _, v := range platform.Include {
			var expr *string
			if v.Os.Expression != nil {
				expr = v.Os.Expression
			}
			m := map[string]interface{}{
				"os_type": v.Os.Type,
				"type":    v.Type,
			}
			if expr != nil {
				m["os_expression"] = *expr
			}
			flattened = append(flattened, m)
		}
	}
	return schema.NewSet(schema.HashResource(platformIncludeResource), flattened)
}
