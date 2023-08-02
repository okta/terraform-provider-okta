package okta

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/cenkalti/backoff/v4"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

var (
	userExcludedSchema = map[string]*schema.Schema{
		"users_excluded": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Set of User IDs to Exclude",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}

	// Basis of policy rules
	baseRuleSchema = map[string]*schema.Schema{
		"policy_id": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
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
			Description: "Policy Rule Priority, this attribute can be set to a valid priority. To avoid endless diff situation we error if an invalid priority is provided. API defaults it to the last (lowest) if not there.",
			// Suppress diff if config is empty.
			DiffSuppressFunc: createValueDiffSuppression("0"),
		},
		"status": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     statusActive,
			Description: "Policy Rule Status: ACTIVE or INACTIVE.",
		},
		"network_connection": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Network selection mode: ANYWHERE, ZONE, ON_NETWORK, or OFF_NETWORK.",
			Default:     "ANYWHERE",
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

	appResource = &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
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
)

func buildBaseRuleSchema(target map[string]*schema.Schema) map[string]*schema.Schema {
	return buildSchema(baseRuleSchema, target)
}

func buildRuleSchema(target map[string]*schema.Schema) map[string]*schema.Schema {
	return buildSchema(baseRuleSchema, target, userExcludedSchema)
}

func createRule(ctx context.Context, d *schema.ResourceData, m interface{}, template sdk.SdkPolicyRule, ruleType string) error {
	logger(m).Info("creating policy rule", "name", d.Get("name").(string))
	err := ensureNotDefaultRule(d)
	if err != nil {
		return err
	}
	policyID := d.Get("policy_id").(string)
	if policyID == "" {
		return fmt.Errorf("'policy_id' field should be set")
	}
	var rule *sdk.SdkPolicyRule
	boc := newExponentialBackOffWithContext(ctx, backoff.DefaultMaxElapsedTime)
	err = backoff.Retry(func() error {
		ruleObj, resp, err := getAPISupplementFromMetadata(m).CreatePolicyRule(ctx, policyID, template)
		if doNotRetry(m, err) {
			return backoff.Permanent(err)
		}
		if resp.StatusCode == http.StatusInternalServerError {
			return err
		}
		if err != nil {
			return backoff.Permanent(err)
		}
		rule = ruleObj
		return nil
	}, boc)
	if err != nil {
		return fmt.Errorf("failed to create policy rule: %v", err)
	}
	status := d.Get("status").(string)
	if status == statusInactive {
		_, err = getOktaClientFromMetadata(m).Policy.DeactivatePolicyRule(ctx, policyID, rule.Id)
		if err != nil {
			return fmt.Errorf("failed to deactivate policy rule on creation: %v", err)
		}
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
			_ = d.Set("policy_id", parts[0])
			d.SetId(parts[1])
			return []*schema.ResourceData{d}, nil
		},
	}
}

func ensureNotDefaultRule(d *schema.ResourceData) error {
	return ensureNotDefault(d, "Rule")
}

func buildPolicyNetworkCondition(d *schema.ResourceData) *sdk.PolicyNetworkCondition {
	return &sdk.PolicyNetworkCondition{
		Connection: d.Get("network_connection").(string),
		Exclude:    convertInterfaceToStringArrNullable(d.Get("network_excludes")),
		Include:    convertInterfaceToStringArrNullable(d.Get("network_includes")),
	}
}

func getPolicyRule(ctx context.Context, d *schema.ResourceData, m interface{}) (*sdk.SdkPolicyRule, error) {
	client := getAPISupplementFromMetadata(m)
	policyID := d.Get("policy_id").(string)
	if policyID == "" {
		return nil, fmt.Errorf("'policy_id' field should be set")
	}
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

func getUsers(d *schema.ResourceData) *sdk.PolicyPeopleCondition {
	var people *sdk.PolicyPeopleCondition

	if exclude, ok := d.GetOk("users_excluded"); ok {
		people = &sdk.PolicyPeopleCondition{
			Users: &sdk.UserCondition{
				Exclude: convertInterfaceToStringSet(exclude),
			},
		}
	}

	return people
}

func syncRuleFromUpstream(d *schema.ResourceData, rule *sdk.SdkPolicyRule) error {
	_ = d.Set("name", rule.Name)
	_ = d.Set("status", rule.Status)
	_ = d.Set("priority", rule.Priority)
	_ = d.Set("network_connection", rule.Conditions.Network.Connection)
	m := map[string]interface{}{
		"users_excluded": convertStringSliceToSetNullable(rule.Conditions.People.Users.Exclude),
	}
	if len(rule.Conditions.Network.Include) > 0 {
		m["network_includes"] = convertStringSliceToInterfaceSlice(rule.Conditions.Network.Include)
	}
	if len(rule.Conditions.Network.Exclude) > 0 {
		m["network_excludes"] = convertStringSliceToInterfaceSlice(rule.Conditions.Network.Exclude)
	}
	if rule.Conditions.Network.Connection != "ANYWHERE" {
		return setNonPrimitives(d, m)
	}
	return setNonPrimitives(d, map[string]interface{}{
		"users_excluded": convertStringSliceToSetNullable(rule.Conditions.People.Users.Exclude),
	})
}

func updateRule(ctx context.Context, d *schema.ResourceData, m interface{}, template sdk.SdkPolicyRule) error {
	logger(m).Info("updating policy rule", "name", d.Get("name").(string))
	if err := ensureNotDefaultRule(d); err != nil {
		return err
	}
	policyID := d.Get("policy_id").(string)
	if policyID == "" {
		return fmt.Errorf("'policy_id' field should be set")
	}
	rule, _, err := getAPISupplementFromMetadata(m).UpdatePolicyRule(ctx, policyID, d.Id(), template)
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
	client := getOktaClientFromMetadata(m).Policy
	policyID := d.Get("policy_id").(string)
	if policyID == "" {
		return fmt.Errorf("'policy_id' field should be set")
	}
	if d.Get("status").(string) == statusActive {
		_, err := client.ActivatePolicyRule(ctx, policyID, d.Id())
		if err != nil {
			return fmt.Errorf("activation has failed: %v", err)
		}
	}
	if d.Get("status").(string) == statusInactive {
		_, err := client.DeactivatePolicyRule(ctx, policyID, d.Id())
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
	if err != nil {
		return err
	}
	if rule == nil {
		return nil
	}
	shouldRemove := true
	if checkIsSystemPolicy {
		if rule.System != nil && *rule.System {
			logger(m).Info(fmt.Sprintf("Policy Rule '%s' is a System Policy, cannot delete from Okta", d.Get("name").(string)))
			shouldRemove = false
		}
	}
	if shouldRemove {
		policyID := d.Get("policy_id").(string)
		if policyID == "" {
			return fmt.Errorf("'policy_id' field should be set")
		}
		_, err = getOktaClientFromMetadata(m).Policy.DeletePolicyRule(ctx, policyID, d.Id())
		if err != nil {
			return err
		}
	}
	return nil
}
