package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
)

func resourcePolicySignOn() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicySignOnCreate,
		ReadContext:   resourcePolicySignOnRead,
		UpdateContext: resourcePolicySignOnUpdate,
		DeleteContext: resourcePolicySignOnDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: basePolicySchema,
	}
}

func resourcePolicySignOnCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	template := buildSignOnPolicy(d)
	err := createPolicy(ctx, d, m, template)
	if err != nil {
		return diag.Errorf("failed to create sign-on policy: %v", err)
	}
	return resourcePolicySignOnRead(ctx, d, m)
}

func resourcePolicySignOnRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	policy, err := getPolicy(ctx, d, m)
	if err != nil {
		return diag.Errorf("failed to get sign-on policy: %v", err)
	}
	if policy == nil {
		return nil
	}
	err = syncPolicyFromUpstream(d, policy)
	if err != nil {
		return diag.Errorf("failed to set sign-on policy: %v", err)
	}
	return nil
}

func resourcePolicySignOnUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	template := buildSignOnPolicy(d)
	err := updatePolicy(ctx, d, m, template)
	if err != nil {
		return diag.Errorf("failed to update sign-on policy: %v", err)
	}
	return resourcePolicySignOnRead(ctx, d, m)
}

func resourcePolicySignOnDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := deletePolicy(ctx, d, m)
	if err != nil {
		return diag.Errorf("failed to delete sign-on policy: %v", err)
	}
	return nil
}

// create or update a sign on policy
func buildSignOnPolicy(d *schema.ResourceData) sdk.Policy {
	template := sdk.SignOnPolicy()
	template.Name = d.Get("name").(string)
	template.Status = d.Get("status").(string)
	if description, ok := d.GetOk("description"); ok {
		template.Description = description.(string)
	}
	if priority, ok := d.GetOk("priority"); ok {
		template.Priority = int64(priority.(int))
	}
	template.Conditions = &okta.PolicyRuleConditions{
		People: getGroups(d),
	}
	return template
}
