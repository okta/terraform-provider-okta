package okta

import (
	"fmt"
	"strings"

	articulateOkta "github.com/articulate/oktasdk-go/okta"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

const passwordPolicyType = "PASSWORD"
const signOnPolicyType = "OKTA_SIGN_ON"
const singOnPolicyRuleType = "SIGN_ON"

// Basis of policy rules
var baseRuleSchema = map[string]*schema.Schema{
	"policyid": &schema.Schema{
		Type:        schema.TypeString,
		ForceNew:    true,
		Required:    true,
		Description: "Policy ID of the Rule",
	},
	"name": &schema.Schema{
		Type:        schema.TypeString,
		ForceNew:    true,
		Required:    true,
		Description: "Policy Rule Name",
	},
	"priority": &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "Policy Rule Priority, this attribute can be set to a valid priority. To avoid endless diff situation we error if an invalid priority is provided. API defaults it to the last/lowest if not there.",
		// Suppress diff if config is empty.
		DiffSuppressFunc: createValueDiffSuppression("0"),
	},
	"status": &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Default:      "ACTIVE",
		ValidateFunc: validation.StringInSlice([]string{"ACTIVE", "INACTIVE"}, false),
		Description:  "Policy Rule Status: ACTIVE or INACTIVE.",
	},
	"users_excluded": {
		Type:          schema.TypeList,
		Optional:      true,
		Description:   "List of User IDs to Exclude",
		ConflictsWith: []string{"users_included"},
		Elem:          &schema.Schema{Type: schema.TypeString},
	},
	"users_included": {
		Type:          schema.TypeList,
		Optional:      true,
		Description:   "List of User IDs to Exclude",
		ConflictsWith: []string{"users_excluded"},
		Elem:          &schema.Schema{Type: schema.TypeString},
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

func buildRuleSchema(target map[string]*schema.Schema) map[string]*schema.Schema {
	for k := range baseRuleSchema {
		target[k] = baseRuleSchema[k]
	}

	return target
}

func syncRuleFromUpstream(d *schema.ResourceData, rule *articulateOkta.Rule) error {
	d.Set("name", rule.Name)
	d.Set("status", rule.Status)
	d.Set("priority", rule.Priority)
	d.Set("network_connection", rule.Conditions.Network.Connection)

	return setNonPrimitives(d, map[string]interface{}{
		"users_excluded":   convertStringArrToInterface(rule.Conditions.People.Users.Exclude),
		"users_included":   convertStringArrToInterface(rule.Conditions.People.Users.Include),
		"network_includes": convertStringArrToInterface(rule.Conditions.Network.Include),
		"network_excludes": convertStringArrToInterface(rule.Conditions.Network.Exclude),
	})
}

func getPolicyRule(d *schema.ResourceData, m interface{}) (*articulateOkta.Rule, error) {
	client := m.(*Config).articulateOktaClient
	policyID := d.Get("policyid").(string)

	_, _, err := client.Policies.GetPolicy(policyID)

	if is404(client) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error Listing Policy in Okta: %v", err)
	}

	rule, _, err := client.Policies.GetPolicyRule(policyID, d.Id())
	if is404(client) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error Listing Policy Rule in Okta: %v", err)
	}

	return rule, nil
}

func getNetwork(d *schema.ResourceData) *articulateOkta.Network {
	network := &articulateOkta.Network{
		Connection: d.Get("network_connection").(string),
	}

	include := d.Get("network_includes")
	exclude := d.Get("network_excludes")

	if include != nil {
		network.Include = convertInterfaceToStringArr(include)
	} else if exclude != nil {
		network.Exclude = convertInterfaceToStringArr(exclude)
	}

	return network
}

func getUsers(d *schema.ResourceData) *articulateOkta.People {
	var people *articulateOkta.People
	include := d.Get("users_included")
	exclude := d.Get("users_excluded")

	if include != nil {
		people = &articulateOkta.People{
			Users: &articulateOkta.Users{
				Include: convertInterfaceToStringArr(include),
			},
		}
	} else if exclude != nil {
		people = &articulateOkta.People{
			Users: &articulateOkta.Users{
				Exclude: convertInterfaceToStringArr(exclude),
			},
		}
	}

	return people
}

func updateRule(d *schema.ResourceData, meta interface{}, updatedRule interface{}) (*articulateOkta.Rule, error) {
	client := getClientFromMetadata(meta)
	_, err := getPolicyRule(d, meta)
	if err != nil {
		return nil, err
	}

	rule, _, err := client.Policies.UpdatePolicyRule(d.Get("policyid").(string), d.Id(), updatedRule)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error runing update against Sign On Policy Rule: %v", err)
	}
	d.Partial(false)
	err = policyRuleActivate(d, meta)

	return rule, err
}

func createRule(d *schema.ResourceData, meta interface{}, template interface{}, ruleType string) (*articulateOkta.Rule, error) {
	client := getClientFromMetadata(meta)
	policyID := d.Get("policyid").(string)

	_, _, err := client.Policies.GetPolicy(policyID)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error Listing Policy in Okta: %v", err)
	}

	currentPolicyRules, _, err := client.Policies.GetPolicyRules(policyID)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error Listing Policy Rules in Okta: %v", err)
	}

	if currentPolicyRules != nil {
		for _, rule := range currentPolicyRules.Rules {
			ruleName := d.Get("name").(string)

			if rule.Name == ruleName {
				return nil, fmt.Errorf("policy rule %v already exists in Okta. Please use import to import it into terrafrom. terraform import %s.%s %s/%s", rule.Name, ruleType, rule.Name, policyID, rule.ID)
			}
		}
	}

	rule, _, err := client.Policies.CreatePolicyRule(policyID, template)
	if err != nil {
		return nil, err
	}

	return rule, err
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

func createPolicyRuleImporter() *schema.ResourceImporter {
	return &schema.ResourceImporter{
		State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
			parts := strings.Split(d.Id(), "/")
			if len(parts) != 2 {
				return nil, fmt.Errorf("Invalid policy rule specifier. Expecting {policyID}/{ruleID}")
			}
			d.Set("policyid", parts[0])
			d.SetId(parts[1])
			return []*schema.ResourceData{d}, nil
		},
	}
}
