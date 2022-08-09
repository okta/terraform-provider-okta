package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/terraform-provider-okta/sdk"
)

// resourcePolicyMfaDefault requires Org Feature Flag OKTA_MFA_POLICY
func resourcePolicyMfaDefault() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyMfaDefaultCreateOrUpdate,
		ReadContext:   resourcePolicyMfaDefaultRead,
		UpdateContext: resourcePolicyMfaDefaultCreateOrUpdate,
		DeleteContext: resourcePolicyMfaDefaultDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				_, err := setDefaultPolicy(ctx, d, m, sdk.MfaPolicyType)
				if err != nil {
					return nil, err
				}
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: buildDefaultMfaPolicySchema(buildFactorSchemaProviders()),
	}
}

func resourcePolicyMfaDefaultCreateOrUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Id()
	if id == "" {
		policy, err := setDefaultPolicy(ctx, d, m, sdk.MfaPolicyType)
		if err != nil {
			return diag.FromErr(err)
		}
		id = policy.Id
	}
	_, _, err := getAPISupplementFromMetadata(m).UpdatePolicy(ctx, id, buildDefaultMFAPolicy(d))
	if err != nil {
		return diag.Errorf("failed to update default MFA policy: %v", err)
	}
	return resourcePolicyMfaDefaultRead(ctx, d, m)
}

func resourcePolicyMfaDefaultRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	policy, err := getPolicy(ctx, d, m)
	if err != nil {
		return diag.Errorf("failed to get default MFA policy: %v", err)
	}
	if policy == nil {
		return nil
	}

	syncSettings(d, policy.Settings)

	return nil
}

// Default policy can not be removed
func resourcePolicyMfaDefaultDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func buildDefaultMFAPolicy(d *schema.ResourceData) sdk.Policy {
	policy := sdk.MfaPolicy()
	policy.Name = d.Get("name").(string)
	policy.Status = d.Get("status").(string)
	policy.Description = d.Get("description").(string)
	policy.Priority = int64(d.Get("priority").(int))
	policy.Settings = buildSettings(d)
	policy.Conditions = &okta.PolicyRuleConditions{
		People: &okta.PolicyPeopleCondition{
			Groups: &okta.GroupCondition{
				Include: []string{d.Get("default_included_group_id").(string)},
			},
		},
	}
	return policy
}
