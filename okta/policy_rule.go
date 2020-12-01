package okta

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
)

var userExcludedSchema = map[string]*schema.Schema{
	"users_excluded": {
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "Set of User IDs to Exclude",
		Elem:        &schema.Schema{Type: schema.TypeString},
	},
}

// Basis of policy rules
var baseRuleSchema = map[string]*schema.Schema{
	// Ugh vestigial incorrect naming. Should switch to policy_id
	"policyid": {
		Type:        schema.TypeString,
		ForceNew:    true,
		Required:    true,
		Description: "Policy ID of the Rule",
	},
	"name": {
		Type:        schema.TypeString,
		ForceNew:    true,
		Required:    true,
		Description: "Policy Rule Name",
	},
	"priority": {
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "Policy Rule Priority, this attribute can be set to a valid priority. To avoid endless diff situation we error if an invalid priority is provided. API defaults it to the last/lowest if not there.",
		// Suppress diff if config is empty.
		DiffSuppressFunc: createValueDiffSuppression("0"),
	},
	"status": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          statusActive,
		ValidateDiagFunc: stringInSlice([]string{statusActive, statusInactive}),
		Description:      "Policy Rule Status: ACTIVE or INACTIVE.",
	},
	"network_connection": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: stringInSlice([]string{"ANYWHERE", "ZONE", "ON_NETWORK", "OFF_NETWORK"}),
		Description:      "Network selection mode: ANYWHERE, ZONE, ON_NETWORK, or OFF_NETWORK.",
		Default:          "ANYWHERE",
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
}

func buildBaseRuleSchema(target map[string]*schema.Schema) map[string]*schema.Schema {
	return buildSchema(baseRuleSchema, target)
}

func buildRuleSchema(target map[string]*schema.Schema) map[string]*schema.Schema {
	return buildSchema(buildSchema(baseRuleSchema, target), userExcludedSchema)
}

func createRule(ctx context.Context, d *schema.ResourceData, m interface{}, template sdk.PolicyRule, ruleType string) error {
	logger(m).Info("creating policy rule", "name", d.Get("name").(string))
	err := ensureNotDefaultRule(d)
	if err != nil {
		return err
	}
	policyID := d.Get("policyid").(string)
	client := getSupplementFromMetadata(m)
	_, resp, err := client.GetPolicy(ctx, policyID)
	if err != nil {
		return fmt.Errorf("failed to get policy by ID: %v", err)
	}
	if is404(resp) {
		return fmt.Errorf("policy with ID %v not found ID", policyID)
	}
	rules, _, err := client.ListPolicyRules(ctx, policyID)
	if err != nil {
		return fmt.Errorf("failed to list policy rules: %v", err)
	}
	ruleName := d.Get("name").(string)
	for i := range rules {
		if rules[i].Name == ruleName {
			return fmt.Errorf("policy rule %v already exists in Okta. Please use 'import' to import it into terrafrom. terraform import %s.%s %s/%s", rules[i].Name, ruleType, rules[i].Name, policyID, rules[i].Id)
		}
	}
	rule, _, err := client.CreatePolicyRule(ctx, policyID, template)
	if err != nil {
		return fmt.Errorf("failed to create policy rule: %v", err)
	}
	// We want to put this under Terraform's control even if priority is invalid.
	d.SetId(rule.Id)
	return validatePriority(template.Priority, rule.Priority)
}

func createPolicyRuleImporter() *schema.ResourceImporter {
	return &schema.ResourceImporter{
		StateContext: func(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
			parts := strings.Split(d.Id(), "/")
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid policy rule specifier. Expecting {policyID}/{ruleID}")
			}
			_ = d.Set("policyid", parts[0])
			d.SetId(parts[1])
			return []*schema.ResourceData{d}, nil
		},
	}
}

func ensureNotDefaultRule(d *schema.ResourceData) error {
	return ensureNotDefault(d, "Rule")
}

func getNetwork(d *schema.ResourceData) *okta.PolicyNetworkCondition {
	return &okta.PolicyNetworkCondition{
		Connection: d.Get("network_connection").(string),
		Exclude:    convertInterfaceToStringArrNullable(d.Get("network_excludes")),
		Include:    convertInterfaceToStringArrNullable(d.Get("network_includes")),
	}
}

func getPolicyRule(ctx context.Context, d *schema.ResourceData, m interface{}) (*sdk.PolicyRule, error) {
	client := getSupplementFromMetadata(m)
	policyID := d.Get("policyid").(string)
	policy, resp, err := client.GetPolicy(ctx, policyID)
	if err := suppressErrorOn404(resp, err); err != nil {
		return nil, err
	}
	if policy == nil {
		d.SetId("")
		return nil, nil
	}
	rule, resp, err := client.GetPolicyRule(ctx, policyID, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return nil, err
	}
	if rule == nil {
		d.SetId("")
		return nil, nil
	}
	return rule, nil
}

func getUsers(d *schema.ResourceData) *okta.PolicyPeopleCondition {
	var people *okta.PolicyPeopleCondition

	if include, ok := d.GetOk("users_excluded"); ok {
		people = &okta.PolicyPeopleCondition{
			Users: &okta.UserCondition{
				Exclude: convertInterfaceToStringSet(include),
			},
		}
	}

	return people
}

func syncRuleFromUpstream(d *schema.ResourceData, rule *sdk.PolicyRule) error {
	_ = d.Set("name", rule.Name)
	_ = d.Set("status", rule.Status)
	_ = d.Set("priority", rule.Priority)
	_ = d.Set("network_connection", rule.Conditions.Network.Connection)
	return setNonPrimitives(d, map[string]interface{}{
		"users_excluded":   convertStringSetToInterface(rule.Conditions.People.Users.Exclude),
		"network_includes": convertStringArrToInterface(rule.Conditions.Network.Include),
		"network_excludes": convertStringArrToInterface(rule.Conditions.Network.Exclude),
	})
}

func updateRule(ctx context.Context, d *schema.ResourceData, m interface{}, template sdk.PolicyRule) error {
	if err := ensureNotDefaultRule(d); err != nil {
		return err
	}
	logger(m).Info("updating policy rule", "name", d.Get("name").(string))
	client := getSupplementFromMetadata(m)
	rule, _, err := client.UpdatePolicyRule(ctx, d.Get("policyid").(string), d.Id(), template)
	if err != nil {
		return err
	}
	err = validatePriority(template.Priority, rule.Priority)
	if err != nil {
		return err
	}
	return policyRuleActivate(ctx, d, m)
}

// activate or deactivate a policy rule according to the terraform schema status field
func policyRuleActivate(ctx context.Context, d *schema.ResourceData, m interface{}) error {
	client := getSupplementFromMetadata(m)

	if d.Get("status").(string) == statusActive {
		_, err := client.ActivatePolicyRule(ctx, d.Get("policyid").(string), d.Id())
		if err != nil {
			return fmt.Errorf("activation has failed: %v", err)
		}
	}
	if d.Get("status").(string) == statusInactive {
		_, err := client.DeactivatePolicyRule(ctx, d.Get("policyid").(string), d.Id())
		if err != nil {
			return fmt.Errorf("deactivation has failed: %v", err)
		}
	}
	return nil
}

func deleteRule(ctx context.Context, d *schema.ResourceData, m interface{}, checkIsSystemPolicy bool) error {
	logger(m).Info("deleting policy rule", "name", d.Get("name").(string))
	if err := ensureNotDefaultRule(d); err != nil {
		return err
	}
	rule, err := getPolicyRule(ctx, d, m)
	if rule == nil {
		return nil
	}
	if err != nil {
		return err
	}
	shouldRemove := true
	if checkIsSystemPolicy {
		if rule.System != nil && *rule.System {
			logger(m).Info(fmt.Sprintf("Policy Rule '%s' is a System Policy, cannot delete from Okta", d.Get("name").(string)))
			shouldRemove = false
		}
	}
	if shouldRemove {
		_, err = getOktaClientFromMetadata(m).Policy.DeletePolicyRule(ctx, d.Get("policyid").(string), d.Id())
		if err != nil {
			return err
		}
	}
	// remove the policy rule resource from terraform
	d.SetId("")
	return nil
}
