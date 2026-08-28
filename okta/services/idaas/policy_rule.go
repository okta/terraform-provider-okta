package idaas

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/cenkalti/backoff/v4"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
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
			Description: "Rule priority. This attribute can be set to a valid priority. To avoid an endless diff situation an error is thrown if an invalid property is provided. The Okta API defaults to the last (lowest) if not provided.",
			// Suppress diff if config is empty.
			DiffSuppressFunc: utils.CreateValueDiffSuppression("0"),
		},
		"status": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     StatusActive,
			Description: "Policy Rule Status: `ACTIVE` or `INACTIVE`. Default: `ACTIVE`",
		},
		"network_connection": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Network selection mode: `ANYWHERE`, `ZONE`, `ON_NETWORK`, or `OFF_NETWORK`. Default: `ANYWHERE`",
			Default:     "ANYWHERE",
		},
		"network_includes": {
			Type:          schema.TypeList,
			Optional:      true,
			Description:   "Required if `network_connection` = `ZONE`. Indicates the network zones to include.",
			ConflictsWith: []string{"network_excludes"},
			Elem:          &schema.Schema{Type: schema.TypeString},
		},
		"network_excludes": {
			Type:          schema.TypeList,
			Optional:      true,
			Description:   "Required if `network_connection` = `ZONE`. Indicates the network zones to exclude.",
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
	return utils.BuildSchema(baseRuleSchema, target)
}

func buildRuleSchema(target map[string]*schema.Schema) map[string]*schema.Schema {
	return utils.BuildSchema(baseRuleSchema, target, userExcludedSchema)
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
	boc := utils.NewExponentialBackOffWithContext(ctx, backoff.DefaultMaxElapsedTime)
	err = backoff.Retry(func() error {
		ruleObj, resp, err := getAPISupplementFromMetadata(m).CreatePolicyRule(ctx, policyID, template)
		if doNotRetry(m, err) {
			return backoff.Permanent(err)
		}
		if err != nil {
			return backoff.Permanent(err)
		}
		if resp.StatusCode == http.StatusInternalServerError {
			return err
		}
		rule = ruleObj
		return nil
	}, boc)
	if err != nil {
		return fmt.Errorf("failed to create policy rule: %v", err)
	}
	status := d.Get("status").(string)
	if status == StatusInactive {
		_, err = getOktaClientFromMetadata(m).Policy.DeactivatePolicyRule(ctx, policyID, rule.Id)
		if err != nil {
			return fmt.Errorf("failed to deactivate policy rule on creation: %v", err)
		}
	}
	// We want to put this under Terraform's control even if priority is invalid.
	d.SetId(rule.Id)
	return utils.ValidatePriority(template.Priority, rule.Priority)
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

// ensureNotDefaultRule is a create-time guard only: Okta creates a policy's default
// rule itself, so a rule by that name can never be POSTed. Do not reintroduce it into
// the update or delete paths, that is exactly the GH-2788 regression. Those paths key
// off the API's `system` flag instead, since a rule's name is not authoritative.
func ensureNotDefaultRule(d *schema.ResourceData) error {
	return utils.EnsureNotDefault(d, "Rule")
}

func buildPolicyNetworkCondition(d *schema.ResourceData) *sdk.PolicyNetworkCondition {
	return &sdk.PolicyNetworkCondition{
		Connection: d.Get("network_connection").(string),
		Exclude:    utils.ConvertInterfaceToStringArrNullable(d.Get("network_excludes")),
		Include:    utils.ConvertInterfaceToStringArrNullable(d.Get("network_includes")),
	}
}

func getPolicyRule(ctx context.Context, d *schema.ResourceData, m interface{}) (*sdk.SdkPolicyRule, error) {
	client := getAPISupplementFromMetadata(m)
	policyID := d.Get("policy_id").(string)
	if policyID == "" {
		return nil, fmt.Errorf("'policy_id' field should be set")
	}
	policy, resp, err := client.GetPolicy(ctx, policyID)
	if err := utils.SuppressErrorOn404(resp, err); err != nil {
		return nil, err
	}
	if policy == nil {
		d.SetId("")
		return nil, nil
	}
	rule, resp, err := client.GetPolicyRule(ctx, policyID, d.Id())
	if err := utils.SuppressErrorOn404(resp, err); err != nil {
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
				Exclude: utils.ConvertInterfaceToStringSet(exclude),
			},
		}
	}

	return people
}

func syncRuleFromUpstream(d *schema.ResourceData, rule *sdk.SdkPolicyRule) error {
	_ = d.Set("name", rule.Name)
	_ = d.Set("status", rule.Status)
	_ = d.Set("priority", rule.Priority)
	_ = d.Set("system", utils.BoolFromBoolPtr(rule.System))
	// A policy's default (system) rule can come back with sparse or entirely absent
	// conditions, so substitute the schema defaults for anything the API omitted.
	network := &sdk.PolicyNetworkCondition{Connection: "ANYWHERE"}
	var usersExcluded []string
	if rule.Conditions != nil {
		if rule.Conditions.Network != nil {
			network = rule.Conditions.Network
		}
		if rule.Conditions.People != nil && rule.Conditions.People.Users != nil {
			usersExcluded = rule.Conditions.People.Users.Exclude
		}
	}
	_ = d.Set("network_connection", network.Connection)
	m := map[string]interface{}{
		"users_excluded": utils.ConvertStringSliceToSetNullable(usersExcluded),
	}
	if len(network.Include) > 0 {
		m["network_includes"] = utils.ConvertStringSliceToInterfaceSlice(network.Include)
	}
	if len(network.Exclude) > 0 {
		m["network_excludes"] = utils.ConvertStringSliceToInterfaceSlice(network.Exclude)
	}
	if network.Connection != "ANYWHERE" {
		return utils.SetNonPrimitives(d, m)
	}
	return utils.SetNonPrimitives(d, map[string]interface{}{
		"users_excluded": utils.ConvertStringSliceToSetNullable(usersExcluded),
	})
}

func updateRule(ctx context.Context, d *schema.ResourceData, m interface{}, template sdk.SdkPolicyRule) error {
	logger(m).Info("updating policy rule", "name", d.Get("name").(string))
	policyID := d.Get("policy_id").(string)
	if policyID == "" {
		return fmt.Errorf("'policy_id' field should be set")
	}

	// Okta allows a policy's default rule to be edited, it only forbids creating and
	// deleting it. A few attributes are managed by Okta on that rule though, so drop
	// them from the payload. This keys off the API's `system` flag, synced into state
	// by the read, rather than the configured name, which is not authoritative.
	if d.Get("system").(bool) {
		// Conditions can't be set on the default/system rule. A nil Conditions
		// marshals to "conditions": null since that field has no omitempty, which is
		// the same payload the profile enrollment resource has always sent for its
		// own system rule. Should Okta ever reject it for this rule type, fetch the
		// rule and send its conditions back verbatim instead.
		template.Conditions = nil
		template.System = utils.BoolPtr(true)
		// usePersistentCookie is read-only on the default rule, and *bool plus
		// omitempty means nil omits it from the request.
		if template.Actions.SignOn != nil && template.Actions.SignOn.Session != nil {
			template.Actions.SignOn.Session.UsePersistentCookie = nil
		}
	}

	rule, _, err := getAPISupplementFromMetadata(m).UpdatePolicyRule(ctx, policyID, d.Id(), template)
	if err != nil {
		return err
	}
	err = utils.ValidatePriority(template.Priority, rule.Priority)
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
	if d.Get("status").(string) == StatusActive {
		_, err := client.ActivatePolicyRule(ctx, policyID, d.Id())
		if err != nil {
			return fmt.Errorf("activation has failed: %v", err)
		}
	}
	if d.Get("status").(string) == StatusInactive {
		_, err := client.DeactivatePolicyRule(ctx, policyID, d.Id())
		if err != nil {
			return fmt.Errorf("deactivation has failed: %v", err)
		}
	}
	return nil
}

func deleteRule(ctx context.Context, d *schema.ResourceData, m interface{}, checkIsSystemPolicy bool) error {
	logger(m).Info("deleting policy rule", "name", d.Get("name").(string))
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

func buildKeepMeSignedIn(d *schema.ResourceData) *sdk.KeepMeSignedIn {
	v, ok := d.GetOk("keep_me_signed_in")
	if !ok {
		return nil
	}
	list := v.([]interface{})
	if len(list) == 0 || list[0] == nil {
		return nil
	}
	m := list[0].(map[string]interface{})
	return &sdk.KeepMeSignedIn{
		PostAuth:                m["post_auth"].(string),
		PostAuthPromptFrequency: m["post_auth_prompt_frequency"].(string),
	}
}

func flattenKeepMeSignedIn(k *sdk.KeepMeSignedIn) []interface{} {
	if k == nil {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"post_auth":                  k.PostAuth,
			"post_auth_prompt_frequency": k.PostAuthPromptFrequency,
		},
	}
}
