package okta

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
		Type:         schema.TypeString,
		Optional:     true,
		Default:      statusActive,
		ValidateFunc: validation.StringInSlice([]string{statusActive, statusInactive}, false),
		Description:  "Policy Rule Status: ACTIVE or INACTIVE.",
	},
	"network_connection": {
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: validation.StringInSlice([]string{"ANYWHERE", "ZONE", "ON_NETWORK", "OFF_NETWORK"}, false),
		Description:  "Network selection mode: ANYWHERE, ZONE, ON_NETWORK, or OFF_NETWORK.",
		Default:      "ANYWHERE",
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

func createRule(d *schema.ResourceData, meta interface{}, template sdk.PolicyRule, ruleType string) (*sdk.PolicyRule, error) {
	ctx := context.Background()
	policyID := d.Get("policyid").(string)
	client := getSupplementFromMetadata(meta)
	_, resp, err := client.GetPolicy(ctx, policyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get policy by ID: %v", err)
	}
	if resp != nil && is404(resp.StatusCode) {
		return nil, fmt.Errorf("policy with ID %v not found ID", policyID)
	}
	rules, _, err := client.ListPolicyRules(ctx, policyID)
	if err != nil {
		return nil, fmt.Errorf("failed to list policy rules: %v", err)
	}
	ruleName := d.Get("name").(string)
	for i := range rules {
		if rules[i].Name == ruleName {
			return nil, fmt.Errorf("policy rule %v already exists in Okta. Please use 'import' to import it into terrafrom. terraform import %s.%s %s/%s", rules[i].Name, ruleType, rules[i].Name, policyID, rules[i].Id)
		}
	}
	rule, _, err := client.CreatePolicyRule(ctx, policyID, template)
	if err != nil {
		return nil, fmt.Errorf("failed to create policy rule: %v", err)
	}
	return rule, err
}

func createPolicyRuleImporter() *schema.ResourceImporter {
	return &schema.ResourceImporter{
		State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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

func getPolicyRule(d *schema.ResourceData, m interface{}) (*sdk.PolicyRule, error) {
	ctx := context.Background()
	client := getSupplementFromMetadata(m)
	policyID := d.Get("policyid").(string)

	_, resp, err := client.GetPolicy(ctx, policyID)
	if resp != nil && is404(resp.StatusCode) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get policy by ID: %v", err)
	}
	rule, resp, err := client.GetPolicyRule(ctx, policyID, d.Id())
	if resp != nil && is404(resp.StatusCode) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get policy rule by ID: %v", err)
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

func resourcePolicyRuleExists(d *schema.ResourceData, m interface{}) (b bool, e error) {
	// Exists - This is called to verify a resource still exists. It is called prior to Read,
	// and lowers the burden of Read to be able to assume the resource exists.
	policy, err := getPolicyRule(d, m)

	if err != nil || policy == nil {
		return false, err
	}

	return true, nil
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

func updateRule(d *schema.ResourceData, meta interface{}, updatedRule sdk.PolicyRule) (*sdk.PolicyRule, error) {
	client := getSupplementFromMetadata(meta)
	_, err := getPolicyRule(d, meta)
	if err != nil {
		return nil, err
	}

	rule, _, err := client.UpdatePolicyRule(context.Background(), d.Get("policyid").(string), d.Id(), updatedRule)
	if err != nil {
		return nil, fmt.Errorf("failed to update policy rule: %v", err)
	}

	err = policyRuleActivate(d, meta)

	return rule, err
}

func deleteRule(d *schema.ResourceData, m interface{}) error {
	if err := ensureNotDefaultRule(d); err != nil {
		return err
	}
	log.Printf("[INFO] Delete Policy Rule %v", d.Get("name").(string))
	client := getOktaClientFromMetadata(m)

	_, err := client.Policy.DeletePolicyRule(context.Background(), d.Get("policyid").(string), d.Id())
	if err != nil {
		return fmt.Errorf("[ERROR] Error Deleting Policy Rule from Okta: %v", err)
	}

	// remove the policy rule resource from terraform
	d.SetId("")

	return nil
}
