package idaas

import (
	"context"
	"fmt"

	"github.com/okta/terraform-provider-okta/okta/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
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
	rule, _, err := getAPISupplementFromMetadata(meta).CreateIdpDiscoveryRule(ctx, policyID, *newRule, nil)
	if err != nil {
		return diag.Errorf("failed to create IDP discovery policy rule: %v", err)
	}
	d.SetId(rule.ID)
	err = setRuleStatus(ctx, d, meta, rule.Status)
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
	rule, resp, err := getAPISupplementFromMetadata(meta).GetIdpDiscoveryRule(ctx, policyID, d.Id())
	if err := utils.SuppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get IDP discovery policy rule: %v", err)
	}
	if rule == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("name", rule.Name)
	_ = d.Set("status", rule.Status)
	_ = d.Set("priority", rule.Priority)
	_ = d.Set("user_identifier_attribute", rule.Conditions.UserIdentifier.Attribute)
	_ = d.Set("user_identifier_type", rule.Conditions.UserIdentifier.Type)
	mm := map[string]interface{}{
		"platform_include":         flattenPlatformInclude(rule.Conditions.Platform),
		"app_include":              flattenDiscoveryRuleAppInclude(rule.Conditions.App),
		"app_exclude":              flattenDiscoveryRuleAppExclude(rule.Conditions.App),
		"user_identifier_patterns": flattenUserIDPatterns(rule.Conditions.UserIdentifier.Patterns),
	}
	_ = d.Set("network_connection", rule.Conditions.Network.Connection)
	if len(rule.Conditions.Network.Include) > 0 {
		mm["network_includes"] = utils.ConvertStringSliceToInterfaceSlice(rule.Conditions.Network.Include)
	}
	if len(rule.Conditions.Network.Exclude) > 0 {
		mm["network_excludes"] = utils.ConvertStringSliceToInterfaceSlice(rule.Conditions.Network.Exclude)
	}
	err = utils.SetNonPrimitives(d, mm)
	if err != nil {
		return diag.Errorf("failed to set IDP discovery policy rule properties: %v", err)
	}
	d.Set("idp_providers", flattenDiscoveryRuleIdpProviders(rule.Actions.IDP.Providers))
	return nil
}

func flattenDiscoveryRuleIdpProviders(providers []*sdk.IdpDiscoveryRuleProvider) []interface{} {
	var flattenedIdpProviders []interface{}
	for _, p := range providers {
		provider := make(map[string]interface{})

		// Default type to OKTA if missing
		if p.Type == "" {
			provider["type"] = "OKTA"
		} else {
			provider["type"] = p.Type
		}

		// Include ID only if not empty
		if p.ID != "" {
			provider["id"] = p.ID
		} else {
			// Optional: explicitly set nil to ensure Terraform understands it should remove `id`
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
	newRule := buildIdpDiscoveryRule(d)
	rule, _, err := getAPISupplementFromMetadata(meta).UpdateIdpDiscoveryRule(ctx, policyID, d.Id(), *newRule, nil)
	if err != nil {
		return diag.Errorf("failed to update IDP discovery policy rule: %v", err)
	}
	err = setRuleStatus(ctx, d, meta, rule.Status)
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
	_, err := getOktaClientFromMetadata(meta).Policy.DeletePolicyRule(ctx, policyID, d.Id())
	if err != nil {
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
	client := getOktaClientFromMetadata(meta)
	if desiredStatus == StatusInactive {
		_, err = client.Policy.DeactivatePolicyRule(ctx, policyID, d.Id())
	} else {
		_, err = client.Policy.ActivatePolicyRule(ctx, policyID, d.Id())
	}
	return err
}

// Build Policy Sign On Rule from resource data
func buildIdpDiscoveryRule(d *schema.ResourceData) *sdk.IdpDiscoveryRule {
	var providers []*sdk.IdpDiscoveryRuleProvider
	if v, ok := d.GetOk("idp_providers"); ok {
		providerList := v.([]interface{})
		for _, provider := range providerList {
			if value, ok := provider.(map[string]any); ok {
				providers = append(providers, &sdk.IdpDiscoveryRuleProvider{
					ID:   utils.GetMapString(value, "id"),
					Type: utils.GetMapString(value, "type"),
				})
			}
		}
	} else {
		providers = append(providers, &sdk.IdpDiscoveryRuleProvider{
			Type: "OKTA",
		})
	}

	rule := &sdk.IdpDiscoveryRule{
		Actions: &sdk.IdpDiscoveryRuleActions{
			IDP: &sdk.IdpDiscoveryRuleIdp{
				Providers: providers,
			},
		},
		Conditions: &sdk.IdpDiscoveryRuleConditions{
			App: buildAppConditions(d),
			Network: &sdk.IdpDiscoveryRuleNetwork{
				Connection: d.Get("network_connection").(string),
				// plural name here is vestigial due to old policy rule resources
				Include: utils.ConvertInterfaceToStringArr(d.Get("network_includes")),
				Exclude: utils.ConvertInterfaceToStringArr(d.Get("network_excludes")),
			},
			Platform:       buildPlatformInclude(d),
			UserIdentifier: buildIdentifier(d),
		},
		Type:   sdk.IdpDiscoveryType,
		Name:   d.Get("name").(string),
		Status: d.Get("status").(string),
	}
	if priority, ok := d.GetOk("priority"); ok {
		rule.Priority = priority.(int)
	}
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

func buildPlatformInclude(d *schema.ResourceData) *sdk.IdpDiscoveryRulePlatform {
	var includeList []*sdk.IdpDiscoveryRulePlatformInclude
	if v, ok := d.GetOk("platform_include"); ok {
		valueList := v.(*schema.Set).List()
		for _, item := range valueList {
			if value, ok := item.(map[string]interface{}); ok {
				includeList = append(includeList, &sdk.IdpDiscoveryRulePlatformInclude{
					Os: &sdk.IdpDiscoveryRulePlatformOS{
						Expression: utils.GetMapString(value, "os_expression"),
						Type:       utils.GetMapString(value, "os_type"),
					},
					Type: utils.GetMapString(value, "type"),
				})
			}
		}
		return &sdk.IdpDiscoveryRulePlatform{
			Include: includeList,
		}
	}
	return nil
}

func buildAppConditions(d *schema.ResourceData) *sdk.IdpDiscoveryRuleApp {
	var includeList []*sdk.IdpDiscoveryRuleAppObj
	if v, ok := d.GetOk("app_include"); ok {
		valueList := v.(*schema.Set).List()
		for _, item := range valueList {
			if value, ok := item.(map[string]interface{}); ok {
				includeList = append(includeList, &sdk.IdpDiscoveryRuleAppObj{
					ID:   utils.GetMapString(value, "id"),
					Type: utils.GetMapString(value, "type"),
					Name: utils.GetMapString(value, "name"),
				})
			}
		}
	}
	var excludeList []*sdk.IdpDiscoveryRuleAppObj
	if v, ok := d.GetOk("app_exclude"); ok {
		valueList := v.(*schema.Set).List()
		for _, item := range valueList {
			if value, ok := item.(map[string]interface{}); ok {
				excludeList = append(excludeList, &sdk.IdpDiscoveryRuleAppObj{
					ID:   utils.GetMapString(value, "id"),
					Type: utils.GetMapString(value, "type"),
					Name: utils.GetMapString(value, "name"),
				})
			}
		}
	}
	return &sdk.IdpDiscoveryRuleApp{
		Include: includeList,
		Exclude: excludeList,
	}
}

func buildUserIDPatterns(d *schema.ResourceData) []*sdk.IdpDiscoveryRulePattern {
	var patternList []*sdk.IdpDiscoveryRulePattern
	if raw, ok := d.GetOk("user_identifier_patterns"); ok {
		patterns := raw.(*schema.Set).List()
		for _, pattern := range patterns {
			if value, ok := pattern.(map[string]interface{}); ok {
				patternList = append(patternList, &sdk.IdpDiscoveryRulePattern{
					MatchType: utils.GetMapString(value, "match_type"),
					Value:     utils.GetMapString(value, "value"),
				})
			}
		}
	}
	return patternList
}

func buildIdentifier(d *schema.ResourceData) *sdk.IdpDiscoveryRuleUserIdentifier {
	uidType := d.Get("user_identifier_type").(string)
	if uidType != "" {
		return &sdk.IdpDiscoveryRuleUserIdentifier{
			Attribute: d.Get("user_identifier_attribute").(string),
			Type:      uidType,
			Patterns:  buildUserIDPatterns(d),
		}
	}
	return nil
}

func flattenUserIDPatterns(patterns []*sdk.IdpDiscoveryRulePattern) *schema.Set {
	flattened := make([]interface{}, len(patterns))
	for i := range patterns {
		flattened[i] = map[string]interface{}{
			"match_type": patterns[i].MatchType,
			"value":      patterns[i].Value,
		}
	}
	return schema.NewSet(schema.HashResource(userIDPatternResource), flattened)
}

func flattenPlatformInclude(platform *sdk.IdpDiscoveryRulePlatform) *schema.Set {
	var flattened []interface{}
	if platform != nil && platform.Include != nil {
		for _, v := range platform.Include {
			flattened = append(flattened, map[string]interface{}{
				"os_expression": v.Os.Expression,
				"os_type":       v.Os.Type,
				"type":          v.Type,
			})
		}
	}
	return schema.NewSet(schema.HashResource(platformIncludeResource), flattened)
}

func flattenDiscoveryRuleAppInclude(app *sdk.IdpDiscoveryRuleApp) *schema.Set {
	if app != nil {
		return flattenAppObj(app.Include)
	}
	return flattenAppObj(nil)
}

func flattenDiscoveryRuleAppExclude(app *sdk.IdpDiscoveryRuleApp) *schema.Set {
	if app != nil {
		return flattenAppObj(app.Exclude)
	}
	return flattenAppObj(nil)
}

func flattenAppObj(appObj []*sdk.IdpDiscoveryRuleAppObj) *schema.Set {
	var flattened []interface{}
	for _, v := range appObj {
		flattened = append(flattened, map[string]interface{}{
			"id":   v.ID,
			"name": v.Name,
			"type": v.Type,
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
