package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

// resourcePolicyMfaDefault requires Org Feature Flag OKTA_MFA_POLICY
func resourcePolicyMfaDefault() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyMfaDefaultCreateOrUpdate,
		ReadContext:   resourcePolicyMfaDefaultRead,
		UpdateContext: resourcePolicyMfaDefaultCreateOrUpdate,
		DeleteContext: resourceFuncNoOp,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				_, err := setDefaultPolicy(ctx, d, m, sdk.MfaPolicyType)
				if err != nil {
					return nil, err
				}
				return []*schema.ResourceData{d}, nil
			},
		},
		Description: `Configures default MFA Policy.
This resource allows you to configure default MFA Policy.
~> Requires Org Feature Flag 'OKTA_MFA_POLICY'. [Contact support](mailto:dev-inquiries@okta.com) to have this feature flag ***enabled***.
~> Unless Org Feature Flag 'ENG_ENABLE_OPTIONAL_PASSWORD_ENROLLMENT' is ***disabled*** 'okta_password' or 'okta_email' must be present and its 'enroll' value set to 'REQUIRED'. [Contact support](mailto:dev-inquiries@okta.com) to have this feature flag ***disabled***.`,
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

func buildDefaultMFAPolicy(d *schema.ResourceData) sdk.SdkPolicy {
	policy := sdk.MfaPolicy()
	policy.Name = d.Get("name").(string)
	policy.Status = d.Get("status").(string)
	policy.Description = d.Get("description").(string)
	policy.PriorityPtr = int64Ptr(d.Get("priority").(int))
	policy.Settings = buildSettings(d)
	policy.Conditions = &sdk.PolicyRuleConditions{
		People: &sdk.PolicyPeopleCondition{
			Groups: &sdk.GroupCondition{
				Include: []string{d.Get("default_included_group_id").(string)},
			},
		},
	}
	return policy
}
