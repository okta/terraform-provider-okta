package okta

import (
	"fmt"

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
	// Need to add an enhancement here, maybe we can provide a hierarchy system, priorityBefore = <refRule>
	"priority": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Policy Rule Priority",
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

	_, _, err := client.Policies.GetPolicy(d.Get("policyid").(string))

	if client.OktaErrorCode == "E0000007" {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error Listing Policy in Okta: %v", err)
	}

	rule, _, err := client.Policies.GetPolicyRule(d.Get("policyid").(string), d.Id())
	if client.OktaErrorCode == "E0000007" {
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
