package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
)

func ResourcePolicySignOn() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicySignOnCreate,
		ReadContext:   resourcePolicySignOnRead,
		UpdateContext: resourcePolicySignOnUpdate,
		DeleteContext: resourcePolicySignOnDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "Creates a Sign On Policy. This resource allows you to create and configure a Sign On Policy.",
		Schema:      basePolicySchema,
	}
}

func resourcePolicySignOnCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	template := buildSignOnPolicy(d)
	err := createPolicy(ctx, d, meta, template)
	if err != nil {
		return diag.Errorf("failed to create sign-on policy: %v", err)
	}
	return resourcePolicySignOnRead(ctx, d, meta)
}

func resourcePolicySignOnRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	policy, err := getPolicy(ctx, d, meta)
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

func resourcePolicySignOnUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	template := buildSignOnPolicy(d)
	err := updatePolicy(ctx, d, meta, template)
	if err != nil {
		return diag.Errorf("failed to update sign-on policy: %v", err)
	}
	return resourcePolicySignOnRead(ctx, d, meta)
}

func resourcePolicySignOnDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := deletePolicy(ctx, d, meta)
	if err != nil {
		return diag.Errorf("failed to delete sign-on policy: %v", err)
	}
	return nil
}

// create or update a sign-on policy
func buildSignOnPolicy(d *schema.ResourceData) sdk.SdkPolicy {
	template := sdk.SignOnPolicy()
	template.Name = d.Get("name").(string)
	template.Status = d.Get("status").(string)
	if description, ok := d.GetOk("description"); ok {
		template.Description = description.(string)
	}
	if priority, ok := d.GetOk("priority"); ok {
		template.PriorityPtr = utils.Int64Ptr(priority.(int))
	}
	template.Conditions = &sdk.PolicyRuleConditions{
		People: getGroups(d),
	}
	return template
}
