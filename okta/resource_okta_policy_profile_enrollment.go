package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourcePolicyProfileEnrollment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyProfileEnrollmentCreate,
		ReadContext:   resourcePolicyProfileEnrollmentRead,
		UpdateContext: resourcePolicyProfileEnrollmentUpdate,
		DeleteContext: resourcePolicyProfileEnrollmentDelete,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the policy",
			},
			"status": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: elemInSlice([]string{statusActive, statusInactive}),
				Description:      "Status of the policy",
				Default:          statusActive,
			},
		},
	}
}

func resourcePolicyProfileEnrollmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	policy, _, err := getSupplementFromMetadata(m).CreatePolicy(ctx, buildPolicyProfileEnrollment(d))
	if err != nil {
		return diag.Errorf("failed to create profile enrollment policy: %v", err)
	}
	d.SetId(policy.Id)
	status, ok := d.GetOk("status")
	if ok && status.(string) == statusInactive {
		_, err = getOktaClientFromMetadata(m).Policy.DeactivatePolicy(ctx, policy.Id)
		if err != nil {
			return diag.Errorf("failed to deactivate profile enrollment policy: %v", err)
		}
	}
	return resourcePolicyProfileEnrollmentRead(ctx, d, m)
}

func resourcePolicyProfileEnrollmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	policy, err := getPolicy(ctx, d, m)
	if err != nil {
		return diag.Errorf("failed to get profile enrollment policy: %v", err)
	}
	if policy == nil {
		return nil
	}
	_ = d.Set("name", policy.Name)
	_ = d.Set("status", policy.Status)
	return nil
}

func resourcePolicyProfileEnrollmentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	_, _, err := getSupplementFromMetadata(m).UpdatePolicy(ctx, d.Id(), buildPolicyProfileEnrollment(d))
	if err != nil {
		return diag.Errorf("failed to update profile enrollment policy: %v", err)
	}
	oldStatus, newStatus := d.GetChange("status")
	if oldStatus != newStatus {
		if newStatus == statusActive {
			_, err = getOktaClientFromMetadata(m).Policy.ActivatePolicy(ctx, d.Id())
		} else {
			_, err = getOktaClientFromMetadata(m).Policy.DeactivatePolicy(ctx, d.Id())
		}
		if err != nil {
			return diag.Errorf("failed to change profile enrollment policy status: %v", err)
		}
	}
	return resourcePolicyProfileEnrollmentRead(ctx, d, m)
}

func resourcePolicyProfileEnrollmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := deletePolicy(ctx, d, m)
	if err != nil {
		return diag.Errorf("failed to delete profile enrollment policy: %v", err)
	}
	return nil
}

// build profile enrollment policy from schema data
func buildPolicyProfileEnrollment(d *schema.ResourceData) sdk.Policy {
	policy := sdk.ProfileEnrollmentPolicy()
	policy.Name = d.Get("name").(string)
	policy.Status = d.Get("status").(string)
	return policy
}
