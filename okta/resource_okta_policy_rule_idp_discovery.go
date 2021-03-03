package okta

import (
	"context"

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
		Schema: buildBaseRuleSchema(map[string]*schema.Schema{
			"idp_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"idp_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "OKTA",
			},
			"app_include": {
				Type:        schema.TypeSet,
				Elem:        appResource,
				Optional:    true,
				Description: "Applications to include in discovery rule",
			},
			"app_exclude": {
				Type:        schema.TypeSet,
				Elem:        appResource,
				Optional:    true,
				Description: "Applications to exclude in discovery rule",
			},
			"platform_include": {
				Type:     schema.TypeSet,
				Elem:     platformIncludeResource,
				Optional: true,
			},
			"user_identifier_type": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: stringInSlice([]string{"IDENTIFIER", "ATTRIBUTE", ""}),
			},
			"user_identifier_attribute": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_identifier_patterns": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     userIDPatternResource,
			},
		}),
	}
}

func resourcePolicyRuleIdpDiscoveryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	logger(m).Info("creating IdP discovery policy rule", "policy_id", d.Get("policyid").(string))
	newRule := buildIdpDiscoveryRule(d)
	rule, _, err := getSupplementFromMetadata(m).CreateIdpDiscoveryRule(ctx, d.Get("policyid").(string), *newRule, nil)
	if err != nil {
		return diag.Errorf("failed to create IDP discovery policy rule: %v", err)
	}
	d.SetId(rule.ID)
	err = setRuleStatus(ctx, d, m, rule.Status)
	if err != nil {
		return diag.Errorf("failed to set IDP discovery policy rule status: %v", err)
	}
	return resourcePolicyRuleIdpDiscoveryRead(ctx, d, m)
}

func resourcePolicyRuleIdpDiscoveryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	logger(m).Info("reading IdP discovery policy rule", "id", d.Id(), "policy_id", d.Get("policyid").(string))
	rule, resp, err := getSupplementFromMetadata(m).GetIdpDiscoveryRule(ctx, d.Get("policyid").(string), d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
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
	_ = d.Set("network_connection", rule.Conditions.Network.Connection)
	err = setNonPrimitives(d, map[string]interface{}{
		"network_includes":         convertStringArrToInterface(rule.Conditions.Network.Include),
		"network_excludes":         convertStringArrToInterface(rule.Conditions.Network.Exclude),
		"platform_include":         flattenPlatformInclude(rule.Conditions.Platform),
		"user_identifier_patterns": flattenUserIDPatterns(rule.Conditions.UserIdentifier.Patterns),
		"app_include":              flattenAppInclude(rule.Conditions.App),
		"app_exclude":              flattenAppExclude(rule.Conditions.App),
	})
	if err != nil {
		return diag.Errorf("failed to set IDP discovery policy rule properties: %v", err)
	}
	return nil
}

func resourcePolicyRuleIdpDiscoveryUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	logger(m).Info("updating IdP discovery policy rule", "id", d.Id(), "policy_id", d.Get("policyid").(string))
	newRule := buildIdpDiscoveryRule(d)
	rule, _, err := getSupplementFromMetadata(m).UpdateIdpDiscoveryRule(ctx, d.Get("policyid").(string), d.Id(), *newRule, nil)
	if err != nil {
		return diag.Errorf("failed to update IDP discovery policy rule: %v", err)
	}
	err = setRuleStatus(ctx, d, m, rule.Status)
	if err != nil {
		return diag.Errorf("failed to set IDP discovery policy rule status: %v", err)
	}
	return resourcePolicyRuleIdpDiscoveryRead(ctx, d, m)
}

func resourcePolicyRuleIdpDiscoveryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	logger(m).Info("deleting IdP discovery policy rule", "id", d.Id(), "policy_id", d.Get("policyid").(string))
	_, err := getOktaClientFromMetadata(m).Policy.DeletePolicyRule(ctx, d.Get("policyid").(string), d.Id())
	if err != nil {
		return diag.Errorf("failed to delete IDP discovery policy rule: %v", err)
	}
	return nil
}

func setRuleStatus(ctx context.Context, d *schema.ResourceData, m interface{}, status string) error {
	desiredStatus := d.Get("status").(string)
	if status == desiredStatus {
		return nil
	}
	logger(m).Info("setting IdP discovery policy rule status", "id", d.Id(),
		"policy_id", d.Get("policyid").(string), "status", desiredStatus)
	var err error
	client := getOktaClientFromMetadata(m)
	if desiredStatus == statusInactive {
		_, err = client.Policy.DeactivatePolicyRule(ctx, d.Get("policyid").(string), d.Id())
	} else {
		_, err = client.Policy.ActivatePolicyRule(ctx, d.Get("policyid").(string), d.Id())
	}
	return err
}

// Build Policy Sign On Rule from resource data
func buildIdpDiscoveryRule(d *schema.ResourceData) *sdk.IdpDiscoveryRule {
	rule := &sdk.IdpDiscoveryRule{
		Actions: &sdk.IdpDiscoveryRuleActions{
			IDP: &sdk.IdpDiscoveryRuleIdp{
				Providers: []*sdk.IdpDiscoveryRuleProvider{
					{
						Type: d.Get("idp_type").(string),
						ID:   d.Get("idp_id").(string),
					},
				},
			},
		},
		Conditions: &sdk.IdpDiscoveryRuleConditions{
			App: buildAppConditions(d),
			Network: &sdk.IdpDiscoveryRuleNetwork{
				Connection: d.Get("network_connection").(string),
				// plural name here is vestigial due to old policy rule resources
				Include: convertInterfaceToStringArr(d.Get("network_includes")),
				Exclude: convertInterfaceToStringArr(d.Get("network_excludes")),
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

var platformIncludeResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"type": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "ANY",
			ValidateDiagFunc: stringInSlice([]string{"ANY", "MOBILE", "DESKTOP"}),
		},
		"os_type": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "ANY",
			ValidateDiagFunc: stringInSlice([]string{"ANY", "IOS", "WINDOWS", "ANDROID", "OTHER", "OSX"}),
		},
		"os_expression": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Only available with OTHER OS type",
		},
	},
}

var userIDPatternResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"match_type": {
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: stringInSlice([]string{"SUFFIX", "EQUALS", "STARTS_WITH", "CONTAINS", "EXPRESSION"}),
		},
		"value": {
			Type:     schema.TypeString,
			Optional: true,
		},
	},
}

var appResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"type": {
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: stringInSlice([]string{"APP", "APP_TYPE"}),
		},
		"name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"id": {
			Type:     schema.TypeString,
			Optional: true,
		},
	},
}

func buildPlatformInclude(d *schema.ResourceData) *sdk.IdpDiscoveryRulePlatform {
	var includeList []*sdk.IdpDiscoveryRulePlatformInclude

	if v, ok := d.GetOk("platform_include"); ok {
		valueList := v.(*schema.Set).List()

		for _, item := range valueList {
			if value, ok := item.(map[string]interface{}); ok {
				includeList = append(includeList, &sdk.IdpDiscoveryRulePlatformInclude{
					Os: &sdk.IdpDiscoveryRulePlatformOS{
						Expression: getMapString(value, "os_expression"),
						Type:       getMapString(value, "os_type"),
					},
					Type: getMapString(value, "type"),
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
					ID:   getMapString(value, "id"),
					Type: getMapString(value, "type"),
					Name: getMapString(value, "name"),
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
					ID:   getMapString(value, "id"),
					Type: getMapString(value, "type"),
					Name: getMapString(value, "name"),
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
					MatchType: getMapString(value, "match_type"),
					Value:     getMapString(value, "value"),
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

func flattenAppInclude(app *sdk.IdpDiscoveryRuleApp) *schema.Set {
	var flattened []interface{}

	if app != nil && app.Include != nil {
		for _, v := range app.Include {
			flattened = append(flattened, map[string]interface{}{
				"id":   v.ID,
				"name": v.Name,
				"type": v.Type,
			})
		}
	}
	return schema.NewSet(schema.HashResource(appResource), flattened)
}

func flattenAppExclude(app *sdk.IdpDiscoveryRuleApp) *schema.Set {
	var flattened []interface{}

	if app != nil && app.Exclude != nil {
		for _, v := range app.Exclude {
			flattened = append(flattened, map[string]interface{}{
				"id":   v.ID,
				"name": v.Name,
				"type": v.Type,
			})
		}
	}
	return schema.NewSet(schema.HashResource(appResource), flattened)
}
