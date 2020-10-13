package okta

import (
	"fmt"
	"strings"

	articulateOkta "github.com/articulate/oktasdk-go/okta"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const (
	passwordPolicyType = "PASSWORD"
	signOnPolicyType   = "OKTA_SIGN_ON"
	mfaPolicyType      = "MFA_ENROLL"
	idpDiscovery       = "IDP_DISCOVERY"
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

func createPolicyRuleImporter() *schema.ResourceImporter {
	return &schema.ResourceImporter{
		State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
			parts := strings.Split(d.Id(), "/")
			if len(parts) != 2 {
				return nil, fmt.Errorf("Invalid policy rule specifier. Expecting {policyID}/{ruleID}")
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

func getNetwork(d *schema.ResourceData) *articulateOkta.Network {
	return &articulateOkta.Network{
		Connection: d.Get("network_connection").(string),
		Exclude:    convertInterfaceToStringArrNullable(d.Get("network_excludes")),
		Include:    convertInterfaceToStringArrNullable(d.Get("network_includes")),
	}
}

func getPolicyRule(d *schema.ResourceData, m interface{}) (*articulateOkta.Rule, error) {
	client := m.(*Config).articulateOktaClient
	policyID := d.Get("policyid").(string)

	_, resp, err := client.Policies.GetPolicy(policyID)

	if is404(resp.StatusCode) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error Listing Policy in Okta: %v", err)
	}

	rule, resp, err := client.Policies.GetPolicyRule(policyID, d.Id())
	if is404(resp.StatusCode) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error Listing Policy Rule in Okta: %v", err)
	}

	return rule, nil
}

func getUsers(d *schema.ResourceData) *articulateOkta.People {
	var people *articulateOkta.People

	if include, ok := d.GetOk("users_excluded"); ok {
		people = &articulateOkta.People{
			Users: &articulateOkta.Users{
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

func syncRuleFromUpstream(d *schema.ResourceData, rule *articulateOkta.Rule) error {
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

	err = policyRuleActivate(d, meta)

	return rule, err
}
