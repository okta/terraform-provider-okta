package okta

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
)

var platformIncludeResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"type": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "ANY",
			ValidateFunc: validation.StringInSlice([]string{"ANY", "MOBILE", "DESKTOP"}, false),
		},
		"os_type": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "ANY",
			ValidateFunc: validation.StringInSlice([]string{"ANY", "IOS", "WINDOWS", "ANDROID", "OTHER", "OSX"}, false),
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
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"SUFFIX", "EQUALS", "STARTS_WITH", "CONTAINS", "EXPRESSION"}, false),
		},
		"value": {
			Type:     schema.TypeString,
			Optional: true,
		},
	},
}

//https://developer.okta.com/docs/reference/api/policy/#application-and-app-instance-condition-object
var appResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"type": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"APP", "APP_TYPE"}, false),
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

func resourcePolicyRuleIdpDiscovery() *schema.Resource {
	return &schema.Resource{
		Exists:   resourcePolicyRuleIdpDiscoveryExists,
		Create:   resourcePolicyRuleIdpDiscoveryCreate,
		Read:     resourcePolicyRuleIdpDiscoveryRead,
		Update:   resourcePolicyRuleIdpDiscoveryUpdate,
		Delete:   resourcePolicyRuleIdpDiscoveryDelete,
		Importer: createPolicyRuleImporter(),

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
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"IDENTIFIER", "ATTRIBUTE", ""}, false),
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
	includeList := []*sdk.IdpDiscoveryRuleAppObj{}

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

	excludeList := []*sdk.IdpDiscoveryRuleAppObj{}

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
	flattened := []interface{}{}

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
	flattened := []interface{}{}

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

func resourcePolicyRuleIdpDiscoveryExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := getSupplementFromMetadata(m)
	rule, _, err := client.GetIdpDiscoveryRule(d.Get("policyid").(string), d.Id())

	return err == nil && rule.ID != "", err
}

func resourcePolicyRuleIdpDiscoveryCreate(d *schema.ResourceData, m interface{}) error {
	newRule := buildIdpDiscoveryRule(d)
	client := getSupplementFromMetadata(m)
	rule, resp, err := client.CreateIdpDiscoveryRule(d.Get("policyid").(string), *newRule, nil)
	if err != nil {
		return responseErr(resp, err)
	}

	d.SetId(rule.ID)
	err = setRuleStatus(d, m, rule.Status)
	if err != nil {
		return err
	}

	return resourcePolicyRuleIdpDiscoveryRead(d, m)
}

func setRuleStatus(d *schema.ResourceData, m interface{}, status string) error {
	desiredStatus := d.Get("status").(string)

	if status != desiredStatus {
		client := getSupplementFromMetadata(m)
		if desiredStatus == statusInactive {
			return responseErr(client.DeactivateRule(d.Get("policyid").(string), d.Id()))
		} else if desiredStatus == statusActive {
			return responseErr(client.ActivateRule(d.Get("policyid").(string), d.Id()))
		}
	}

	return nil
}

func resourcePolicyRuleIdpDiscoveryRead(d *schema.ResourceData, m interface{}) error {
	client := getSupplementFromMetadata(m)
	rule, resp, err := client.GetIdpDiscoveryRule(d.Get("policyid").(string), d.Id())

	if resp != nil && is404(resp.StatusCode) {
		d.SetId("")
		return nil
	}

	if err != nil {
		return responseErr(resp, err)
	}

	_ = d.Set("name", rule.Name)
	_ = d.Set("status", rule.Status)
	_ = d.Set("priority", rule.Priority)
	_ = d.Set("user_identifier_attribute", rule.Conditions.UserIdentifier.Attribute)
	_ = d.Set("user_identifier_type", rule.Conditions.UserIdentifier.Type)
	_ = d.Set("network_connection", rule.Conditions.Network.Connection)

	return setNonPrimitives(d, map[string]interface{}{
		"network_includes":         convertStringArrToInterface(rule.Conditions.Network.Include),
		"network_excludes":         convertStringArrToInterface(rule.Conditions.Network.Exclude),
		"platform_include":         flattenPlatformInclude(rule.Conditions.Platform),
		"user_identifier_patterns": flattenUserIDPatterns(rule.Conditions.UserIdentifier.Patterns),
		"app_include":              flattenAppInclude(rule.Conditions.App),
		"app_exclude":              flattenAppExclude(rule.Conditions.App),
	})
}

func resourcePolicyRuleIdpDiscoveryUpdate(d *schema.ResourceData, m interface{}) error {
	newRule := buildIdpDiscoveryRule(d)
	client := getSupplementFromMetadata(m)
	rule, resp, err := client.UpdateIdpDiscoveryRule(d.Get("policyid").(string), d.Id(), *newRule, nil)
	if err != nil {
		return responseErr(resp, err)
	}

	err = setRuleStatus(d, m, rule.Status)
	if err != nil {
		return err
	}

	return resourcePolicyRuleIdpDiscoveryRead(d, m)
}

func resourcePolicyRuleIdpDiscoveryDelete(d *schema.ResourceData, m interface{}) error {
	client := getSupplementFromMetadata(m)

	resp, err := client.DeleteIdpDiscoveryRule(d.Get("policyid").(string), d.Id())
	return suppressErrorOn404(resp, err)
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
