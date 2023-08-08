package okta

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourceAppSignOnPolicyRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppSignOnPolicyRuleCreate,
		ReadContext:   resourceAppSignOnPolicyRuleRead,
		UpdateContext: resourceAppSignOnPolicyRuleUpdate,
		DeleteContext: resourceAppSignOnPolicyRuleDelete,
		Importer:      createPolicyRuleImporter(),
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
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Status of the rule",
				Default:     statusActive,
			},
			"priority": {
				Type:     schema.TypeInt,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					p, n := d.GetChange("priority")
					return p == n && new == "0"
				},
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
			},
			"inactivity_period": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The inactivity duration after which the end user must re-authenticate. Use the ISO 8601 Period format for recurring time intervals.",
				Default:     "PT1H",
			},
			"constraints": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: stringIsJSON,
					StateFunc:        normalizeDataJSON,
				},
				Optional:    true,
				Description: "An array that contains nested Authenticator Constraint objects that are organized by the Authenticator class",
			},
		},
	}
}

func resourceAppSignOnPolicyRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if isClassicOrg(ctx, m) {
		return resourceOIEOnlyFeatureError(appSignOnPolicyRule)
	}

	rule, _, err := getAPISupplementFromMetadata(m).CreateAppSignOnPolicyRule(ctx, d.Get("policy_id").(string), buildAppSignOnPolicyRule(d))
	if err != nil {
		return diag.Errorf("failed to create app sign on policy rule: %v", err)
	}
	d.SetId(rule.Id)
	if status, ok := d.GetOk("status"); ok {
		if status.(string) == statusInactive {
			_, err = getAPISupplementFromMetadata(m).DeactivateAppSignOnPolicyRule(ctx, d.Get("policy_id").(string), d.Id())
			if err != nil {
				return diag.Errorf("failed to deactivate app sign on policy rule: %v", err)
			}
		}
	}
	return resourceAppSignOnPolicyRuleRead(ctx, d, m)
}

func resourceAppSignOnPolicyRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if isClassicOrg(ctx, m) {
		return resourceOIEOnlyFeatureError(appSignOnPolicyRule)
	}

	rule, resp, err := getAPISupplementFromMetadata(m).GetAppSignOnPolicyRule(ctx, d.Get("policy_id").(string), d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get app sign on policy rule: %v", err)
	}
	if rule == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("name", rule.Name)
	if rule.PriorityPtr != nil {
		_ = d.Set("priority", *rule.PriorityPtr)
	}
	_ = d.Set("status", rule.Status)
	if rule.Actions.AppSignOn != nil {
		_ = d.Set("access", rule.Actions.AppSignOn.Access)
		if rule.Actions.AppSignOn.VerificationMethod != nil {
			_ = d.Set("type", rule.Actions.AppSignOn.VerificationMethod.Type)
			_ = d.Set("factor_mode", rule.Actions.AppSignOn.VerificationMethod.FactorMode)
			_ = d.Set("re_authentication_frequency", rule.Actions.AppSignOn.VerificationMethod.ReauthenticateIn)
			_ = d.Set("inactivity_period", rule.Actions.AppSignOn.VerificationMethod.InactivityPeriod)
			arr := make([]interface{}, len(rule.Actions.AppSignOn.VerificationMethod.Constraints))
			for i := range rule.Actions.AppSignOn.VerificationMethod.Constraints {
				b, _ := json.Marshal(rule.Actions.AppSignOn.VerificationMethod.Constraints[i])
				arr[i] = string(b)
			}
			_ = d.Set("constraints", arr)
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
			m["network_includes"] = convertStringSliceToInterfaceSlice(rule.Conditions.Network.Include)
		}
		if len(rule.Conditions.Network.Exclude) > 0 {
			m["network_excludes"] = convertStringSliceToInterfaceSlice(rule.Conditions.Network.Exclude)
		}
		if rule.Conditions.Device != nil {
			_ = d.Set("device_is_managed", rule.Conditions.Device.Managed)
			_ = d.Set("device_is_registered", rule.Conditions.Device.Registered)
			m["device_assurances_included"] = convertStringSliceToSetNullable(rule.Conditions.Device.Assurance.Include)
		}
		if rule.Conditions.People != nil {
			if rule.Conditions.People.Users != nil {
				m["users_excluded"] = convertStringSliceToSetNullable(rule.Conditions.People.Users.Exclude)
				m["users_included"] = convertStringSliceToSetNullable(rule.Conditions.People.Users.Include)
			}
			if rule.Conditions.People.Groups != nil {
				m["groups_excluded"] = convertStringSliceToSetNullable(rule.Conditions.People.Groups.Exclude)
				m["groups_included"] = convertStringSliceToSetNullable(rule.Conditions.People.Groups.Include)
			}
		}
		if rule.Conditions.UserType != nil {
			m["user_types_excluded"] = convertStringSliceToSetNullable(rule.Conditions.UserType.Exclude)
			m["user_types_included"] = convertStringSliceToSetNullable(rule.Conditions.UserType.Include)
		}
		_ = setNonPrimitives(d, m)
	}
	return nil
}

func resourceAppSignOnPolicyRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if isClassicOrg(ctx, m) {
		return resourceOIEOnlyFeatureError(appSignOnPolicyRule)
	}

	_, _, err := getAPISupplementFromMetadata(m).UpdateAppSignOnPolicyRule(ctx, d.Get("policy_id").(string), d.Id(), buildAppSignOnPolicyRule(d))
	if err != nil {
		return diag.Errorf("failed to create app sign on policy rule: %v", err)
	}
	oldStatus, newStatus := d.GetChange("status")
	if oldStatus != newStatus {
		if newStatus == statusActive {
			_, err = getAPISupplementFromMetadata(m).ActivateAppSignOnPolicyRule(ctx, d.Get("policy_id").(string), d.Id())
		} else {
			_, err = getAPISupplementFromMetadata(m).DeactivateAppSignOnPolicyRule(ctx, d.Get("policy_id").(string), d.Id())
		}
		if err != nil {
			return diag.Errorf("failed to change app sign on policy rule status: %v", err)
		}
	}
	return resourceAppSignOnPolicyRuleRead(ctx, d, m)
}

func resourceAppSignOnPolicyRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if isClassicOrg(ctx, m) {
		return resourceOIEOnlyFeatureError(appSignOnPolicyRule)
	}

	if d.Get("name") == "Catch-all Rule" {
		// You cannot delete a default rule in a policy
		return nil
	}
	_, err := getAPISupplementFromMetadata(m).DeleteAppSignOnPolicyRule(ctx, d.Get("policy_id").(string), d.Id())
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
		PriorityPtr: int64Ptr(d.Get("priority").(int)),
		Type:        "ACCESS_POLICY",
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
	// if this is a default rule, the conditions attribute is read-only.
	if d.Get("name") == "Catch-all Rule" {
		return rule
	}
	rule.Conditions = &sdk.AccessPolicyRuleConditions{
		Network: buildPolicyNetworkCondition(d),
		Platform: &sdk.PlatformPolicyRuleCondition{
			Include: buildAccessPolicyPlatformInclude(d),
		},
		ElCondition: &sdk.AccessPolicyRuleCustomCondition{
			Condition: d.Get("custom_expression").(string),
		},
	}
	isRegistered, ok := d.GetOk("device_is_registered")
	if ok && isRegistered.(bool) {
		rule.Conditions.Device = &sdk.DeviceAccessPolicyRuleCondition{
			Managed:    boolPtr(d.Get("device_is_managed").(bool)),
			Registered: boolPtr(isRegistered.(bool)),
		}
	}
	deviceAssurancesIncluded, deviceAssurancesIncludedOk := d.GetOk("device_assurances_included")
	if deviceAssurancesIncludedOk {
		if rule.Conditions.Device != nil {
			rule.Conditions.Device.Assurance = &sdk.DeviceAssurancePolicyRuleCondition{
				Include: convertInterfaceToStringSetNullable(deviceAssurancesIncluded),
			}
		} else {
			rule.Conditions.Device = &sdk.DeviceAccessPolicyRuleCondition{
				Assurance: &sdk.DeviceAssurancePolicyRuleCondition{
					Include: convertInterfaceToStringSetNullable(deviceAssurancesIncluded),
				},
			}
		}
	}

	usersExcluded, usersExcludedOk := d.GetOk("users_excluded")
	usersIncluded, usersIncludedOk := d.GetOk("users_included")
	if usersExcludedOk || usersIncludedOk {
		rule.Conditions.People = &sdk.PolicyPeopleCondition{
			Users: &sdk.UserCondition{
				Exclude: convertInterfaceToStringSetNullable(usersExcluded),
				Include: convertInterfaceToStringSetNullable(usersIncluded),
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
			Exclude: convertInterfaceToStringSetNullable(groupsExcluded),
			Include: convertInterfaceToStringSetNullable(groupsIncluded),
		}
	}
	userTypesExcluded, userTypesExcludedOk := d.GetOk("user_types_excluded")
	userTypesIncluded, userTypesIncludedOk := d.GetOk("user_types_included")
	if userTypesExcludedOk || userTypesIncludedOk {
		rule.Conditions.UserType = &sdk.UserTypeCondition{
			Exclude: convertInterfaceToStringSetNullable(userTypesExcluded),
			Include: convertInterfaceToStringSetNullable(userTypesIncluded),
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
			var expr string
			if typ := getMapString(value, "os_type"); typ == "OTHER" {
				if v := getMapString(value, "os_expression"); v != "" {
					expr = v
				}
			}
			includeList = append(includeList, &sdk.PlatformConditionEvaluatorPlatform{
				Os: &sdk.PlatformConditionEvaluatorPlatformOperatingSystem{
					Expression: expr,
					Type:       getMapString(value, "os_type"),
				},
				Type: getMapString(value, "type"),
			})
		}
	}
	return includeList
}

func flattenAccessPolicyPlatformInclude(platform *sdk.PlatformPolicyRuleCondition) *schema.Set {
	var flattened []interface{}
	if platform != nil && platform.Include != nil {
		for _, v := range platform.Include {
			var expr string
			if v.Os.Expression != "" {
				expr = v.Os.Expression
			}
			flattened = append(flattened, map[string]interface{}{
				"os_expression": expr,
				"os_type":       v.Os.Type,
				"type":          v.Type,
			})
		}
	}
	return schema.NewSet(schema.HashResource(platformIncludeResource), flattened)
}
