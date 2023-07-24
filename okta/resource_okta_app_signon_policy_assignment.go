package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppSignOnAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppSignOnAssignmentCreate,
		ReadContext:   resourceAppSignOnAssignmentRead,
		DeleteContext: resourceAppSignOnAssignmentDelete,
		UpdateContext: resourceAppSignOnAssignmentUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				_ = d.Set("app_id", d.Id())
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAppSignOnAssignmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getAPISupplementFromMetadata(m)
	appId := d.Get("app_id").(string)
	policyId := d.Get("policy_id").(string)
	assignment, err := client.GetAppSignOnPolicyRuleAssigment(ctx, appId)
	if err != nil {
		return diag.Errorf("failed get app by ID: %v", err)
	}

	if assignment.PolicyId != policyId {
		_, err = client.SetAppSignOnPolicyRuleAssigment(ctx, appId, policyId)
		if err != nil {
			return diag.Errorf("failed to update the application sign on policy assignment: %v", err)
		}
	}

	err = d.Set("policy_id", policyId)
	if err != nil {
		return diag.Errorf("failed to set policy_id: %v", err)
	}

	d.SetId(appId)

	return nil
}

func resourceAppSignOnAssignmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	appId := d.Get("app_id").(string)
	assignment, err := getAPISupplementFromMetadata(m).GetAppSignOnPolicyRuleAssigment(ctx, appId)
	if err != nil {
		return diag.Errorf("failed get app by ID: %v", err)
	}

	err = d.Set("policy_id", assignment.PolicyId)
	if err != nil {
		return diag.Errorf("failed to set policy_id: %v", err)
	}

	d.SetId(appId)

	return nil
}

func resourceAppSignOnAssignmentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getAPISupplementFromMetadata(m)
	appId := d.Get("app_id").(string)
	policyId := d.Get("policy_id").(string)
	assignment, err := client.GetAppSignOnPolicyRuleAssigment(ctx, appId)
	if err != nil {
		return diag.Errorf("failed get app by ID: %v", err)
	}

	if assignment.PolicyId != policyId {
		_, err = client.SetAppSignOnPolicyRuleAssigment(ctx, appId, policyId)
		if err != nil {
			return diag.Errorf("failed to update the application sign on policy assignment: %v", err)
		}
	}

	err = d.Set("policy_id", assignment.PolicyId)
	if err != nil {
		return diag.Errorf("failed to set policy_id: %v", err)
	}

	d.SetId(appId)

	return nil
}

func resourceAppSignOnAssignmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}
