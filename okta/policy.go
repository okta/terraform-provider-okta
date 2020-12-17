package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
)

// Basis of policy schema
var (
	basePolicySchema = map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Policy Name",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Policy Description",
		},
		"priority": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Policy Priority, this attribute can be set to a valid priority. To avoid endless diff situation we error if an invalid priority is provided. API defaults it to the last/lowest if not there.",
			// Suppress diff if config is empty.
			DiffSuppressFunc: createValueDiffSuppression("0"),
		},
		"status": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          statusActive,
			ValidateDiagFunc: stringInSlice([]string{statusActive, statusInactive}),
			Description:      "Policy Status: ACTIVE or INACTIVE.",
		},
		"groups_included": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "List of Group IDs to Include",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}

	statusSchema = &schema.Schema{
		Type:             schema.TypeString,
		Optional:         true,
		Default:          statusActive,
		ValidateDiagFunc: stringInSlice([]string{statusActive, statusInactive}),
	}
)

func getPeopleConditions(d *schema.ResourceData) *okta.GroupRulePeopleCondition {
	return &okta.GroupRulePeopleCondition{
		Groups: &okta.GroupRuleGroupCondition{
			Include: convertInterfaceToStringSet(d.Get("group_whitelist")),
			Exclude: convertInterfaceToStringSet(d.Get("group_blacklist")),
		},
		Users: &okta.GroupRuleUserCondition{
			Include: convertInterfaceToStringSet(d.Get("user_whitelist")),
			Exclude: convertInterfaceToStringSet(d.Get("user_blacklist")),
		},
	}
}

func buildPolicySchema(target map[string]*schema.Schema) map[string]*schema.Schema {
	return buildSchema(basePolicySchema, target)
}

func createPolicy(ctx context.Context, d *schema.ResourceData, m interface{}, template sdk.Policy) error {
	logger(m).Info("creating policy", "name", template.Name, "type", template.Type)
	if err := ensureNotDefaultPolicy(d); err != nil {
		return err
	}
	policy, _, err := getSupplementFromMetadata(m).CreatePolicy(ctx, template)
	if err != nil {
		return err
	}
	d.SetId(policy.Id)
	// Even if priority is invalid we want to add the policy to Terraform to reflect upstream.
	err = validatePriority(template.Priority, policy.Priority)
	if err != nil {
		return err
	}

	return policyActivate(ctx, d, m)
}

func ensureNotDefaultPolicy(d *schema.ResourceData) error {
	return ensureNotDefault(d, "Policy")
}

func getGroups(d *schema.ResourceData) *okta.PolicyPeopleCondition {
	var people *okta.PolicyPeopleCondition
	if include, ok := d.GetOk("groups_included"); ok {
		people = &okta.PolicyPeopleCondition{
			Groups: &okta.GroupCondition{
				Include: convertInterfaceToStringSet(include),
			},
		}
	}
	return people
}

// Grabs policy from upstream, if the resource does not exist the returned policy will be nil which is not considered an error
func getPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) (*sdk.Policy, error) {
	logger(m).Info("getting policy", "id", d.Id())
	policy, resp, err := getSupplementFromMetadata(m).GetPolicy(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return nil, err
	}
	if policy == nil {
		d.SetId("")
		return nil, nil
	}
	return policy, err
}

// activate or deactivate a policy according to the terraform schema status field
func policyActivate(ctx context.Context, d *schema.ResourceData, m interface{}) error {
	logger(m).Info("changing policy's status", "id", d.Id(), "status", d.Get("status").(string))
	client := getOktaClientFromMetadata(m)

	if d.Get("status").(string) == statusActive {
		_, err := client.Policy.ActivatePolicy(ctx, d.Id())
		if err != nil {
			return fmt.Errorf("activation has failed: %v", err)
		}
	}
	if d.Get("status").(string) == statusInactive {
		_, err := client.Policy.DeactivatePolicy(ctx, d.Id())
		if err != nil {
			return fmt.Errorf("dctivation has failed: %v", err)
		}
	}
	return nil
}

func updatePolicy(ctx context.Context, d *schema.ResourceData, m interface{}, template sdk.Policy) error {
	logger(m).Info("updating policy", "name", d.Get("name").(string))
	if err := ensureNotDefaultPolicy(d); err != nil {
		return err
	}
	policy, _, err := getSupplementFromMetadata(m).UpdatePolicy(ctx, d.Id(), template)
	if err != nil {
		return err
	}
	// avoiding perpetual diffs by erroring when the configured priority is not valid and the API defaults it.
	err = validatePriority(template.Priority, policy.Priority)
	if err != nil {
		return err
	}
	return policyActivate(ctx, d, m)
}

func deletePolicy(ctx context.Context, d *schema.ResourceData, m interface{}) error {
	if err := ensureNotDefaultPolicy(d); err != nil {
		return err
	}
	logger(m).Info("deleting policy", "id", d.Id())
	client := getOktaClientFromMetadata(m)
	_, err := client.Policy.DeletePolicy(ctx, d.Id())
	if err != nil {
		return err
	}
	// remove the policy resource from terraform
	d.SetId("")
	return nil
}

func syncPolicyFromUpstream(d *schema.ResourceData, policy *sdk.Policy) error {
	_ = d.Set("name", policy.Name)
	_ = d.Set("description", policy.Description)
	_ = d.Set("status", policy.Status)
	_ = d.Set("priority", policy.Priority)
	return setNonPrimitives(d, map[string]interface{}{
		"groups_included": convertStringSetToInterface(policy.Conditions.People.Groups.Include),
	})
}
