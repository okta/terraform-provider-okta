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
		Description:   "Creates an MFA Policy Rule. This resource allows you to create and configure an MFA Policy Rule.",
		Schema: buildRuleSchema(map[string]*schema.Schema{
			"enroll": {
				Type:        schema.TypeString,
				Default:     "CHALLENGE",
				Optional:    true,
				Description: "When a user should be prompted for MFA. It can be `CHALLENGE`, `LOGIN`, or `NEVER`.",
			},
			"app_include": {
				Type:     schema.TypeSet,
				Elem:     appResource,
				Optional: true,
				Description: `Applications to include in discovery rule. **IMPORTANT**: this field is only available in Classic Organizations.
	- 'id' - (Optional) Use if 'type' is 'APP' to indicate the application id to include.
	- 'name' - (Optional) Use if the 'type' is 'APP_TYPE' to indicate the type of application(s) to include in instances where an entire group (i.e. 'yahoo_mail') of applications should be included.
	- 'type' - (Required) One of: 'APP', 'APP_TYPE'`,
			},
			"app_exclude": {
				Type:     schema.TypeSet,
				Elem:     appResource,
				Optional: true,
				Description: `Applications to exclude in discovery rule. **IMPORTANT**: this field is only available in Classic Organizations.
	- 'id' - (Optional) Use if 'type' is 'APP' to indicate the application id to include.
	- 'name' - (Optional) Use if the 'type' is 'APP_TYPE' to indicate the type of application(s) to include in instances where an entire group (i.e. 'yahoo_mail') of applications should be included.
	- 'type' - (Required) One of: 'APP', 'APP_TYPE'`,
			},
		}),
	}
}

func resourcePolicyMfaRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	template := buildMfaPolicyRule(d)
	err := createRule(ctx, d, meta, template, policyRulePassword)
	if err != nil {
		return diag.Errorf("failed to create MFA policy rule: %v", err)
	}
	return resourcePolicyMfaRuleRead(ctx, d, meta)
}

func resourcePolicyMfaRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rule, err := getPolicyRule(ctx, d, meta)
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

func resourcePolicyMfaRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	template := buildMfaPolicyRule(d)
	err := updateRule(ctx, d, meta, template)
	if err != nil {
		return diag.Errorf("failed to update MFA policy rule: %v", err)
	}
	return resourcePolicyMfaRuleRead(ctx, d, meta)
}

func resourcePolicyMfaRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := deleteRule(ctx, d, meta, false)
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
