package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourcePolicyMfaRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyMfaRuleCreate,
		ReadContext:   resourcePolicyMfaRuleRead,
		UpdateContext: resourcePolicyMfaRuleUpdate,
		DeleteContext: resourcePolicyMfaRuleDelete,
		Importer:      createPolicyRuleImporter(),
		Schema: buildRuleSchema(map[string]*schema.Schema{
			"enroll": {
				Type:        schema.TypeString,
				Default:     "CHALLENGE",
				Optional:    true,
				Description: "Should the user be enrolled the first time they LOGIN, the next time they are CHALLENGED, or NEVER?",
			},
			"app_include": {
				Type:        schema.TypeSet,
				Elem:        appResource,
				Optional:    true,
				Description: "Applications to include",
			},
			"app_exclude": {
				Type:        schema.TypeSet,
				Elem:        appResource,
				Optional:    true,
				Description: "Applications to exclude",
			},
		}),
	}
}

func resourcePolicyMfaRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	template := buildMfaPolicyRule(d)
	err := createRule(ctx, d, m, template, policyRulePassword)
	if err != nil {
		return diag.Errorf("failed to create MFA policy rule: %v", err)
	}
	return resourcePolicyMfaRuleRead(ctx, d, m)
}

func resourcePolicyMfaRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	rule, err := getPolicyRule(ctx, d, m)
	if err != nil {
		return diag.Errorf("failed to get MFA policy rule: %v", err)
	}
	if rule == nil {
		return nil
	}
	err = syncRuleFromUpstream(d, rule)
	if err != nil {
		return diag.Errorf("failed to sync MFA policy rule: %v", err)
	}
	if (rule.Conditions.App) != nil {
		if len(rule.Conditions.App.Include) != 0 {
			_ = d.Set("app_include", flattenApps(rule.Conditions.App.Include))
		}
		if len(rule.Conditions.App.Exclude) != 0 {
			_ = d.Set("app_exclude", flattenApps(rule.Conditions.App.Exclude))
		}
	}
	if rule.Actions.PasswordPolicyRuleActions != nil && rule.Actions.PasswordPolicyRuleActions.Enroll != nil {
		_ = d.Set("enroll", rule.Actions.PasswordPolicyRuleActions.Enroll.Self)
	}
	return nil
}

func resourcePolicyMfaRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	template := buildMfaPolicyRule(d)
	err := updateRule(ctx, d, m, template)
	if err != nil {
		return diag.Errorf("failed to update MFA policy rule: %v", err)
	}
	return resourcePolicyMfaRuleRead(ctx, d, m)
}

func resourcePolicyMfaRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := deleteRule(ctx, d, m, false)
	if err != nil {
		return diag.Errorf("failed to delete MFA policy rule: %v", err)
	}
	return nil
}

// build password policy rule from schema data
func buildMfaPolicyRule(d *schema.ResourceData) sdk.SdkPolicyRule {
	rule := sdk.MfaPolicyRule()
	rule.Name = d.Get("name").(string)
	rule.Status = d.Get("status").(string)
	if priority, ok := d.GetOk("priority"); ok {
		rule.Priority = int64(priority.(int))
	}
	rule.Conditions = &sdk.PolicyRuleConditions{
		Network: buildPolicyNetworkCondition(d),
		People:  getUsers(d),
		App:     buildMFAPolicyAppCondition(d),
	}
	if enroll, ok := d.GetOk("enroll"); ok {
		rule.Actions = sdk.SdkPolicyRuleActions{
			PasswordPolicyRuleActions: &sdk.PasswordPolicyRuleActions{
				Enroll: &sdk.PolicyRuleActionsEnroll{
					Self: enroll.(string),
				},
			},
		}
	}
	return rule
}

func buildMFAPolicyAppCondition(d *schema.ResourceData) *sdk.AppAndInstancePolicyRuleCondition {
	incl, okInclude := d.GetOk("app_include")
	excl, okExclude := d.GetOk("app_exclude")
	if !okInclude && !okExclude {
		return nil
	}
	rc := &sdk.AppAndInstancePolicyRuleCondition{}
	if okInclude {
		valueList := incl.(*schema.Set).List()
		var includeList []*sdk.AppAndInstanceConditionEvaluatorAppOrInstance
		for _, item := range valueList {
			if value, ok := item.(map[string]interface{}); ok {
				includeList = append(includeList, &sdk.AppAndInstanceConditionEvaluatorAppOrInstance{
					Id:   getMapString(value, "id"),
					Type: getMapString(value, "type"),
					Name: getMapString(value, "name"),
				})
			}
		}
		rc.Include = includeList
	}
	if okExclude {
		valueList := excl.(*schema.Set).List()
		var excludeList []*sdk.AppAndInstanceConditionEvaluatorAppOrInstance
		for _, item := range valueList {
			if value, ok := item.(map[string]interface{}); ok {
				excludeList = append(excludeList, &sdk.AppAndInstanceConditionEvaluatorAppOrInstance{
					Id:   getMapString(value, "id"),
					Type: getMapString(value, "type"),
					Name: getMapString(value, "name"),
				})
			}
		}
		rc.Exclude = excludeList
	}
	return rc
}

func flattenApps(appObj []*sdk.AppAndInstanceConditionEvaluatorAppOrInstance) *schema.Set {
	var flattened []interface{}
	for _, v := range appObj {
		flattened = append(flattened, map[string]interface{}{
			"id":   v.Id,
			"name": v.Name,
			"type": v.Type,
		})
	}
	return schema.NewSet(schema.HashResource(appResource), flattened)
}
