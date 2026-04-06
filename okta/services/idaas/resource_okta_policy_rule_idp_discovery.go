package idaas

import (
	"context"
	"fmt"

	"github.com/okta/terraform-provider-okta/okta/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	v6okta "github.com/okta/okta-sdk-golang/v6/okta"
)

func resourcePolicyRuleIdpDiscovery() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyRuleIdpDiscoveryCreate,
		ReadContext:   resourcePolicyRuleIdpDiscoveryRead,
		UpdateContext: resourcePolicyRuleIdpDiscoveryUpdate,
		DeleteContext: resourcePolicyRuleIdpDiscoveryDelete,
		Importer:      createPolicyRuleImporter(),
		Description: `Creates an IdP Discovery Policy Rule.

This resource allows you to create and configure an IdP Discovery Policy Rule.
-> If you receive the error 'You do not have permission to access the feature
you are requesting' [contact support](mailto:dev-inquiries@okta.com) and
request feature flag 'ADVANCED_SSO' be applied to your org.`,
		Schema: buildBaseRuleSchema(map[string]*schema.Schema{
			"idp_providers": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The identifier for the Idp the rule should route to if all conditions are met.",
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "Type of IdP. One of: `AMAZON`, `APPLE`, `DISCORD`, `FACEBOOK`, `GITHUB`, `GITLAB`, " +
								"`GOOGLE`, `IDV_CLEAR`, `IDV_INCODE`, `IDV_PERSONA`, `LINKEDIN`, `LOGINGOV`, `LOGINGOV_SANDBOX`, " +
								"`MICROSOFT`, `OIDC`, `PAYPAL`, `PAYPAL_SANDBOX`, `SALESFORCE`, `SAML2`, `SPOTIFY`, `X509`, `XERO`, " +
								"`YAHOO`, `YAHOOJP`, Default: `OKTA`",
						},
					},
				},
			},
			"app_include": {
				Type:     schema.TypeSet,
				Elem:     appResource,
				Optional: true,
				Description: `Applications to include in discovery rule.
- 'id' - (Optional) Use if 'type' is 'APP' to indicate the application id to include.
- 'name' - (Optional) Use if the 'type' is 'APP_TYPE' to indicate the type of application(s) to include in instances where an entire group (i.e. 'yahoo_mail') of applications should be included.
- 'type' - (Required) One of: 'APP', 'APP_TYPE'`,
			},
			"app_exclude": {
				Type:        schema.TypeSet,
				Elem:        appResource,
				Optional:    true,
				Description: "Applications to exclude in discovery. See `app_include` for details.",
			},
			"platform_include": {
				Type:     schema.TypeSet,
				Elem:     platformIncludeResource,
				Optional: true,
				Description: `Platform to include in discovery rule.
- 'type' - (Optional) One of: 'ANY', 'MOBILE', 'DESKTOP'
- 'os_expression - (Optional) Only available when using os_type = 'OTHER'
- 'os_type' - (Optional) One of: 'ANY', 'IOS', 'WINDOWS', 'ANDROID', 'OTHER', 'OSX'`,
			},
			"user_identifier_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "One of: `IDENTIFIER`, `ATTRIBUTE`",
			},
			"user_identifier_attribute": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Profile attribute matching can only have a single value that describes the type indicated in `user_identifier_type`. This is the attribute or identifier that the `user_identifier_patterns` are checked against.",
			},
			"user_identifier_patterns": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     userIDPatternResource,
				Description: `Specifies a User Identifier pattern condition to match against. If 'match_type' of 'EXPRESSION' is used, only a *single* element can be set, otherwise multiple elements of matching patterns may be provided.
- 'match_type' - (Optional) The kind of pattern. For regex, use 'EXPRESSION'. For simple string matches, use one of the following: 'SUFFIX', 'EQUALS', 'STARTS_WITH', 'CONTAINS'
- 'value' - (Optional) The regex or simple match string to match against.`,
			},
			"selection_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "SPECIFIC",
				Description: "Determines how the IdP is selected. One of: `SPECIFIC`, `DYNAMIC`. Default: `SPECIFIC`. When `DYNAMIC`, the IdP is selected based on the evaluated `provider_expression`.",
			},
			"provider_expression": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "An Okta Expression Language expression that is evaluated against the Login Context and used to dynamically select an IdP. " +
					"Only applicable when `selection_type` is `DYNAMIC`. Maps to `actions.idp.matchCriteria[0].providerExpression` in the API. " +
					"Example: `login.identifier.substringAfter('@')`",
			},
			"property_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The IdP property to match the evaluated expression against when `selection_type` is `DYNAMIC`. Maps to `actions.idp.matchCriteria[0].propertyName` in the API.",
			},
			"should_fall_back_to_okta": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Specifies whether to fall back to Okta if authentication with the matched IdP fails. Only applicable when `selection_type` is `DYNAMIC`. Default: `false`.",
			},
		}),
	}
}

func resourcePolicyRuleIdpDiscoveryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := validatePolicyRuleIdpDiscovery(d)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := d.Get("policy_id").(string)
	if policyID == "" {
		return diag.Errorf("'policy_id' field should be set")
	}
	logger(meta).Info("creating IdP discovery policy rule", "policy_id", policyID)

	newRule := buildIdpDiscoveryRule(d)
	ruleRequest := v6okta.IdpDiscoveryPolicyRuleAsListPolicyRules200ResponseInner(newRule)
	ruleResp, _, err := getOktaV6ClientFromMetadata(meta).PolicyAPI.CreatePolicyRule(ctx, policyID).PolicyRule(ruleRequest).Execute()
	if err != nil {
		return diag.Errorf("failed to create IDP discovery policy rule: %v", err)
	}
	if ruleResp == nil || ruleResp.IdpDiscoveryPolicyRule == nil {
		return diag.Errorf("API response did not contain a valid IdP discovery policy rule")
	}

	rule := ruleResp.IdpDiscoveryPolicyRule
	d.SetId(rule.GetId())

	err = setRuleStatus(ctx, d, meta, rule.GetStatus())
	if err != nil {
		return diag.Errorf("failed to set IDP discovery policy rule status: %v", err)
	}
	return resourcePolicyRuleIdpDiscoveryRead(ctx, d, meta)
}

func resourcePolicyRuleIdpDiscoveryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	policyID := d.Get("policy_id").(string)
	if policyID == "" {
		return diag.Errorf("'policy_id' field should be set")
	}
	logger(meta).Info("reading IdP discovery policy rule", "id", d.Id(), "policy_id", policyID)

	ruleResp, resp, err := getOktaV6ClientFromMetadata(meta).PolicyAPI.GetPolicyRule(ctx, policyID, d.Id()).Execute()
	if err := utils.SuppressErrorOn404_V6(resp, err); err != nil {
		return diag.Errorf("failed to get IDP discovery policy rule: %v", err)
	}
	if ruleResp == nil || ruleResp.IdpDiscoveryPolicyRule == nil {
		d.SetId("")
		return nil
	}

	rule := ruleResp.IdpDiscoveryPolicyRule
	_ = d.Set("name", rule.GetName())
	_ = d.Set("status", rule.GetStatus())
	if rule.Priority.IsSet() {
		_ = d.Set("priority", int(rule.GetPriority()))
	}

	mm := map[string]interface{}{}
	conditions := rule.Conditions
	if conditions != nil {
		if conditions.UserIdentifier != nil {
			_ = d.Set("user_identifier_attribute", conditions.UserIdentifier.GetAttribute())
			_ = d.Set("user_identifier_type", conditions.UserIdentifier.GetType())
			mm["user_identifier_patterns"] = flattenUserIDPatterns(conditions.UserIdentifier.GetPatterns())
		}
		if conditions.Platform != nil {
			mm["platform_include"] = flattenPlatformInclude(conditions.Platform)
		}
		if conditions.App != nil {
			mm["app_include"] = flattenDiscoveryRuleAppInclude(conditions.App)
			mm["app_exclude"] = flattenDiscoveryRuleAppExclude(conditions.App)
		}
		if conditions.Network != nil {
			_ = d.Set("network_connection", conditions.Network.GetConnection())
			if len(conditions.Network.Include) > 0 {
				mm["network_includes"] = utils.ConvertStringSliceToInterfaceSlice(conditions.Network.Include)
			}
			if len(conditions.Network.Exclude) > 0 {
				mm["network_excludes"] = utils.ConvertStringSliceToInterfaceSlice(conditions.Network.Exclude)
			}
		}
	}
	if setErr := utils.SetNonPrimitives(d, mm); setErr != nil {
		return diag.Errorf("failed to set IDP discovery policy rule properties: %v", setErr)
	}

	if rule.Actions != nil && rule.Actions.Idp != nil {
		idp := rule.Actions.Idp
		selectionType := idp.GetIdpSelectionType()
		if selectionType == "" {
			selectionType = "SPECIFIC"
		}
		_ = d.Set("selection_type", selectionType)
		_ = d.Set("should_fall_back_to_okta", idp.GetShouldFallBackToOkta())

		if criteria := idp.MatchCriteria; len(criteria) > 0 {
			_ = d.Set("provider_expression", criteria[0].GetProviderExpression())
			_ = d.Set("property_name", criteria[0].GetPropertyName())
		} else {
			_ = d.Set("provider_expression", "")
			_ = d.Set("property_name", "")
		}

		_ = d.Set("idp_providers", flattenDiscoveryRuleIdpProviders(idp.Providers))
	}

	return nil
}

func flattenDiscoveryRuleIdpProviders(providers []v6okta.IdpPolicyRuleActionProvider) []interface{} {
	var flattenedIdpProviders []interface{}
	for _, p := range providers {
		provider := make(map[string]interface{})

		// Default type to OKTA if missing
		if p.GetType() == "" {
			provider["type"] = "OKTA"
		} else {
			provider["type"] = p.GetType()
		}

		// Include ID only if not empty
		if p.GetId() != "" {
			provider["id"] = p.GetId()
		} else {
			provider["id"] = nil
		}

		flattenedIdpProviders = append(flattenedIdpProviders, provider)
	}
	return flattenedIdpProviders
}

func resourcePolicyRuleIdpDiscoveryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := validatePolicyRuleIdpDiscovery(d)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := d.Get("policy_id").(string)
	if policyID == "" {
		return diag.Errorf("'policy_id' field should be set")
	}
	logger(meta).Info("updating IdP discovery policy rule", "id", d.Id(), "policy_id", policyID)

	updatedRule := buildIdpDiscoveryRule(d)
	ruleID := d.Id()
	updatedRule.Id = &ruleID
	ruleRequest := v6okta.IdpDiscoveryPolicyRuleAsListPolicyRules200ResponseInner(updatedRule)
	ruleResp, _, err := getOktaV6ClientFromMetadata(meta).PolicyAPI.ReplacePolicyRule(ctx, policyID, d.Id()).PolicyRule(ruleRequest).Execute()
	if err != nil {
		return diag.Errorf("failed to update IDP discovery policy rule: %v", err)
	}
	if ruleResp == nil || ruleResp.IdpDiscoveryPolicyRule == nil {
		return diag.Errorf("API response did not contain a valid IdP discovery policy rule after update")
	}

	err = setRuleStatus(ctx, d, meta, ruleResp.IdpDiscoveryPolicyRule.GetStatus())
	if err != nil {
		return diag.Errorf("failed to set IDP discovery policy rule status: %v", err)
	}
	return resourcePolicyRuleIdpDiscoveryRead(ctx, d, meta)
}

func resourcePolicyRuleIdpDiscoveryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	policyID := d.Get("policy_id").(string)
	if policyID == "" {
		return diag.Errorf("'policy_id' field should be set")
	}
	logger(meta).Info("deleting IdP discovery policy rule", "id", d.Id(), "policy_id", policyID)

	resp, err := getOktaV6ClientFromMetadata(meta).PolicyAPI.DeletePolicyRule(ctx, policyID, d.Id()).Execute()
	if err := utils.SuppressErrorOn404_V6(resp, err); err != nil {
		return diag.Errorf("failed to delete IDP discovery policy rule: %v", err)
	}
	return nil
}

func setRuleStatus(ctx context.Context, d *schema.ResourceData, meta interface{}, status string) error {
	desiredStatus := d.Get("status").(string)
	if status == desiredStatus {
		return nil
	}
	policyID := d.Get("policy_id").(string)
	if policyID == "" {
		return fmt.Errorf("'policy_id' field should be set")
	}
	logger(meta).Info("setting IdP discovery policy rule status", "id", d.Id(),
		"policy_id", policyID, "status", desiredStatus)

	var err error
	client := getOktaV6ClientFromMetadata(meta).PolicyAPI
	if desiredStatus == StatusInactive {
		_, err = client.DeactivatePolicyRule(ctx, policyID, d.Id()).Execute()
	} else {
		_, err = client.ActivatePolicyRule(ctx, policyID, d.Id()).Execute()
	}
	return err
}

// buildIdpDiscoveryRule constructs a v6 IdpDiscoveryPolicyRule from the resource data.
func buildIdpDiscoveryRule(d *schema.ResourceData) *v6okta.IdpDiscoveryPolicyRule {
	rule := v6okta.NewIdpDiscoveryPolicyRule()

	name := d.Get("name").(string)
	rule.Name = &name

	ruleType := "IDP_DISCOVERY"
	rule.Type = &ruleType

	status := d.Get("status").(string)
	rule.Status = &status

	if priority, ok := d.GetOk("priority"); ok {
		p := int32(priority.(int))
		rule.Priority = *v6okta.NewNullableInt32(&p)
	}

	// Build IDP action
	idp := v6okta.NewIdpPolicyRuleActionIdp()

	selectionType := d.Get("selection_type").(string)
	if selectionType != "" {
		idp.SetIdpSelectionType(selectionType)
	}

	shouldFallBack := d.Get("should_fall_back_to_okta").(bool)
	idp.SetShouldFallBackToOkta(shouldFallBack)

	if selectionType == "DYNAMIC" {
		// DYNAMIC rules use matchCriteria for IdP selection; providers list is empty
		idp.SetProviders([]v6okta.IdpPolicyRuleActionProvider{})
		if expr := d.Get("provider_expression").(string); expr != "" {
			criteria := v6okta.NewIdpPolicyRuleActionMatchCriteria()
			criteria.SetProviderExpression(expr)
			if propName := d.Get("property_name").(string); propName != "" {
				criteria.SetPropertyName(propName)
			}
			idp.SetMatchCriteria([]v6okta.IdpPolicyRuleActionMatchCriteria{*criteria})
		}
	} else {
		// SPECIFIC mode: use the providers list
		var providers []v6okta.IdpPolicyRuleActionProvider
		if v, ok := d.GetOk("idp_providers"); ok {
			for _, provider := range v.([]interface{}) {
				if value, ok := provider.(map[string]any); ok {
					p := v6okta.NewIdpPolicyRuleActionProvider()
					if id := utils.GetMapString(value, "id"); id != "" {
						p.SetId(id)
					}
					if t := utils.GetMapString(value, "type"); t != "" {
						p.SetType(t)
					}
					providers = append(providers, *p)
				}
			}
		}
		if len(providers) == 0 {
			p := v6okta.NewIdpPolicyRuleActionProvider()
			p.SetType("OKTA")
			providers = []v6okta.IdpPolicyRuleActionProvider{*p}
		}
		idp.SetProviders(providers)
	}

	action := v6okta.NewIdpPolicyRuleAction()
	action.SetIdp(*idp)
	rule.SetActions(*action)

	// Build conditions
	conditions := v6okta.NewIdpDiscoveryPolicyRuleCondition()

	// Network condition
	network := v6okta.NewPolicyNetworkCondition()
	network.SetConnection(d.Get("network_connection").(string))
	if includes := utils.ConvertInterfaceToStringArr(d.Get("network_includes")); len(includes) > 0 {
		network.SetInclude(includes)
	}
	if excludes := utils.ConvertInterfaceToStringArr(d.Get("network_excludes")); len(excludes) > 0 {
		network.SetExclude(excludes)
	}
	conditions.SetNetwork(*network)

	// Platform condition
	if platform := buildPlatformIncludeV6(d); platform != nil {
		conditions.SetPlatform(*platform)
	}

	// App condition
	if app := buildAppConditionsV6(d); app != nil {
		conditions.SetApp(*app)
	}

	// User identifier condition
	if uid := buildIdentifierV6(d); uid != nil {
		conditions.SetUserIdentifier(*uid)
	}

	rule.SetConditions(*conditions)

	return rule
}

var (
	platformIncludeResource = &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ANY",
			},
			"os_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ANY",
			},
			"os_expression": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Only available with OTHER OS type",
			},
		},
	}

	userIDPatternResource = &schema.Resource{
		Schema: map[string]*schema.Schema{
			"match_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"value": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
)

func buildPlatformIncludeV6(d *schema.ResourceData) *v6okta.IdpDiscoveryPlatformPolicyRuleCondition {
	v, ok := d.GetOk("platform_include")
	if !ok {
		return nil
	}
	var includeList []v6okta.IdpDiscoveryPlatformConditionEvaluatorPlatform
	for _, item := range v.(*schema.Set).List() {
		if value, ok := item.(map[string]interface{}); ok {
			platform := v6okta.NewIdpDiscoveryPlatformConditionEvaluatorPlatform()
			if t := utils.GetMapString(value, "type"); t != "" {
				platform.SetType(t)
			}
			os := v6okta.NewIdpDiscoveryPlatformConditionEvaluatorPlatformOperatingSystem()
			if osType := utils.GetMapString(value, "os_type"); osType != "" {
				os.SetType(osType)
			}
			if osExpr := utils.GetMapString(value, "os_expression"); osExpr != "" {
				os.SetExpression(osExpr)
			}
			platform.SetOs(*os)
			includeList = append(includeList, *platform)
		}
	}
	result := v6okta.NewIdpDiscoveryPlatformPolicyRuleCondition()
	result.SetInclude(includeList)
	return result
}

func buildAppConditionsV6(d *schema.ResourceData) *v6okta.AppAndInstancePolicyRuleCondition {
	var includeList []v6okta.AppAndInstanceConditionEvaluatorAppOrInstance
	if v, ok := d.GetOk("app_include"); ok {
		for _, item := range v.(*schema.Set).List() {
			if value, ok := item.(map[string]interface{}); ok {
				app := v6okta.NewAppAndInstanceConditionEvaluatorAppOrInstance()
				if id := utils.GetMapString(value, "id"); id != "" {
					app.SetId(id)
				}
				if name := utils.GetMapString(value, "name"); name != "" {
					app.SetName(name)
				}
				if t := utils.GetMapString(value, "type"); t != "" {
					app.SetType(t)
				}
				includeList = append(includeList, *app)
			}
		}
	}

	var excludeList []v6okta.AppAndInstanceConditionEvaluatorAppOrInstance
	if v, ok := d.GetOk("app_exclude"); ok {
		for _, item := range v.(*schema.Set).List() {
			if value, ok := item.(map[string]interface{}); ok {
				app := v6okta.NewAppAndInstanceConditionEvaluatorAppOrInstance()
				if id := utils.GetMapString(value, "id"); id != "" {
					app.SetId(id)
				}
				if name := utils.GetMapString(value, "name"); name != "" {
					app.SetName(name)
				}
				if t := utils.GetMapString(value, "type"); t != "" {
					app.SetType(t)
				}
				excludeList = append(excludeList, *app)
			}
		}
	}

	if len(includeList) == 0 && len(excludeList) == 0 {
		return nil
	}

	cond := v6okta.NewAppAndInstancePolicyRuleCondition()
	if len(includeList) > 0 {
		cond.SetInclude(includeList)
	}
	if len(excludeList) > 0 {
		cond.SetExclude(excludeList)
	}
	return cond
}

func buildUserIDPatternsV6(d *schema.ResourceData) []v6okta.UserIdentifierConditionEvaluatorPattern {
	var patternList []v6okta.UserIdentifierConditionEvaluatorPattern
	if raw, ok := d.GetOk("user_identifier_patterns"); ok {
		for _, pattern := range raw.(*schema.Set).List() {
			if value, ok := pattern.(map[string]interface{}); ok {
				p := v6okta.NewUserIdentifierConditionEvaluatorPattern()
				p.SetMatchType(utils.GetMapString(value, "match_type"))
				p.SetValue(utils.GetMapString(value, "value"))
				patternList = append(patternList, *p)
			}
		}
	}
	return patternList
}

func buildIdentifierV6(d *schema.ResourceData) *v6okta.UserIdentifierPolicyRuleCondition {
	uidType := d.Get("user_identifier_type").(string)
	if uidType == "" {
		return nil
	}
	patterns := buildUserIDPatternsV6(d)
	uid := v6okta.NewUserIdentifierPolicyRuleCondition()
	uid.SetType(uidType)
	uid.SetPatterns(patterns)
	attribute := d.Get("user_identifier_attribute").(string)
	if attribute != "" {
		uid.SetAttribute(attribute)
	}
	return uid
}

func flattenUserIDPatterns(patterns []v6okta.UserIdentifierConditionEvaluatorPattern) *schema.Set {
	flattened := make([]interface{}, len(patterns))
	for i, p := range patterns {
		flattened[i] = map[string]interface{}{
			"match_type": p.GetMatchType(),
			"value":      p.GetValue(),
		}
	}
	return schema.NewSet(schema.HashResource(userIDPatternResource), flattened)
}

func flattenPlatformInclude(platform *v6okta.IdpDiscoveryPlatformPolicyRuleCondition) *schema.Set {
	var flattened []interface{}
	if platform != nil {
		for _, v := range platform.Include {
			item := map[string]interface{}{
				"type":          v.GetType(),
				"os_type":       "",
				"os_expression": "",
			}
			if v.Os != nil {
				item["os_type"] = v.Os.GetType()
				item["os_expression"] = v.Os.GetExpression()
			}
			flattened = append(flattened, item)
		}
	}
	return schema.NewSet(schema.HashResource(platformIncludeResource), flattened)
}

func flattenDiscoveryRuleAppInclude(app *v6okta.AppAndInstancePolicyRuleCondition) *schema.Set {
	if app != nil {
		return flattenAppObj(app.Include)
	}
	return flattenAppObj(nil)
}

func flattenDiscoveryRuleAppExclude(app *v6okta.AppAndInstancePolicyRuleCondition) *schema.Set {
	if app != nil {
		return flattenAppObj(app.Exclude)
	}
	return flattenAppObj(nil)
}

func flattenAppObj(appObj []v6okta.AppAndInstanceConditionEvaluatorAppOrInstance) *schema.Set {
	var flattened []interface{}
	for _, v := range appObj {
		flattened = append(flattened, map[string]interface{}{
			"id":   v.GetId(),
			"name": v.GetName(),
			"type": v.GetType(),
		})
	}
	return schema.NewSet(schema.HashResource(appResource), flattened)
}

var (
	errFDiscoveryRuleIdPAppConditionID   = "either 'name' or 'id' should be provided in the '%s' block"
	errFDiscoveryRuleIdPAppConditionName = "'name' is required if the type is 'APP_TYPE' in the '%s' block"
)

func validatePolicyRuleIdpDiscovery(d *schema.ResourceData) error {
	for _, appCondition := range []string{"app_include", "app_exclude"} {
		v, ok := d.GetOk(appCondition)
		if !ok {
			continue
		}
		for _, item := range v.(*schema.Set).List() {
			if value, ok := item.(map[string]interface{}); ok {
				id := utils.GetMapString(value, "id")
				name := utils.GetMapString(value, "name")
				if id == "" && name == "" {
					return fmt.Errorf(errFDiscoveryRuleIdPAppConditionID, appCondition)
				}
				if utils.GetMapString(value, "type") == "APP_TYPE" && name == "" {
					return fmt.Errorf(errFDiscoveryRuleIdPAppConditionName, appCondition)
				}
			}
		}
	}
	return nil
}
