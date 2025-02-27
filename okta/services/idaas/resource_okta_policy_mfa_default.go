package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
)

// resourcePolicyMfaDefault requires Org Feature Flag OKTA_MFA_POLICY
func resourcePolicyMfaDefault() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyMfaDefaultCreateOrUpdate,
		ReadContext:   resourcePolicyMfaDefaultRead,
		UpdateContext: resourcePolicyMfaDefaultCreateOrUpdate,
		DeleteContext: utils.ResourceFuncNoOp,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				_, err := setDefaultPolicy(ctx, d, meta, sdk.MfaPolicyType)
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

func resourcePolicyMfaDefaultCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var id string
	// Issue #2107, where adding a new MFA_ENROLL policy change the priority of the default policy, leading to the default policy unable to update
	// It is now required that the default policy is set again for every create and update, and the only thing that can be change is factor/authenticator
	policy, err := setDefaultPolicy(ctx, d, meta, sdk.MfaPolicyType)
	if err != nil {
		return diag.FromErr(err)
	}
	id = policy.Id

	_, _, err = getAPISupplementFromMetadata(meta).UpdatePolicy(ctx, id, buildDefaultMFAPolicy(d))
	if err != nil {
		return diag.Errorf("failed to update default MFA policy: %v", err)
	}
	return resourcePolicyMfaDefaultRead(ctx, d, meta)
}

func resourcePolicyMfaDefaultRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	policy, err := getPolicy(ctx, d, meta)
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
	policy.PriorityPtr = utils.Int64Ptr(d.Get("priority").(int))
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
