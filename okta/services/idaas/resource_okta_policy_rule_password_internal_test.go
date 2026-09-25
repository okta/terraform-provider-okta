package idaas

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	v6okta "github.com/okta/okta-sdk-golang/v6/okta"
)

// currentPasswordRuleJSON is a GET /api/v1/policies/{policyId}/rules/{ruleId}
// response for a password rule. Alongside the properties the schema models it
// carries two Okta does not: actions.selfServicePasswordReset.settings, the case
// that exposed GH-2960, and a top level unmodeledFutureProperty standing in for
// whatever Okta adds next.
const currentPasswordRuleJSON = `{
  "id": "0pr1a2b3c4d5e6f7g8h9",
  "name": "Example Rule",
  "type": "PASSWORD",
  "status": "ACTIVE",
  "priority": 1,
  "system": false,
  "created": "2026-01-02T03:04:05.000Z",
  "lastUpdated": "2026-01-02T03:04:05.000Z",
  "unmodeledFutureProperty": "keep me",
  "conditions": {
    "network": {
      "connection": "ANYWHERE",
      "unmodeledNetworkProperty": "keep me"
    },
    "people": {
      "users": {
        "exclude": ["00u1a2b3c4d5e6f7g8h9"]
      }
    }
  },
  "actions": {
    "passwordChange": {
      "access": "ALLOW"
    },
    "selfServiceUnlock": {
      "access": "DENY"
    },
    "selfServicePasswordReset": {
      "access": "ALLOW",
      "settings": {
        "allowRecoveryEmailWithoutEnrollment": true
      },
      "requirement": {
        "accessControl": "LEGACY",
        "primary": {
          "methods": ["email"]
        },
        "stepUp": {
          "required": false
        }
      }
    }
  },
  "_links": {
    "self": {
      "href": "https://example.okta.com/api/v1/policies/00p1a2b3c4d5e6f7g8h9/rules/0pr1a2b3c4d5e6f7g8h9"
    }
  }
}`

func currentPasswordRule(t *testing.T) *v6okta.PasswordPolicyRule {
	t.Helper()
	var rule v6okta.PasswordPolicyRule
	if err := json.Unmarshal([]byte(currentPasswordRuleJSON), &rule); err != nil {
		t.Fatalf("failed to unmarshal the upstream rule fixture: %v", err)
	}
	return &rule
}

// baseRuleConfig is the configuration equivalent of currentPasswordRuleJSON,
// minus the properties the schema does not model. password_unlock differs so
// that tests can tell a managed change apart from a preserved one.
func baseRuleConfig() map[string]interface{} {
	return map[string]interface{}{
		"policy_id":                     "00p1a2b3c4d5e6f7g8h9",
		"name":                          "Example Rule",
		"status":                        StatusActive,
		"priority":                      1,
		"network_connection":            "ANYWHERE",
		"password_change":               "ALLOW",
		"password_reset":                "ALLOW",
		"password_unlock":               "ALLOW",
		"users_excluded":                []interface{}{"00u1a2b3c4d5e6f7g8h9"},
		"password_reset_access_control": "LEGACY",
		"password_reset_requirement": []interface{}{
			map[string]interface{}{
				"primary_methods": []interface{}{"email"},
				"step_up_enabled": false,
			},
		},
	}
}

// buildRuleFromConfig runs the real build path: raw configuration in, the struct
// the update path would PUT out.
func buildRuleFromConfig(t *testing.T, config map[string]interface{}) v6okta.PasswordPolicyRule {
	t.Helper()
	d := schema.TestResourceDataRaw(t, resourcePolicyPasswordRule().Schema, config)
	return buildPolicyRulePassword(d)
}

// replacePolicyRuleBody marshals a rule the way the SDK marshals it on the way
// into ReplacePolicyRule, so assertions run against the bytes Okta would see.
func replacePolicyRuleBody(t *testing.T, rule *v6okta.PasswordPolicyRule) map[string]interface{} {
	t.Helper()
	payload := v6okta.PasswordPolicyRuleAsListPolicyRules200ResponseInner(rule)
	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal the rule: %v", err)
	}
	var body map[string]interface{}
	if err := json.Unmarshal(raw, &body); err != nil {
		t.Fatalf("failed to unmarshal the request body: %v", err)
	}
	return body
}

// lookup walks a decoded JSON body. It reports whether the path exists, so that
// tests can assert on a property being gone as well as on its value.
func lookup(body map[string]interface{}, path ...string) (interface{}, bool) {
	var current interface{} = body
	for _, key := range path {
		object, ok := current.(map[string]interface{})
		if !ok {
			return nil, false
		}
		current, ok = object[key]
		if !ok {
			return nil, false
		}
	}
	return current, true
}

func assertValue(t *testing.T, body map[string]interface{}, want interface{}, path ...string) {
	t.Helper()
	got, ok := lookup(body, path...)
	if !ok {
		t.Fatalf("%v is missing from the request body, expected %v", path, want)
	}
	if got != want {
		t.Errorf("%v = %v, expected %v", path, got, want)
	}
}

func assertAbsent(t *testing.T, body map[string]interface{}, path ...string) {
	t.Helper()
	if got, ok := lookup(body, path...); ok {
		t.Errorf("%v = %v, expected it to be absent from the request body", path, got)
	}
}

// TestPreserveUnmodeledPropertiesKeepsUpstreamProperties is the GH-2960
// regression test. An update must not delete server side properties the schema
// does not model.
func TestPreserveUnmodeledPropertiesKeepsUpstreamProperties(t *testing.T) {
	rule := buildRuleFromConfig(t, baseRuleConfig())
	preserveUnmodeledProperties(&rule, currentPasswordRule(t))
	body := replacePolicyRuleBody(t, &rule)

	// The property that exposed the defect, and the whole object holding it.
	assertValue(t, body, true, "actions", "selfServicePasswordReset", "settings", "allowRecoveryEmailWithoutEnrollment")
	// Unknown properties on any other node survive too.
	assertValue(t, body, "keep me", "unmodeledFutureProperty")
	assertValue(t, body, "keep me", "conditions", "network", "unmodeledNetworkProperty")

	// Configuration still wins for everything the schema does model.
	assertValue(t, body, "ALLOW", "actions", "selfServiceUnlock", "access")
	assertValue(t, body, "ALLOW", "actions", "passwordChange", "access")
	assertValue(t, body, "LEGACY", "actions", "selfServicePasswordReset", "requirement", "accessControl")
	assertValue(t, body, "Example Rule", "name")
	assertValue(t, body, "ANYWHERE", "conditions", "network", "connection")

	// Read-only properties are not echoed back: they belong to Okta, and the
	// build path deliberately leaves them out of the payload.
	for _, property := range []string{"id", "created", "lastUpdated", "_links"} {
		assertAbsent(t, body, property)
	}
}

// TestPreserveUnmodeledPropertiesKeepsRemovalsRemoved guards the other half of
// the contract: carrying unknown properties over must not resurrect a managed
// value the configuration dropped.
func TestPreserveUnmodeledPropertiesKeepsRemovalsRemoved(t *testing.T) {
	config := baseRuleConfig()
	// Drop everything the upstream rule has that this configuration no longer asks for.
	delete(config, "users_excluded")
	delete(config, "password_reset_access_control")
	delete(config, "password_reset_requirement")

	rule := buildRuleFromConfig(t, config)
	preserveUnmodeledProperties(&rule, currentPasswordRule(t))
	body := replacePolicyRuleBody(t, &rule)

	assertAbsent(t, body, "conditions", "people")
	assertAbsent(t, body, "actions", "selfServicePasswordReset", "requirement")

	// An unknown property under a node configuration kept is still carried over.
	assertValue(t, body, true, "actions", "selfServicePasswordReset", "settings", "allowRecoveryEmailWithoutEnrollment")
}

// TestPreserveUnmodeledPropertiesPerNode checks every node of the rule that can
// hold unknown properties, including the ones a plain configuration does not
// reach, so a new nesting level in the SDK cannot quietly go unhandled.
func TestPreserveUnmodeledPropertiesPerNode(t *testing.T) {
	// A configuration that populates every node the build path can produce.
	config := baseRuleConfig()
	config["users_included"] = []interface{}{"00u1a2b3c4d5e6f7g8h9"}
	config["groups_included"] = []interface{}{"00g1a2b3c4d5e6f7g8h9"}
	config["network_connection"] = "ZONE"
	config["network_includes"] = []interface{}{"nzo1a2b3c4d5e6f7g8h9"}
	config["password_reset_requirement"] = []interface{}{
		map[string]interface{}{
			"primary_methods": []interface{}{"email"},
			"step_up_enabled": true,
			"step_up_methods": []interface{}{"security_question"},
		},
	}

	// Each case names a node and the JSON path the property should come out at.
	tests := []struct {
		name string
		set  func(current *v6okta.PasswordPolicyRule, properties map[string]interface{})
		path []string
	}{
		{
			name: "rule",
			set:  func(c *v6okta.PasswordPolicyRule, p map[string]interface{}) { c.AdditionalProperties = p },
			path: []string{},
		},
		{
			name: "conditions",
			set:  func(c *v6okta.PasswordPolicyRule, p map[string]interface{}) { c.Conditions.AdditionalProperties = p },
			path: []string{"conditions"},
		},
		{
			name: "conditions.network",
			set: func(c *v6okta.PasswordPolicyRule, p map[string]interface{}) {
				c.Conditions.Network.AdditionalProperties = p
			},
			path: []string{"conditions", "network"},
		},
		{
			name: "conditions.people",
			set: func(c *v6okta.PasswordPolicyRule, p map[string]interface{}) {
				c.Conditions.People.AdditionalProperties = p
			},
			path: []string{"conditions", "people"},
		},
		{
			name: "conditions.people.users",
			set: func(c *v6okta.PasswordPolicyRule, p map[string]interface{}) {
				c.Conditions.People.Users.AdditionalProperties = p
			},
			path: []string{"conditions", "people", "users"},
		},
		{
			name: "conditions.people.groups",
			set: func(c *v6okta.PasswordPolicyRule, p map[string]interface{}) {
				c.Conditions.People.Groups.AdditionalProperties = p
			},
			path: []string{"conditions", "people", "groups"},
		},
		{
			name: "actions",
			set:  func(c *v6okta.PasswordPolicyRule, p map[string]interface{}) { c.Actions.AdditionalProperties = p },
			path: []string{"actions"},
		},
		{
			name: "actions.passwordChange",
			set: func(c *v6okta.PasswordPolicyRule, p map[string]interface{}) {
				c.Actions.PasswordChange.AdditionalProperties = p
			},
			path: []string{"actions", "passwordChange"},
		},
		{
			name: "actions.selfServiceUnlock",
			set: func(c *v6okta.PasswordPolicyRule, p map[string]interface{}) {
				c.Actions.SelfServiceUnlock.AdditionalProperties = p
			},
			path: []string{"actions", "selfServiceUnlock"},
		},
		{
			name: "actions.selfServicePasswordReset",
			set: func(c *v6okta.PasswordPolicyRule, p map[string]interface{}) {
				c.Actions.SelfServicePasswordReset.AdditionalProperties = p
			},
			path: []string{"actions", "selfServicePasswordReset"},
		},
		{
			name: "actions.selfServicePasswordReset.requirement",
			set: func(c *v6okta.PasswordPolicyRule, p map[string]interface{}) {
				c.Actions.SelfServicePasswordReset.Requirement.AdditionalProperties = p
			},
			path: []string{"actions", "selfServicePasswordReset", "requirement"},
		},
		{
			name: "actions.selfServicePasswordReset.requirement.primary",
			set: func(c *v6okta.PasswordPolicyRule, p map[string]interface{}) {
				c.Actions.SelfServicePasswordReset.Requirement.Primary.AdditionalProperties = p
			},
			path: []string{"actions", "selfServicePasswordReset", "requirement", "primary"},
		},
		{
			name: "actions.selfServicePasswordReset.requirement.stepUp",
			set: func(c *v6okta.PasswordPolicyRule, p map[string]interface{}) {
				c.Actions.SelfServicePasswordReset.Requirement.StepUp.AdditionalProperties = p
			},
			path: []string{"actions", "selfServicePasswordReset", "requirement", "stepUp"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			current := currentPasswordRule(t)
			// The fixture does not carry every node, so fill in the ones it lacks.
			current.Conditions.People.SetGroups(v6okta.GroupCondition{})

			test.set(current, map[string]interface{}{"unmodeled": "keep me"})

			rule := buildRuleFromConfig(t, config)
			preserveUnmodeledProperties(&rule, current)
			body := replacePolicyRuleBody(t, &rule)

			assertValue(t, body, "keep me", append(test.path, "unmodeled")...)
		})
	}
}

// TestPreserveUnmodeledPropertiesNilSafety covers the shapes the update path can
// hand over when Okta answers with a sparse rule, or with no rule at all.
func TestPreserveUnmodeledPropertiesNilSafety(t *testing.T) {
	built := buildRuleFromConfig(t, baseRuleConfig())

	tests := []struct {
		name    string
		built   *v6okta.PasswordPolicyRule
		current *v6okta.PasswordPolicyRule
	}{
		{name: "nil built", built: nil, current: currentPasswordRule(t)},
		{name: "nil current", built: &built, current: nil},
		{name: "both nil", built: nil, current: nil},
		{name: "empty current", built: &built, current: v6okta.NewPasswordPolicyRule()},
		{
			name:  "current with empty conditions and actions",
			built: &built,
			current: func() *v6okta.PasswordPolicyRule {
				rule := v6okta.NewPasswordPolicyRule()
				rule.SetConditions(v6okta.PasswordPolicyRuleConditions{})
				rule.SetActions(v6okta.PasswordPolicyRuleActions{})
				return rule
			}(),
		},
		{
			name: "built with empty conditions and actions",
			built: func() *v6okta.PasswordPolicyRule {
				rule := v6okta.NewPasswordPolicyRule()
				rule.SetConditions(v6okta.PasswordPolicyRuleConditions{})
				rule.SetActions(v6okta.PasswordPolicyRuleActions{})
				return rule
			}(),
			current: currentPasswordRule(t),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// A panic here is the failure. There is nothing else to assert: a
			// node missing on either side has no properties to carry over.
			preserveUnmodeledProperties(test.built, test.current)
		})
	}
}
